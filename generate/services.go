package generate

import (
	"go/types"

	"github.com/lucassabreu/go-container/scan"
)

type basicServiceDef struct {
	ServiceName       string
	ServiceRetultType types.Type
}

func (b basicServiceDef) Name() string {
	return b.ServiceName
}

func (b basicServiceDef) ResultType() types.Type {
	return b.ServiceRetultType
}

type serviceByFactoryDef struct {
	basicServiceDef
	factoryFunc scan.Func
	arguments   []valueDef
}

func (sd serviceByFactoryDef) Generate(cg ContainerGenerator) string {
	return ""
}

type serviceByFailableFactoryDef struct {
	serviceByFactoryDef
}
