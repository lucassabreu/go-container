package main

import (
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	pkgName := "github.com/Coderockr/vitrine-social/server/handlers"
	// typeName = "NewOrganizationHandler"

	pkg, err := build.Import(pkgName, wd, 0)
	if err != nil {
		panic(err)
	}

	fset := token.NewFileSet()
	for _, file := range pkg.GoFiles {
		f, err := parser.ParseFile(fset, filepath.Join(pkg.Dir, file), nil, 0)
		if err != nil {
			panic(err)
		}

		for _, decl := range f.Decls {
			// decl, ok := decl.(*ast.GenDecl)
			decl, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}
			// if !ok || (decl.Tok != token.TYPE && decl.Tok != token.FUNC) {
			// 	continue
			// }

			if decl.Name.Name != "NewOrganizationHandler" {
				continue
			}

			println(decl.Name.Name)
			for _, p := range decl.Type.Params.List {
				for _, n := range p.Names {
					print("\t" + n.Name)
				}
				p, ok := p.Type.(*ast.GenDecl)
				if !ok || p.Tok != token.TYPE {
					continue
				}

				for _, spec := range p.Specs {
					spec := spec.(*ast.TypeSpec)
					println(" " + spec.Name.Name)
				}
			}

			// for _, spec := range decl.Specs {
			// 	spec := spec.(*ast.TypeSpec)
			// 	println(spec.Name.Name)
			// 	// if spec.Name.Name != id {
			// 	// 	continue
			// 	// }
			// 	// return Pkg{Package: pkg, FileSet: fset}, spec, nil
			// }
		}
	}
}
