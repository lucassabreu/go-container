package generate_test

import (
	"log"
	"testing"

	"github.com/lucassabreu/go-container/def"
	"github.com/lucassabreu/go-container/generate"
	"github.com/stretchr/testify/require"
)

func TestCompile(t *testing.T) {

	c := def.Container{
		Name: "SpecialContainer",
		Docs: "SpecialContainer is a container, but special",
		Packages: []def.Package{
			def.NewPackage("github.com/lucassabreu/go-container/examples/test"),
			def.NewPackageWithAlias("github.com/lucassabreu/go-container/examples/basicapp", "ex"),
		},
		Services: map[string]def.Service{
			"IDo":      def.NewFactoryService("test.NewIDo"),
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
			"ToDo": def.NewFactoryService(
				"test.NewToDo",
				def.NewStructValue(map[string]def.Value{
					"first":  def.NewSingleValue("wake up"),
					"second": def.NewSingleValue("drink coffe"),
					"third":  def.NewSingleValue("smile"),
				}),
			),
		},
	}

	cg, err := generate.NewContainerGenerator(c)

	require.Nil(t, err)

	err = cg.Compile()
	if err != nil {
		log.Printf(err.Error())
		require.Nil(t, err)
	}
	require.Equal(t, "<>", cg.String())
}
