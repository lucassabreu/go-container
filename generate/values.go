package generate

import (
	"fmt"
	"go/types"
	"strconv"
)

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

type slideValue struct {
	values []Value
}

func (slideValue) Generate(ContainerGenerator) string {
	panic("not implemented")
}

func (slideValue) NeedsPointer() bool {
	panic("not implemented")
}
