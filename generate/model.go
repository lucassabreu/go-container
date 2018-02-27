package generate

import (
	"fmt"
	"io"
)

const DURATION = "Duration"
const DURATIONSLIDE = "DurationSlide"
const FLOAT64 = "Float64"
const INT = "Int"
const INT64 = "Int64"
const STRING = "String"
const STRINGMAP = "StringMap"
const STRINGMAPSTRING = "StringMapString"
const STRINGMAPSTRINGSLICE = "StringMapStringSlice"
const STRINGSLICE = "StringSlice"
const TIME = "Time"

// ContainerFormatter generate golang code for the a Contaienr based on commands
type ContainerFormatter struct {
	ContainerPackage string
	ContainerName    string

	Packages map[string]string
	Services map[string]service
}

// GetPackageAlias return the alias for the package
func (cf *ContainerFormatter) GetPackageAlias(pkg string) *string {
	if alias, ok := cf.Packages[pkg]; ok {
		return &alias
	}

	return nil
}

// AddPackageAlias to the container
func (cf *ContainerFormatter) AddPackageAlias(pkg string, alias string) {
	cf.Packages[pkg] = alias
}

type value interface {
	Generate(c ContainerFormatter, castTo string, w io.Writer) error
}

type serviceValue struct {
	serviceName string
}

func (s serviceValue) Generate(c ContainerFormatter, castTo string, w io.Writer) error {
	_, err := w.Write([]byte(fmt.Sprintf("%s(c.Get%s())", castTo, s.serviceName)))
	return err
}

type parameterValue struct {
	parameterName string
}

func (s parameterValue) Generate(c ContainerFormatter, castTo string, w io.Writer) error {
	switch castTo {
	case DURATION:
	case DURATIONSLIDE:
	case FLOAT64:
	case INT:
	case INT64:
	case STRING:
	case STRINGMAP:
	case STRINGMAPSTRING:
	case STRINGMAPSTRINGSLICE:
	case STRINGSLICE:
	case TIME:
		_, err := w.Write([]byte(fmt.Sprintf("c.GetParametersBag().Get%s(\"%s\"))", castTo, s.parameterName)))
		return err
	}

	_, err := w.Write([]byte(fmt.Sprintf("%s(c.GetParametersBag().Get(\"%s\")))", castTo, s.parameterName)))
	return err
}

type staticValue struct {
	value interface{}
}

func (s staticValue) Generate(c ContainerFormatter, castTo string, w io.Writer) error {
	switch castTo {
	case INT, INT64:
		_, err := w.Write([]byte(fmt.Sprintf("%d", s.value)))
		return err
	case FLOAT64:
		_, err := w.Write([]byte(fmt.Sprintf("%f", s.value)))
		return err
	case STRING:
		_, err := w.Write([]byte(fmt.Sprintf("\"%s\"", s.value)))
		return err
	case STRINGSLICE:
		if _, err := w.Write([]byte("[...]string{\n")); err != nil {
			return err
		}

		strings := s.value.([]string)
		for _, value := range strings {
			if _, err := w.Write([]byte(fmt.Sprintf("\"%s\",", value))); err != nil {
				return err
			}
		}

		if _, err := w.Write([]byte("}")); err != nil {
			return err
		}
		return nil
	case STRINGMAPSTRING:
		if _, err := w.Write([]byte("map[string]string{\n")); err != nil {
			return err
		}

		strings := s.value.(map[string]string)
		for key, value := range strings {
			if _, err := w.Write([]byte(fmt.Sprintf("\"%s\":\"%s\",\n", key, value))); err != nil {
				return err
			}
		}

		if _, err := w.Write([]byte("}")); err != nil {
			return err
		}
		return nil
	case STRINGMAPSTRINGSLICE:
		if _, err := w.Write([]byte("map[string][]string{\n")); err != nil {
			return err
		}

		strings := s.value.(map[string][]string)
		for key, value := range strings {
			if _, err := w.Write([]byte(fmt.Sprintf("\"%s\": []string {", key))); err != nil {
				return err
			}

			for _, str := range value {
				if _, err := w.Write([]byte(fmt.Sprintf("\"%s\",\n", str))); err != nil {
					return err
				}
			}

			if _, err := w.Write([]byte("},\n")); err != nil {
				return err
			}
		}

		if _, err := w.Write([]byte("}")); err != nil {
			return err
		}
		return nil
	case DURATION:
		timePackage := c.GetPackageAlias("time")
		if timePackage == nil {
			s := "time"
			timePackage = &s
			c.AddPackageAlias(*timePackage, *timePackage)
		}

		_, err := w.Write([]byte(fmt.Sprintf("%s.ParseDuration(\"%s\")", *timePackage, s.value.(string))))
		return err
	case DURATIONSLIDE:
		timePackage := c.GetPackageAlias("time")
		if timePackage == nil {
			s := "time"
			timePackage = &s
			c.AddPackageAlias(*timePackage, *timePackage)
		}

		if _, err := w.Write([]byte(fmt.Sprintf("[...]%s.Duration{\n", *timePackage))); err != nil {
			return err
		}

		strings := s.value.([]string)
		for _, value := range strings {
			if _, err := w.Write([]byte(fmt.Sprintf("%s.Parse(%s.RFC3339, \"%s\"),", *timePackage, *timePackage, value))); err != nil {
				return err
			}
		}

		if _, err := w.Write([]byte("}")); err != nil {
			return err
		}
		return nil
	case TIME:
		_, err := w.Write([]byte(fmt.Sprintf("c.GetParametersBag().Get%s(\"%s\"))", castTo, s.value)))
		return err
	}

	return fmt.Errorf(
		"Static values should be one of this types: %v",
		[...]string{DURATION, DURATIONSLIDE, FLOAT64, INT, INT64, STRING, STRINGMAP, STRINGMAPSTRING, STRINGMAPSTRINGSLICE, STRINGSLICE, TIME})
}

type service interface {
	Name() string
	Generate(c ContainerFormatter, w io.Writer) error
}

type funcFactoryService struct {
	name        string
	factoryName string
	arguments   []value
}

type methodFactoryService struct {
	name        string
	serviceName string
	arguments   []value
}

type initStructService struct {
	name       string
	structName string
	values     map[string]value
}
