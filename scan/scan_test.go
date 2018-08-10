package scan_test

import (
	"strings"
	"testing"

	"github.com/lucassabreu/go-container/scan"
	"github.com/stretchr/testify/require"
)

func TestImportPackage(t *testing.T) {
	pkgName := "github.com/lucassabreu/go-container/examples/test"
	pkg, err := scan.ImportPackage(pkgName)
	if err != nil {
		t.Errorf("Should not fail, error received: %v", err)
	}

	require.Equal(t, pkg.Name, "example")
	require.Equal(t, pkg.ImportPath, pkgName)

	require.Equal(t, 5, len(pkg.Funcs))

	if _, ok := pkg.Funcs["NewIDo"]; !ok {
		t.Fatalf("Should have found 'NewIDo' exported func")
	}

	if _, ok := pkg.Funcs["NewTheyDo"]; !ok {
		t.Fatalf("Should have found 'NewTheyDo' exported func")
	}

	require.Equal(t, 3, len(pkg.Structs))

	if _, ok := pkg.Structs["TheyDo"]; !ok {
		t.Fatalf("Should have found 'TheyDo' exported struct")
	}

	f := pkg.Funcs["NewTheyDo"]
	require.Equal(t, "NewTheyDo", f.Name)

	require.Equal(t, 1, len(f.Params))
	typ := f.Params[0]
	require.Equal(t, "func(string)", typ.String())

	require.Equal(t, 1, len(f.Results))
	typ = f.Results[0]
	require.Equal(t, pkgName+".TheyDo", typ.String())

	f = pkg.Funcs["NewIDo"]
	require.Equal(t, "NewIDo", f.Name)

	require.Equal(t, 0, len(f.Params))

	require.Equal(t, 1, len(f.Results))
	typ = f.Results[0]
	require.Equal(t, pkgName+".Doer", typ.String())

	require.Equal(
		t,
		`Package: example (github.com/lucassabreu/go-container/examples/test)
	Funcs:
		NewDoALot([]github.com/lucassabreu/go-container/examples/test.Doer) (github.com/lucassabreu/go-container/examples/test.doALot)
		NewIDo() (github.com/lucassabreu/go-container/examples/test.Doer)
		NewJustDo(string) (github.com/lucassabreu/go-container/examples/test.JustDo)
		NewSomethingDo(github.com/lucassabreu/go-container/examples/test.Doer) (github.com/lucassabreu/go-container/examples/test.SomethingDo)
		NewTheyDo(func(string)) (github.com/lucassabreu/go-container/examples/test.TheyDo)

	Structs:
		JustDo{
			That string,
		}
		SomethingDo{
			Something github.com/lucassabreu/go-container/examples/test.Doer,
		}
		TheyDo{
			ToDo func(string),
		}
`,
		pkg.String(),
	)
}

func TestImportPackageShouldFailWhenPackageNotExists(t *testing.T) {
	pkgName := "github.com/lucassabreu/go-container/examples/test1"
	_, err := scan.ImportPackage(pkgName)
	if err == nil {
		t.Errorf("Should fail, no error received")
	}

	if !strings.HasPrefix(err.Error(), "cannot find package") {
		t.Errorf("Should've not found the package, got: %v", err)
	}
}
