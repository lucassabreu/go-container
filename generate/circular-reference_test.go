package generate_test

import (
	"testing"

	"github.com/lucassabreu/go-container/def"
	"github.com/lucassabreu/go-container/generate"
	"github.com/stretchr/testify/require"
)

func TestServiceReferences(t *testing.T) {
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
							"OtherDependency": def.NewServiceValue("OtherDependency"),
						},
					),
				},
			},
		},
		"factory_uses_struct_that_uses_factory": testCase{
			Err: "There is a circular reference for @(Factory|Struct) -> @(Struct|Factory) -> @(Factory|Struct)",
			Container: def.Container{
				Services: map[string]def.Service{
					"Factory": def.NewFactoryService("test.NewService", def.NewServiceValue("Struct")),
					"Struct": def.NewInitializationService("test.Service", map[string]def.Value{
						"Service": def.NewServiceValue("Factory"),
					}),
				},
			},
		},
		"factory1_that_uses_factory2_that_uses_factory1": testCase{
			Err: "There is a circular reference for @Service\\w -> @Service\\w -> @Service\\w",
			Container: def.Container{
				Services: map[string]def.Service{
					"Service1": def.NewFactoryService("test.NewService", def.NewServiceValue("Service2")),
					"Service2": def.NewFactoryService("test.NewService", def.NewServiceValue("Service1")),
				},
			},
		},
		"deep_circular": testCase{
			Err: "There is a circular reference for @\\w* -> @\\w* -> @\\w* -> @\\w*",
			Container: def.Container{
				Services: map[string]def.Service{
					"Factory":  def.NewFactoryService("test.NewService", def.NewServiceValue("Struct")),
					"Struct2":  def.NewInitializationService("test.Service", nil),
					"Factory2": def.NewFactoryService("test.NewService"),
					"MiddleOne": def.NewFactoryService("test.NewService", def.NewStructValue(map[string]def.Value{
						"field": def.NewServiceValue("Factory"),
					})),
					"Struct": def.NewInitializationService("test.Service", map[string]def.Value{
						"Services": def.NewSliceValue([]def.Value{
							def.NewServiceValue("Struct2"),
							def.NewServiceValue("Factory2"),
							def.NewServiceValue("MiddleOne"),
						}),
					}),
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := generate.CheckCircularReference(test.Container)
			require.Regexp(t, test.Err, err)
		})
	}
}
