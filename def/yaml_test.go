package def_test

import (
	"testing"

	"github.com/lucassabreu/go-container/def"
	"github.com/stretchr/testify/require"
	yaml "gopkg.in/yaml.v2"
)

func TestYamlUnmarshaller(t *testing.T) {

	expected := def.Container{
		Packages: []def.Package{
			def.Package{Package: "github.com/lucassabreu/go-container/test"},
		},
		Services: map[string]def.Service{
			"IDo":      def.NewFactoryService("test.IDo"),
			"JustDoIt": def.NewFactoryService("test.NewJustDo", def.NewSingleValue("it")),
			"JustDo": def.NewInitializationService(
				"test.JustDo",
				map[string]def.Value{
					"That": def.NewSingleValue("other thing"),
				},
			),
			"SomethingDo": def.NewFactoryService("test.NewSomethingDo", def.NewServiceValue("IDo")),
			"DoALot": def.NewFactoryService(
				"test.NewDoALot",
				def.NewSliceValue([]def.Value{
					def.NewServiceValue("IDo"),
					def.NewServiceValue("SomethingDo"),
				}),
			),
		},
	}
	yamlStr := `
packages:
  - github.com/lucassabreu/go-container/test

services:
  IDo:
    factory: test.IDo

  JustDoIt:
    factory: test.NewJustDo
    arguments:
      - "it"

  JustDo:
    struct: test.JustDo
    fields:
      That: "other thing"

  SomethingDo:
    factory: test.NewSomethingDo
    arguments: [ "@IDo" ]

  DoALot:
    factory: test.NewDoALot
    arguments: [ [ "@IDo", "@SomethingDo" ] ]
`

	converted := def.Container{}
	if err := yaml.Unmarshal([]byte(yamlStr), &converted); err != nil {
		t.Fatal(err)
	}

	require.Equal(t, expected, converted)
}
