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
	packages             []packageDef

	Services map[string]serviceDef
}

// Packages imported for the container
func (cg ContainerGenerator) Packages() []packageDef {
	return cg.packages
}

type packageDef struct {
	Name     string
	FullName string
	Alias    *string
	Package  scan.Package
}

func (pDef packageDef) UniqueName() string {
	if pDef.Alias != nil {
		return *pDef.Alias
	}

	return pDef.Name
}

type serviceDef interface {
	Name() string
	ResultType() types.Type
	Generate(ContainerGenerator) string
}

type valueDef interface {
	Generate(ContainerGenerator, types.Type) string
}
