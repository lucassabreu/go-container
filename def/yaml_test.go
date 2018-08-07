package def_test

import (
	"testing"

	"github.com/lucassabreu/go-container/def"
	"github.com/stretchr/testify/require"
	yaml "gopkg.in/yaml.v2"
)

func TestYamlUnmarshaller(t *testing.T) {
	type testCase struct {
		expected def.Container
		yaml     string
	}

	tests := map[string]testCase{
		"simple": testCase{
			expected: def.Container{
				Packages: []def.Package{
					def.Package{Package: "github.com/lucassabreu/go-container/test"},
				},
				Services: map[string]def.Service{
					"IDo": def.NewFactoryService("test.IDo"),
					"JustDo": def.NewInitializationService(
						"test.JustDo",
						map[string]def.Value{
							"That": def.NewSingleValue("other thing"),
						},
					),
				},
			},
			yaml: `
packages:
  - github.com/lucassabreu/go-container/test

services:
  IDo:
    factory: test.IDo

  JustDo:
    struct: test.JustDo
    fields:
      That: "other thing"
`,
		},
	}

	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			converted := def.Container{}
			if err := yaml.Unmarshal([]byte(testCase.yaml), &converted); err != nil {
				t.Fatal(err)
			}

			require.Equal(t, testCase.expected, converted)
		})
	}
}
