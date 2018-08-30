package generate

import (
	"bytes"
	"fmt"
	"go/types"
	"strconv"

	"github.com/lucassabreu/go-container/def"
)

func (cg ContainerGenerator) createValue(t types.Type, arg def.Value) (v Value, err error) {
	switch arg.ValueType() {
	case def.ValueSingle:
		v = &ConstantValue{
			value: arg.GetSingleValue(),
			typ:   t,
		}
		return
	case def.ValueService:
		v = &ServiceReference{
			ServiceName: arg.GetService(),
			Type:        t,
		}
		return
	case def.ValueSlice:
		v = &SliceValue{
			Value: arg,
			Type:  t,
		}
		return
	case def.ValueStruct:
		v = &StructValue{
			Value: arg,
			Type:  t,
		}
		return
	default:
		return nil, fmt.Errorf("Value type %s was not recognized", arg.ValueType().String())
	}
}

// Value will generate a value definition or use
type Value interface {
	Build(*ContainerGenerator)
	NeedsVariable() bool
	GenerateVariable(varName string) string
	GenerateUse() string
}

type ConstantValue struct {
	value         string
	typ           types.Type
	needsVariable bool
	varName       string
}

func (v *ConstantValue) Build(c *ContainerGenerator) {
	_, ok := v.typ.(*types.Pointer)
	v.needsVariable = ok
}

func (v *ConstantValue) NeedsVariable() bool {
	return v.needsVariable
}

func constantValueToString(v string, kind types.BasicKind) string {
	if kind == types.String {
		return strconv.Quote(v)

	}
	return v
}

func (v *ConstantValue) GenerateVariable(varName string) string {
	if !v.NeedsVariable() {
		return ""
	}

	typ := v.typ.(*types.Pointer).Elem().(*types.Basic)
	v.varName = varName
	return fmt.Sprintf("%s := %s", varName, constantValueToString(v.value, typ.Kind()))
}

func (v *ConstantValue) GenerateUse() string {
	if !v.NeedsVariable() {
		typ := v.typ.(*types.Basic)
		return constantValueToString(v.value, typ.Kind())
	}
	return fmt.Sprintf("&%s", v.varName)
}

type ServiceReference struct {
	Type        types.Type
	ServiceName string
}

func (v *ServiceReference) Build(c *ContainerGenerator) {
}

func (v *ServiceReference) NeedsVariable() bool {
	return false
}

func (v *ServiceReference) GenerateVariable(varName string) string {
	return ""
}

func (v *ServiceReference) GenerateUse() string {
	return fmt.Sprintf("c.Get%s()", v.ServiceName)
}

type SliceValue struct {
	Value         def.Value
	Type          types.Type
	values        []Value
	needsVariable bool
	pointer       bool
	sliceType     SliceType
	varName       string
}

func (v *SliceValue) Build(c *ContainerGenerator) {
	var t *types.Slice
	_, ok := v.Type.(*types.Pointer)
	v.pointer = ok
	v.needsVariable = ok

	if v.pointer {
		t = v.Type.(*types.Pointer).Elem().(*types.Slice)
	} else {
		t = v.Type.(*types.Slice)
	}

	v.sliceType = c.RegisterType(t).(SliceType)

	values := v.Value.GetSlice()
	v.values = make([]Value, len(values))
	for i, d := range values {
		v.values[i], _ = c.createValue(t.Elem(), d)
		v.values[i].Build(c)
		v.needsVariable = v.needsVariable || v.values[i].NeedsVariable()
	}
}

func (v *SliceValue) NeedsVariable() bool {
	return v.needsVariable
}

func (v *SliceValue) GenerateVariable(varName string) string {
	v.varName = varName
	b := &bytes.Buffer{}
	for i, value := range v.values {
		b.WriteString(value.GenerateVariable(fmt.Sprintf("%s_%d", varName, i)))
		b.WriteRune('\n')
	}

	if v.pointer {
		b.WriteString(varName)
		b.WriteString(" := ")
		sliceValueToUseString(v, b)
	}

	return b.String()
}

func sliceValueToUseString(v *SliceValue, b *bytes.Buffer) {
	b.WriteString(v.sliceType.String())
	b.WriteString("{\n")
	for _, value := range v.values {
		b.WriteString(value.GenerateUse())
		b.WriteString(",\n")
	}
	b.WriteRune('}')
}

func (v *SliceValue) GenerateUse() string {
	if v.pointer {
		return fmt.Sprintf("&%s", v.varName)
	}
	b := &bytes.Buffer{}
	sliceValueToUseString(v, b)
	return b.String()
}

type StructValue struct {
	Value def.Value
	Type  types.Type
}

func (v *StructValue) Build(c *ContainerGenerator) {
}

func (v *StructValue) NeedsVariable() bool {
	return false
}

func (v *StructValue) GenerateVariable(varName string) string {
	return ""
}

func (v *StructValue) GenerateUse() string {
	return fmt.Sprintf("%#v", v.Value)
}
