package generate

import (
	"fmt"
	"go/types"
	"strings"

	"github.com/lucassabreu/go-container/def"
	"github.com/lucassabreu/go-container/scan"
)

// ContainerGenerator represents a container to be generated
type ContainerGenerator struct {
	ContainerName string
	ContainerDocs string

	importedPackageNames []string
	packages             []Package

	services map[string]Service
}

// NewContainerGenerator creates a ContainerGenerator for the the def.Container
func NewContainerGenerator(cDef def.Container) (cg ContainerGenerator, err error) {
	if err = CheckCircularReference(cDef); err != nil {
		return
	}

	cg.ContainerName = cDef.Name
	cg.ContainerDocs = cDef.Docs

	for _, pkg := range cDef.Packages {
		err := cg.RegisterPackage(pkg.Package, pkg.Alias)
		if err != nil {
			return ContainerGenerator{}, err
		}
	}

	for name, serv := range cDef.Services {
		if serv.IsByFactory() {
			if err = cg.registerServiceByFactory(name, *serv.Factory, serv.Arguments); err != nil {
				return ContainerGenerator{}, err
			}
		}

		if err = cg.registerServiceByInitialization(name, *serv.Struct, serv.Fields); err != nil {
			return ContainerGenerator{}, err
		}
	}

	return
}

// GetPackageByUniqueName returns the package by its "import" name
func (cg ContainerGenerator) GetPackageByUniqueName(name string) *Package {
	for _, pkg := range cg.packages {
		if pkg.UniqueName() == name {
			return &pkg
		}
	}
	return nil
}

// Services returns the registered services on the container
func (cg ContainerGenerator) Services() map[string]Service {
	return cg.services
}

// RegisterPackage add the package into the ContainerGenerator
func (cg ContainerGenerator) RegisterPackage(pkgPath string, alias *string) error {
	scannedPackage, err := scan.ImportPackage(pkgPath)
	if err != nil {
		return err
	}

	uniqueName := scannedPackage.Name
	if alias != nil {
		uniqueName = *alias
	}

	if contains(cg.importedPackageNames, uniqueName) {
		return fmt.Errorf("Aready exists a package with the imported name \"%s\"", uniqueName)
	}

	cg.packages = append(cg.packages, Package{
		alias:    alias,
		fullName: scannedPackage.ImportPath,
		name:     scannedPackage.Name,
		scanned:  scannedPackage,
	})
	return nil
}

func breakIntoPackageAndDef(ref string) (pkg, def string) {
	pieces := strings.Split(ref, ".")
	pkg = pieces[0]
	def = pieces[1]
	return
}

func (cg ContainerGenerator) registerServiceByFactory(name, factoryFunc string, args []def.Value) error {
	pkg, factoryFunc := breakIntoPackageAndDef(factoryFunc)

	pkgGen := cg.GetPackageByUniqueName(pkg)
	if pkgGen == nil {
		return fmt.Errorf("There is no imported package with name \"%s\"", pkg)
	}

	fnc, ok := pkgGen.ScannedPackage().Funcs[factoryFunc]
	if !ok {
		return fmt.Errorf("There is no func named \"%s\" at the package %s", factoryFunc, pkgGen.ScannedPackage().ImportPath)
	}

	if !fnc.Variadic && len(args) != len(fnc.Params) {
		return fmt.Errorf("Func %s.%s expects %d parameters, %d informmed", pkg, factoryFunc, len(fnc.Params), len(args))
	}

	if fnc.Variadic && len(args) < (len(fnc.Params)-1) {
		return fmt.Errorf("Func %s.%s expects at least %d parameters, only %d informmed", pkg, factoryFunc, len(fnc.Params)-1, len(args))
	}

	if len(fnc.Results) == 0 || len(fnc.Results) > 2 || (len(fnc.Results) == 2 && fnc.Results[1].String() != "error") {
		return fmt.Errorf("Func %s.%s should return one value, or one value and one error", pkg, factoryFunc)
	}

	paramTypes := fnc.Params
	for i := len(paramTypes); i < len(args); i++ {
		paramTypes = append(paramTypes, fnc.Params[:1][0])
	}

	values := make([]Value, len(paramTypes))
	for i, t := range paramTypes {
		v, err := cg.createValue(t, args[i])
		if err != nil {
			return err
		}
		values[i] = v
	}

	cg.services[name] = serviceByFactoryGen{
		basicServiceGen: basicServiceGen{
			ServiceName:       name,
			ServiceResultType: fnc.Results[0],
		},
		arguments: values,
	}

	if len(fnc.Results) == 2 {
		cg.services[name] = serviceByFailableFactoryGen{cg.services[name].(serviceByFactoryGen)}
	}

	return nil
}

func (cg ContainerGenerator) registerServiceByInitialization(name, structName string, fields map[string]def.Value) error {
	pkg, structName := breakIntoPackageAndDef(structName)

	pkgGen := cg.GetPackageByUniqueName(pkg)
	if pkgGen == nil {
		return fmt.Errorf("There is no imported package with name \"%s\"", pkg)
	}

	structType, ok := pkgGen.ScannedPackage().Structs[structName]
	if !ok {
		return fmt.Errorf("There is no struct named \"%s\" at the package %s", structName, pkgGen.ScannedPackage().ImportPath)
	}

	initValues := make(map[string]Value, len(fields))
	for fieldName, defValue := range fields {
		fieldType, ok := structType.Fields[fieldName]
		if !ok {
			return fmt.Errorf("There is no field \"%s\" at the struct %s.%s (service %s)", fieldName, pkg, structName, name)
		}

		v, err := cg.createValue(fieldType, defValue)
		if err != nil {
			return err
		}

		initValues[fieldName] = v
	}

	cg.services[name] = serviceByInitializationGen{
		basicServiceGen: basicServiceGen{
			ServiceName:       name,
			ServiceResultType: structType.Type,
		},
		initStruct: structType,
		values:     initValues,
	}

	return nil
}

func (cg ContainerGenerator) createValue(t types.Type, arg def.Value) (v Value, err error) {
	if arg.ValueType() == def.ValueSingle {
		v = constValue{
			value: arg.GetSingleValue(),
			typ:   t,
		}
		return
	}
	return nil, nil
}
