package generate_test

import (
	"regexp"
	"testing"

	"github.com/lucassabreu/go-container/def"

	"github.com/lucassabreu/go-container/generate"
)

var bytes = []byte(`

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

`)

func TestCircularReference(t *testing.T) {
	type testCase struct {
		Container def.Container
		Err       string
	}

	tests := map[string]testCase{
		"dependency_not_found_factory": testCase{
			Err: "Service Dependency not found",
			Container: def.Container{
				Services: map[string]def.Service{
					"Dependent": def.NewFactoryService("test.NewDependent", def.NewServiceValue("Dependency")),
				},
			},
		},
		"dependency_not_found_struct": testCase{
			Err: "Service OtherDependency not found",
			Container: def.Container{
				Services: map[string]def.Service{
					"Dependent": def.NewInitializationService(
						"test.Dependent",
						map[string]def.Value{
							"OtherDependency": def.NewServiceValue("@OtherDependency"),
						},
					),
				},
			},
		},
		"no_problem_with_struct_and_factory": testCase {
			Container: def.Container{
				Services: map[string]def.Service {
					"Factory": def.NewFactoryService("test.NewService", def.NewServiceValue("Struct")
				},
			},
		},
- container:
  services:
    Factory:
      factory: test.NewService
      arguments: [ "@Struct" ]
    Struct:
      struct: test.Service
      fields:
        - Service: "@Factory"
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := generate.Generate(test.Container)

			if len(test.Err) == 0 && err == nil {
				return
			}

			if len(test.Err) == 0 && err != nil {
				t.Errorf("expected no error, got '%s' error", err.Error())
				return
			}

			if err == nil {
				t.Errorf("expected error '%s', got no error", test.Err)
				return
			}

			if test.Err != err.Error() && regexp.MustCompile(test.Err).MatchString(err.Error()) {
				t.Errorf("expected error '%s', got '%s'", test.Err, err.Error())
				return
			}
		})
	}
}
