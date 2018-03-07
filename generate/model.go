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
	packages             []packageGen

	Services map[string]serviceGen
}

type packageGen struct {
	Name     string
	FullName string
	Alias    *string
	Package  scan.Package
}

func (pDef packageGen) UniqueName() string {
	if pDef.Alias != nil {
		return *pDef.Alias
	}

	return pDef.Name
}

func (cg ContainerGenerator) getPackageByUniqueName(name string) *packageGen {
	for _, pkg := range cg.packages {
		if pkg.UniqueName() == name {
			return &pkg
		}
	}
	return nil
}

type serviceGen interface {
	Name() string
	ResultType() types.Type
	Generate(ContainerGenerator) string
}

type valueGen interface {
	Generate(ContainerGenerator) string
	NeedsPointer() bool
}
