package generate

import (
	"fmt"

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

	cg.packages = append(cg.packages, packageDef{
		Alias:    alias,
		FullName: scannedPackage.ImportPath,
		Name:     scannedPackage.Name,
		Package:  scannedPackage,
	})
	return nil
}
