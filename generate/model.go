package generate

import (
	"fmt"
	"io"
)

const durationStr = "Duration"
const durationSlideStr = "DurationSlide"
const float64Str = "Float64"
const intStr = "Int"
const int64Str = "Int64"
const stringStr = "String"
const stringMapStr = "StringMap"
const stringMapStringStr = "StringMapString"
const stringMapStringSliceStr = "StringMapStringSlice"
const stringSliceStr = "StringSlice"
const timeStr = "Time"

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
	case durationStr:
	case durationSlideStr:
	case float64Str:
	case intStr:
	case int64Str:
	case stringStr:
	case stringMapStr:
	case stringMapStringStr:
	case stringMapStringSliceStr:
	case stringSliceStr:
	case timeStr:
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
	case intStr, int64Str:
		_, err := w.Write([]byte(fmt.Sprintf("%d", s.value)))
		return err
	case float64Str:
		_, err := w.Write([]byte(fmt.Sprintf("%f", s.value)))
		return err
	case stringStr:
		_, err := w.Write([]byte(fmt.Sprintf("\"%s\"", s.value)))
		return err
	case stringSliceStr:
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
	case stringMapStringStr:
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
	case stringMapStringSliceStr:
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
	case durationStr:
		timePackage := c.GetPackageAlias("time")
		if timePackage == nil {
			s := "time"
			timePackage = &s
			c.AddPackageAlias(*timePackage, *timePackage)
		}

		_, err := w.Write([]byte(fmt.Sprintf("%s.ParseDuration(\"%s\")", *timePackage, s.value.(string))))
		return err
	case durationSlideStr:
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
	case timeStr:
		_, err := w.Write([]byte(fmt.Sprintf("c.GetParametersBag().Get%s(\"%s\"))", castTo, s.value)))
		return err
	}

	return fmt.Errorf(
		"Static values should be one of this types: %v",
		[...]string{durationStr, durationSlideStr, float64Str, intStr, int64Str, stringStr, stringMapStr, stringMapStringStr, stringMapStringSliceStr, stringSliceStr, timeStr})
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
