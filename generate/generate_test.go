package generate_test

import (
	"regexp"
	"testing"

	"github.com/lucassabreu/go-container/def"

	"github.com/lucassabreu/go-container/generate"
	yaml "gopkg.in/yaml.v2"
)

func TestCircularReference(t *testing.T) {
	var tests []struct {
		Container def.Container
		Err       *string
	}

	bytes := []byte(`
- err: "Service Dependency not found"
  container:
    services:
      Dependent:
        factory: test.NewDependent
        arguments:
          - "@Dependency"

- err: "Service OtherDependency not found"
  container:
    services:
      Dependent:
        struct: test.Dependent
        fields:
          OtherDependency: "@OtherDependency"

- err: "There is a circular reference for @Service\\\\w -> @Service\\\\w -> @Service\\\\w"
  container:
    services:
      Service1:
        factory: test.NewService
        arguments:
          - "@Service2"
      Service2:
        factory: test.NewService
        arguments:
          - "@Service1"

- err: "There is a circular reference for @\\\\w* -> @\\\\w* -> @\\\\w* -> @\\\\w*"
  container:
    services:
      Factory:
        factory: test.NewService
        arguments:
          - Service: "@Struct"
      Struct2:
        struct: test.Service
      Factory2:
        factory: test.Service
      MiddleOne:
        factory: test.NewService
        arguments:
          - "@Factory"
      Struct:
        struct: test.Service
        fields:
          Services:
            - "@Struct2"
            - "@Factory2"
            - "@MiddleOne"

- container:
  services:
    Factory:
      factory: test.NewService
      arguments: [ "@Struct" ]
    Struct:
      struct: test.Service
      fields:
        - Service: "@Factory"
`)

	err := yaml.Unmarshal(bytes, &tests)
	if err != nil {
		t.Error(err, ":", string(bytes))
	}

	for _, test := range tests {
		_, err := generate.Generate(test.Container)

		if test.Err == nil && err == nil {
			continue
		}

		if test.Err == nil && err != nil {
			t.Errorf("expected no error, got '%s' error", err.Error())
			continue
		}

		if err == nil {
			t.Errorf("expected error '%s', got no error", *test.Err)
			continue
		}

		if *test.Err != err.Error() && regexp.MustCompile(*test.Err).MatchString(err.Error()) {
			t.Errorf("expected error '%s', got '%s'", *test.Err, err.Error())
			continue
		}
	}
}
