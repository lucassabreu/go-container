package generate

import (
	"go/types"
	"strconv"
)

// Value will generate a value definition or use
type Value interface {
	Generate(ContainerGenerator) string
	NeedsPointer() bool
}

type constValue struct {
	value string
	typ   types.Type
}

func (c constValue) Generate(ContainerGenerator) string {
	var typ types.Type
	typ, ok := c.typ.(*types.Pointer)
	if !ok {
		typ = c.typ
	}

	b := typ.(*types.Basic)
	if b.Kind() == types.String {
		return strconv.Quote(c.value)

	}

	return c.value
}

func (c constValue) NeedsPointer() bool {
	_, ok := c.typ.(*types.Pointer)
	return ok
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
