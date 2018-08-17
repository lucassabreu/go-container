package generate

import (
	"fmt"
	"go/types"
	"strings"
)

// Type represents a type with its Package (if it have one)
type Type struct {
	typ     types.Type
	pointer bool
	name    string
	pkg     *Package
}

func (t Type) String() string {
	if t.pkg == nil {
		return t.name
	}
	s := fmt.Sprintf("%s.%s", t.pkg.UniqueName(), t.name)
	if t.pointer {
		return "*" + s
	}
	return s
}

// AsPointer formats the type to a pointer type declaration
func (t Type) AsPointer() string {
	if t.pointer {
		return t.String()
	}
	return "*" + t.String()
}

// RegisterType will create a new Type to be used in the generator,
// if the types.Type was aready registered, than the previews one
// will be returned
func (cg ContainerGenerator) RegisterType(t types.Type) Type {
	if typ, ok := cg.types[t]; ok {
		return typ
	}

	pointer := false
	fullName := t.String()
	if strings.HasPrefix(fullName, "*") {
		pointer = true
		fullName = fullName[1:]
	}

	path, name := breakIntoPackageAndDef(fullName)
	var pkg *Package
	if len(path) > 0 {
		cg.RegisterPackage(path, nil)
		pkg = cg.GetPackageByPath(path)
	}

	return Type{
		typ:     t,
		name:    name,
		pkg:     pkg,
		pointer: pointer,
	}
}
