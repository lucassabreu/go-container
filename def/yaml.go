package def

import (
	"fmt"
	"strings"
)

// UnmarshalYAML for the Package
func (p *Package) UnmarshalYAML(unmarshal func(interface{}) error) error {
	defaultErr := fmt.Errorf("the package should a single string or a map of strings with a element")

	var aliasedPackage map[string]string
	if err := unmarshal(&aliasedPackage); err == nil {
		if len(aliasedPackage) != 1 {
			return defaultErr
		}
		for packageStr, alias := range aliasedPackage {
			p.Package = packageStr
			p.Alias = &alias
			return nil
		}
	}

	var packageName string
	if err := unmarshal(&packageName); err == nil {
		p.Package = packageName
		return nil
	}

	return defaultErr
}

// UnmarshalYAML for the Value
func (v *Value) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var structValue map[string]Value
	if err := unmarshal(&structValue); err == nil {
		*v = NewStructValue(structValue)
		return nil
	}

	var sliceValue []Value
	if err := unmarshal(&sliceValue); err == nil {
		*v = NewSliceValue(sliceValue)
		return nil
	}

	var value string
	if err := unmarshal(&value); err == nil {
		if strings.HasPrefix(value, "@") {
			*v = NewServiceValue(value[1:])
			return nil
		}

		*v = NewSingleValue(value)
		return nil
	}

	return fmt.Errorf("the value type was not recognized")
}
