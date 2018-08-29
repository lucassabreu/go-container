package generate

import (
	"bytes"
	"fmt"

	"github.com/lucassabreu/go-container/scan"
)

// Service creates a new func to generate a service
type Service interface {
	Name() string
	Package() *Package
	ResultType() Type
	Build(*ContainerGenerator) error
	String() string
}

type ServiceByFactory struct {
	name       string
	pkg        *Package
	resultType Type
	factory    scan.Func
	arguments  []Value
	buffer     *bytes.Buffer
}

func NewServiceByFactory(name string, pkg *Package, factory scan.Func, args []Value, resultType Type) *ServiceByFactory {
	return &ServiceByFactory{
		name:       name,
		pkg:        pkg,
		factory:    factory,
		arguments:  args,
		resultType: resultType,
	}
}

func writeReturnService(s Service, b *bytes.Buffer) {
	b.WriteString("return c.")
	b.WriteString(ToVarName(s.Name()))
	b.WriteRune('\n')
}

func writeServiceMethodOpening(s Service, cg *ContainerGenerator, b *bytes.Buffer) {
	b.WriteString("func (c *")
	b.WriteString(cg.ContainerName)
	b.WriteString(") Get")
	b.WriteString(s.Name())
	b.WriteString("() ")
	b.WriteString(s.ResultType().StringAsPointer())
	b.WriteString("{\nif c.")
	b.WriteString(ToVarName(s.Name()))
	b.WriteString(" != nil {\n")
	writeReturnService(s, b)
	b.WriteString("}\n")
}

func (s *ServiceByFactory) Build(cg *ContainerGenerator) error {
	b := &bytes.Buffer{}

	writeServiceMethodOpening(s, cg, b)

	for i, v := range s.arguments {
		v.Build(cg)
		if v.NeedsVariable() {
			v.GenerateVariable(fmt.Sprintf("v%d", i))
		}
	}

	b.WriteString("c.")
	b.WriteString(ToVarName(s.Name()))
	b.WriteRune('=')

	b.WriteString(s.Package().UniqueName())
	b.WriteRune('.')
	b.WriteString(s.factory.Name)
	b.WriteString("(\n")
	for _, v := range s.arguments {
		b.WriteString(v.GenerateUse())
		b.WriteString(",\n")
	}
	b.WriteString("\n)\n")

	writeReturnService(s, b)
	b.WriteString("}")

	s.buffer = b

	return nil
}

func (s *ServiceByFactory) String() string {
	return s.buffer.String()
}

func (s *ServiceByFactory) ResultType() Type {
	return s.resultType
}

func (s *ServiceByFactory) Name() string {
	return s.name
}

func (s *ServiceByFactory) Package() *Package {
	return s.pkg
}

type basicServiceGen struct {
	ServiceName       string
	ServiceResultType Type
	pkg               *Package
}

func (b basicServiceGen) Name() string {
	return b.ServiceName
}

func (b basicServiceGen) ResultType() Type {
	return b.ServiceResultType
}

func (b basicServiceGen) Package() *Package {
	return b.pkg
}

type serviceByFailableFactoryGen struct {
	*ServiceByFactory
}

type serviceByInitializationGen struct {
	basicServiceGen
	initStruct scan.Struct
	values     map[string]Value
}

func (sd serviceByInitializationGen) Build(cg *ContainerGenerator) error {
	return nil
}

func (sd serviceByInitializationGen) String() string {
	return ""
}
