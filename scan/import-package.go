package scan

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"strings"
)

// Package represents a Go package definition
type Package struct {
	Name       string
	ImportPath string
	Funcs      map[string]Func
	Structs    map[string]Struct
}

// Func represents a Go func definition
type Func struct {
	Name    string
	Params  []types.Type
	Results []types.Type
}

// Struct represents a Go struct definition
type Struct struct {
	Name   string
	Fields map[string]types.Type
}

// ImportPackage will find and read the definitions on a package and return it
func ImportPackage(pkgName string) (Package, error) {
	wd, err := os.Getwd()
	pkgDef := Package{}

	pkg, err := build.Import(pkgName, wd, 0)
	if err != nil {
		return pkgDef, err
	}

	fset := token.NewFileSet()
	astFiles := []*ast.File{}
	for _, file := range pkg.GoFiles {
		f, err := parser.ParseFile(fset, filepath.Join(pkg.Dir, file), nil, 0)
		if err != nil {
			return pkgDef, err
		}

		astFiles = append(astFiles, f)
	}

	conf := types.Config{
		IgnoreFuncBodies: true,
		FakeImportC:      true,
		Importer:         importer.For("source", nil),
	}
	checkedPkg, err := conf.Check(pkgName, fset, astFiles, nil)
	if err != nil {
		return pkgDef, fmt.Errorf("Type check failed: %v", err)
	}

	pkgDef.ImportPath = checkedPkg.Path()
	pkgDef.Name = checkedPkg.Name()
	pkgDef.Funcs = make(map[string]Func)
	pkgDef.Structs = make(map[string]Struct)

	scope := checkedPkg.Scope()
	for _, name := range scope.Names() {
		if ast.IsExported(name) {
			obj := scope.Lookup(name)
			switch obj := obj.(type) {
			case *types.Func:
				sig := obj.Type().(*types.Signature)
				if sig.Recv() != nil {
					continue
				}

				pkgDef.Funcs[obj.Name()] = Func{
					Name:    obj.Name(),
					Params:  tupleToTypes(sig.Params()),
					Results: tupleToTypes(sig.Results()),
				}
			case *types.TypeName:
				typ, ok := obj.Type().Underlying().(*types.Struct)
				if !ok {
					continue
				}

				structDef := Struct{
					Name:   obj.Name(),
					Fields: make(map[string]types.Type),
				}

				for i := 0; i < typ.NumFields(); i++ {
					field := typ.Field(i)
					if !field.Exported() {
						continue
					}

					structDef.Fields[field.Name()] = field.Type()
				}

				pkgDef.Structs[obj.Name()] = structDef
			}
		}
	}

	return pkgDef, nil
}

func tupleToTypes(t *types.Tuple) []types.Type {
	typeList := make([]types.Type, t.Len())
	for i := 0; i < t.Len(); i++ {
		typeList[i] = t.At(i).Type()
	}
	return typeList
}

func (p Package) String() string {
	b := strings.Builder{}

	b.WriteString("Package: ")
	b.WriteString(p.Name)
	b.WriteString(" (")
	b.WriteString(p.ImportPath)
	b.WriteString(")\n\tFuncs:\n")
	for _, f := range p.Funcs {
		b.WriteString("\t\t")
		b.WriteString(f.Name)
		b.WriteRune('(')
		for i, t := range f.Params {
			b.WriteString(t.String())
			if i < (len(f.Params) - 1) {
				b.WriteString(", ")
			}
		}
		b.WriteString(") (")
		for i, t := range f.Results {
			b.WriteString(t.String())
			if i < (len(f.Results) - 1) {
				b.WriteString(", ")
			}
		}
		b.WriteString(")\n")
	}

	b.WriteString("\n\tStructs:\n")
	for _, f := range p.Structs {
		b.WriteString("\t\t")
		b.WriteString(f.Name)
		b.WriteString("{\n")
		for name, t := range f.Fields {
			b.WriteString("\t\t\t")
			b.WriteString(name)
			b.WriteRune(' ')
			b.WriteString(t.String())
			b.WriteString(",\n")
		}
		b.WriteString("\t\t}\n")
	}
	return b.String()
}
