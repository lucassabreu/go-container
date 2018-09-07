package generate

import (
	"bytes"
	"fmt"
	"go/types"
)

// Type represents a type with its Package (if it have one)
type Type interface {
	Type() types.Type
	String() string
	StringAsPointer() string
}

// BasicType represents a basic type
type BasicType struct {
	InnerType *types.Basic
}

func (t BasicType) String() string {
	return t.InnerType.String()
}

// StringAsPointer strings the object to a pointer
func (t BasicType) StringAsPointer() string {
	return "*" + t.String()
}

// Type returns the types.Type from the Type
func (t BasicType) Type() types.Type {
	return t.InnerType
}

// PointerType represents a type that is a pointer
type PointerType struct {
	InnerType Type
}

func (t PointerType) String() string {
	return "*" + t.InnerType.String()
}

// StringAsPointer strings the object to a pointer
func (t PointerType) StringAsPointer() string {
	return t.String()
}

// Type returns the types.Type from the Type
func (t PointerType) Type() types.Type {
	return t.InnerType.Type()
}

// NamedType represents a struct type
type NamedType struct {
	InnerType types.Type
	Name      string
	Package   *Package
}

func (t NamedType) String() string {
	return fmt.Sprintf("%s.%s", t.Package.UniqueName(), t.Name)
}

// StringAsPointer strings the object to a pointer
func (t NamedType) StringAsPointer() string {
	return "*" + t.String()
}

// Type returns the types.Type from the Type
func (t NamedType) Type() types.Type {
	return t.InnerType
}

// GetStruct get the underlying struct, if it exists
func (t NamedType) GetStruct() *types.Struct {
	return t.Type().(*types.Struct)
}

// InterfaceType represents a interface type
type InterfaceType struct {
	InnerType *types.Interface
	Name      string
	Package   *Package
}

func (t InterfaceType) String() string {
	return fmt.Sprintf("%s.%s", t.Package.UniqueName(), t.Name)
}

// StringAsPointer strings the object to a pointer
func (t InterfaceType) StringAsPointer() string {
	return t.String()
}

// Type returns the types.Type from the Type
func (t InterfaceType) Type() types.Type {
	return t.InnerType
}

// ArrayType represents a array type
type ArrayType struct {
	InnerType *types.Array
	ElemType  Type
}

func (t ArrayType) String() string {
	return fmt.Sprintf("[%d]%s", t.InnerType.Len(), t.ElemType.String())
}

// StringAsPointer strings the object to a pointer
func (t ArrayType) StringAsPointer() string {
	return "*" + t.String()
}

// Type returns the types.Type from the Type
func (t ArrayType) Type() types.Type {
	return t.InnerType
}

// SliceType represents a slice type
type SliceType struct {
	InnerType *types.Slice
	ElemType  Type
}

func (t SliceType) String() string {
	return fmt.Sprintf("[]%s", t.ElemType.String())
}

// StringAsPointer strings the object to a pointer
func (t SliceType) StringAsPointer() string {
	return "*" + t.String()
}

// Type returns the types.Type from the Type
func (t SliceType) Type() types.Type {
	return t.InnerType
}

// ChanType represents a channel type
type ChanType struct {
	InnerType *types.Chan
	ElemType  Type
}

func (t ChanType) String() string {
	s := bytes.Buffer{}
	if t.InnerType.Dir() == types.SendOnly {
		s.WriteString("<-")
	}
	s.WriteString("chan")
	if t.InnerType.Dir() == types.RecvOnly {
		s.WriteString("<-")
	}
	s.WriteRune(' ')
	s.WriteString(t.ElemType.String())
	return s.String()
}

// StringAsPointer strings the object to a pointer
func (t ChanType) StringAsPointer() string {
	return "*" + t.String()
}

// Type returns the types.Type from the Type
func (t ChanType) Type() types.Type {
	return t.InnerType
}

// MapType represents a map type
type MapType struct {
	InnerType *types.Map
	ElemType  Type
	KeyType   Type
}

func (t MapType) String() string {
	return fmt.Sprintf("map[%s]%s", t.KeyType.String(), t.ElemType.String())
}

// StringAsPointer strings the object to a pointer
func (t MapType) StringAsPointer() string {
	return "*" + t.String()
}

// Type returns the types.Type from the Type
func (t MapType) Type() types.Type {
	return t.InnerType
}

// RegisterType will create a new Type to be used in the generator,
// if the types.Type was aready registered, than the previews one
// will be returned
func (cg ContainerGenerator) RegisterType(t types.Type) Type {
	if typ, ok := cg.types[t]; ok {
		return typ
	}

	var typ Type
	switch t := t.(type) {
	case *types.Pointer:
		typ = PointerType{
			InnerType: cg.RegisterType(t.Elem()),
		}
	case *types.Basic:
		typ = BasicType{InnerType: t}
	case *types.Array:
		typ = ArrayType{
			ElemType:  cg.RegisterType(t.Elem()),
			InnerType: t,
		}
	case *types.Slice:
		typ = SliceType{
			ElemType:  cg.RegisterType(t.Elem()),
			InnerType: t,
		}
	case *types.Map:
		typ = MapType{
			ElemType:  cg.RegisterType(t.Elem()),
			KeyType:   cg.RegisterType(t.Key()),
			InnerType: t,
		}
	case *types.Chan:
		typ = ChanType{
			ElemType:  cg.RegisterType(t.Elem()),
			InnerType: t,
		}
	case *types.Named:
		fullName := t.String()
		path, name := breakIntoPackageAndDef(fullName)
		cg.RegisterPackage(path, nil)
		pkg := cg.GetPackageByPath(path)

		switch u := t.Underlying().(type) {
		case *types.Interface:
			typ = InterfaceType{
				InnerType: u,
				Name:      name,
				Package:   pkg,
			}
		case *types.Struct:
			typ = NamedType{
				InnerType: u,
				Name:      name,
				Package:   pkg,
			}
		default:
			typ = NamedType{
				InnerType: t,
				Name:      name,
				Package:   pkg,
			}
		}
	}

	cg.types[t] = typ

	return typ
}
