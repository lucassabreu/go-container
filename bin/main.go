package main

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
)

// Package represents a Go package definition
type Package struct {
	Name       string
	ImportPath string
	Funcs      map[string]interface{}
	Structs    map[string]interface{}
}

// Func represents a Go func definition
type Func struct {
	Name    string
	Params  []types.Type
	Results []types.Type
}

func main() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	testImport("github.com/Coderockr/vitrine-social/server/handlers", wd)
	testImport("gopkg.in/yaml.v2", wd)
}

func testImport(pkgName, wd string) {

	pkg, err := build.Import(pkgName, wd, 0)
	if err != nil {
		panic(err)
	}

	fset := token.NewFileSet()
	astFiles := []*ast.File{}
	for _, file := range pkg.GoFiles {
		f, err := parser.ParseFile(fset, filepath.Join(pkg.Dir, file), nil, 0)
		if err != nil {
			panic(err)
		}

		astFiles = append(astFiles, f)
	}

	conf := types.Config{
		IgnoreFuncBodies: true,
		FakeImportC:      true,
		Importer:         importer.Default(),
	}
	checkedPkg, err := conf.Check(pkgName, fset, astFiles, nil)
	if err != nil {
		panic(fmt.Errorf("Type check failed: %v", err))
	}

	pkgDef := Package{
		ImportPath: checkedPkg.Path(),
		Name:       checkedPkg.Name(),
		Funcs:      make(map[string]interface{}),
		Structs:    make(map[string]interface{}),
	}

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

			}

		}
	}

	fmt.Println(pkgDef)
}
