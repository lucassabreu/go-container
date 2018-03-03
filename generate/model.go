package generate

import (
	"go/types"

	"github.com/lucassabreu/go-container/scan"
)

// ContainerGenerator represents a container to be generated
type ContainerGenerator struct {
	ContainerName string
	ContainerDocs string

	Packages []packageDef

	Services map[string]serviceDef
}

type packageDef struct {
	Name     string
	FullName string
	Alias    *string
	Package  scan.Package
}

type serviceDef interface {
	Name() string
	ResultType() types.Type
	Generate(ContainerGenerator) string
}

type valueDef interface {
	Generate(ContainerGenerator, types.Type) string
}
