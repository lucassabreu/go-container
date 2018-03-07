package generate

import (
	"fmt"
	"go/types"
	"strings"

	"github.com/lucassabreu/go-container/def"
	"github.com/lucassabreu/go-container/scan"
)

// Generate creates a ContainerGenerator for the the def.Container, may fail
func Generate(cDef def.Container) (cg ContainerGenerator, err error) {
	if err = hasCircularReference(cDef); err != nil {
		return
	}

	cg.ContainerName = cDef.Name
	cg.ContainerDocs = cDef.Docs

	for _, pkg := range cDef.Packages {
		err := cg.registerPackage(pkg.Package, pkg.Alias)
		if err != nil {
			return ContainerGenerator{}, err
		}
	}

	for name, serv := range cDef.Services {
		if serv.IsByFactory() {
			if err = cg.registerServiceByFactory(name, *serv.Factory, *serv.Arguments); err != nil {
				return ContainerGenerator{}, err
			}
		}

		if err = cg.registerServiceByInitialization(name, *serv.Struct, *serv.Fields); err != nil {
			return ContainerGenerator{}, err
		}
	}

	return
}

func (cg ContainerGenerator) registerPackage(pkgPath string, alias *string) error {
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

	cg.packages = append(cg.packages, packageGen{
		Alias:    alias,
		FullName: scannedPackage.ImportPath,
		Name:     scannedPackage.Name,
		Package:  scannedPackage,
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

	pkgGen := cg.getPackageByUniqueName(pkg)
	if pkgGen == nil {
		return fmt.Errorf("There is no imported package with name \"%s\"", pkg)
	}

	fnc, ok := pkgGen.Package.Funcs[factoryFunc]
	if !ok {
		return fmt.Errorf("There is no func named \"%s\" at the package %s", factoryFunc, pkgGen.Package.ImportPath)
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

	values := make([]valueGen, len(paramTypes))
	for i, t := range paramTypes {
		v, err := cg.createValueGen(t, args[i])
		if err != nil {
			return err
		}
		values[i] = v
	}

	cg.Services[name] = serviceByFactoryGen{
		basicServiceGen: basicServiceGen{
			ServiceName:       name,
			ServiceRetultType: fnc.Results[0],
		},
		arguments: values,
	}

	if len(fnc.Results) == 2 {
		cg.Services[name] = serviceByFailableFactoryGen{cg.Services[name].(serviceByFactoryGen)}
	}

	return nil
}

func (cg ContainerGenerator) registerServiceByInitialization(name, structName string, fields map[string]def.Value) error {
	return nil
}

func (cg ContainerGenerator) createValueGen(t types.Type, arg def.Value) (v valueGen, err error) {
	switch arg.ValueType() {
	case def.ValueSingle:
		v = constValueGen{
			value: arg.GetSingleValue(),
			typ:   t,
		}
		return
	}
	return nil, nil
}
