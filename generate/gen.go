package generate

import (
	"github.com/lucassabreu/go-container/def"
)

// Generate creates a ContainerGenerator for the the def.Container, may fail
func Generate(cDef def.Container) (cg ContainerGenerator, err error) {
	if err = hasCircularReference(cDef); err != nil {
		return
	}

	return
}
