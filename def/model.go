package def

// Container definitions
type Container struct {
	// Container struct metadata
	Name string
	Docs string

	// Packages that will be used and their aliases (when needed)
	Packages []Package

	// Services to be managed by the container
	Services map[string]Service
}

// Package represents a package in the Container
type Package struct {
	Package string
	Alias   *string
}

// Service represents a entry in the services config
type Service struct {
	Factory   *string
	Arguments []Value
	Struct    *string
	Fields    map[string]Value
}

func NewFactoryService(factoryFuncName string, args ...Value) Service {
	s := Service{
		Factory:   &factoryFuncName,
		Arguments: args,
	}
	return s
}

func NewInitializationService(structName string, fields map[string]Value) Service {
	return Service{
		Struct: &structName,
		Fields: fields,
	}
}

// IsByFactory when a factory function factory should be called
func (s Service) IsByFactory() bool {
	return s.Factory != nil
}

// IsByInitialization when the struct should be initialized
func (s Service) IsByInitialization() bool {
	return s.Struct != nil
}

// ValueType represents the type of the value
type ValueType int16

const (
	// ValueSingle is a value that is static and single valued
	ValueSingle ValueType = iota << 1
	// ValueStruct is a map or struct of values
	ValueStruct
	// ValueSlice is a slice of values
	ValueSlice
	// ValueService is a reference to a service
	ValueService
)

// Value represents a generic value (may be a reference to a service, a static value or static struct)
type Value struct {
	value     interface{}
	valueType ValueType
}

// ValueType represents the type of the value
func (v Value) ValueType() ValueType {
	return v.valueType
}

// GetSingleValue should be used to get the value of a ValueSingle type
func (v Value) GetSingleValue() string {
	return v.value.(string)
}

// GetStruct should be used to get the value of a ValueStruct type
func (v Value) GetStruct() map[string]Value {
	return v.value.(map[string]Value)
}

// GetSlice should be used to get the value of a ValueSlice type
func (v Value) GetSlice() []Value {
	return v.value.([]Value)
}

// GetService should be used to get the value of a ValueService type
func (v Value) GetService() string {
	return v.value.(string)
}

func NewSingleValue(v string) Value {
	return Value{
		valueType: ValueSingle,
		value:     v,
	}
}

func NewServiceValue(v string) Value {
	return Value{
		valueType: ValueService,
		value:     v,
	}
}

func NewSliceValue(v []Value) Value {
	return Value{
		valueType: ValueSlice,
		value:     v,
	}
}

func NewStructValue(v map[string]Value) Value {
	return Value{
		valueType: ValueStruct,
		value:     v,
	}
}
