package generate

import (
	"go/types"

	"github.com/lucassabreu/go-container/scan"
)

type basicServiceGen struct {
	ServiceName       string
	ServiceResultType types.Type
}

func (b basicServiceGen) Name() string {
	return b.ServiceName
}

func (b basicServiceGen) ResultType() types.Type {
	return b.ServiceResultType
}

type serviceByFactoryGen struct {
	basicServiceGen
	factoryFunc scan.Func
	arguments   []valueGen
}

func (sd serviceByFactoryGen) Generate(cg ContainerGenerator) string {
	return ""
}

type serviceByFailableFactoryGen struct {
	serviceByFactoryGen
}

type serviceByInitializationGen struct {
	basicServiceGen
	initStruct scan.Struct
	values     map[string]valueGen
}

func (sd serviceByInitializationGen) Generate(cg ContainerGenerator) string {
	return ""
}
