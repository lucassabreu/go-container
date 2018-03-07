package generate

import (
	"go/types"
	"strconv"
)

type constValueDef struct {
	value string
	typ   types.Type
}

func (c constValueDef) Generate(ContainerGenerator) string {
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

func (c constValueDef) NeedsPointer() bool {
	_, ok := c.typ.(*types.Pointer)
	return ok
}

type slideValueDef struct {
	values []valueDef
}

func (slideValueDef) Generate(ContainerGenerator) string {
	panic("not implemented")
}

func (slideValueDef) NeedsPointer() bool {
	panic("not implemented")
}
