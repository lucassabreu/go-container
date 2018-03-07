package generate

import (
	"go/types"
	"strconv"
)

type constValueGen struct {
	value string
	typ   types.Type
}

func (c constValueGen) Generate(ContainerGenerator) string {
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

func (c constValueGen) NeedsPointer() bool {
	_, ok := c.typ.(*types.Pointer)
	return ok
}

type slideValueGen struct {
	values []valueGen
}

func (slideValueGen) Generate(ContainerGenerator) string {
	panic("not implemented")
}

func (slideValueGen) NeedsPointer() bool {
	panic("not implemented")
}
