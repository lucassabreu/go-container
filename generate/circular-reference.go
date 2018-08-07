package generate

import (
	"fmt"
	"strings"

	"github.com/lucassabreu/go-container/def"
)

func hasCircularReference(cDef def.Container) error {
	for name := range cDef.Services {
		if err := lookForItself(name, cDef); err != nil {
			return err
		}
	}

	return nil
}

func lookForItself(service string, cDef def.Container, referenced ...string) error {
	if contains(referenced, service) {
		return fmt.Errorf(
			"There is a circular reference for @%s -> @%s",
			strings.Join(referenced, " -> @"),
			service,
		)
	}

	referenced = append(referenced, service)

	sDef, ok := cDef.Services[service]
	if !ok {
		return fmt.Errorf("Service %s not found", service)
	}

	if sDef.IsByFactory() {
		if sDef.Arguments == nil {
			return nil
		}

		return lookupSlice(sDef.Arguments, cDef, referenced...)
	}

	if sDef.Fields == nil {
		return nil
	}

	return lookupStruct(sDef.Fields, cDef, referenced...)
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func testValue(vDef def.Value, cDef def.Container, referenced ...string) error {
	switch vDef.ValueType() {
	case def.ValueService:
		if err := lookForItself(vDef.GetService(), cDef, referenced...); err != nil {
			return err
		}
	case def.ValueSlice:
		if err := lookupSlice(vDef.GetSlice(), cDef, referenced...); err != nil {
			return err
		}
	case def.ValueStruct:
		if err := lookupStruct(vDef.GetStruct(), cDef, referenced...); err != nil {
			return err
		}
	}
	return nil
}

func lookupSlice(values []def.Value, cDef def.Container, referenced ...string) error {
	for _, vDef := range values {
		if err := testValue(vDef, cDef, referenced...); err != nil {
			return err
		}
	}
	return nil
}

func lookupStruct(values map[string]def.Value, cDef def.Container, referenced ...string) error {
	for _, vDef := range values {
		if err := testValue(vDef, cDef, referenced...); err != nil {
			return err
		}
	}
	return nil
}
