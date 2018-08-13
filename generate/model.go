package generate

import (
	"go/types"

	"github.com/lucassabreu/go-container/scan"
)

// ContainerGenerator represents a container to be generated
type ContainerGenerator struct {
	ContainerName string
	ContainerDocs string

	importedPackageNames []string
	packages             []Package

	Services map[string]ServiceGenerator
}

// Package represents a package to be imported by the container
type Package struct {
	name     string
	fullName string
	alias    *string
	scanned  scan.Package
}

// UniqueName package name or alias (if informed), should be unique
func (p Package) UniqueName() string {
	if p.alias != nil {
		return *p.alias
	}

	return p.name
}

// ScannedPackage returns the scan.Package which it is based
func (p Package) ScannedPackage() scan.Package {
	return p.scanned
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

// ServiceGenerator creates a new func to generate a service
type ServiceGenerator interface {
	Name() string
	ResultType() types.Type
	Generate(ContainerGenerator) string
}

type valueGen interface {
	Generate(ContainerGenerator) string
	NeedsPointer() bool
}
