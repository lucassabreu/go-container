package generate

import (
	"github.com/lucassabreu/go-container/scan"
)

// Service creates a new func to generate a service
type Service interface {
	Name() string
	ResultType() Type
	Package() *Package
	Generate(ContainerGenerator) string
}

type basicServiceGen struct {
	ServicePackage    *Package
	ServiceName       string
	ServiceResultType Type
}

func (b basicServiceGen) Name() string {
	return b.ServiceName
}

func (b basicServiceGen) ResultType() Type {
	return b.ServiceResultType
}

func (b basicServiceGen) Package() *Package {
	return b.ServicePackage
}

type serviceByFactoryGen struct {
	basicServiceGen
	factoryFunc scan.Func
	arguments   []Value
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
	values     map[string]Value
}

func (sd serviceByInitializationGen) Generate(cg ContainerGenerator) string {
	return ""
}
