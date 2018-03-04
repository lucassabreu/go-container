package scan_test

import (
	"strings"
	"testing"

	"github.com/lucassabreu/go-container/scan"
)

func TestImportPackage(t *testing.T) {
	pkgName := "github.com/lucassabreu/go-container/scan/test"
	pkg, err := scan.ImportPackage(pkgName)
	if err != nil {
		t.Errorf("Should not fail, error received: %v", err)
	}

	assertEqual(t, pkg.Name, "example")
	assertEqual(t, pkg.ImportPath, pkgName)

	assertEqual(t, 2, len(pkg.Funcs))

	if _, ok := pkg.Funcs["NewIDo"]; !ok {
		t.Fatalf("Should have found 'NewIDo' exported func")
	}

	if _, ok := pkg.Funcs["NewTheyDo"]; !ok {
		t.Fatalf("Should have found 'NewTheyDo' exported func")
	}

	assertEqual(t, 1, len(pkg.Structs))

	if _, ok := pkg.Structs["TheyDo"]; !ok {
		t.Fatalf("Should have found 'TheyDo' exported struct")
	}

	f := pkg.Funcs["NewTheyDo"]
	assertEqual(t, "NewTheyDo", f.Name)

	assertEqual(t, 1, len(f.Params))
	typ := f.Params[0]
	assertEqual(t, "func(string)", typ.String())

	assertEqual(t, 1, len(f.Results))
	typ = f.Results[0]
	assertEqual(t, pkgName+".TheyDo", typ.String())

	f = pkg.Funcs["NewIDo"]
	assertEqual(t, "NewIDo", f.Name)

	assertEqual(t, 0, len(f.Params))

	assertEqual(t, 1, len(f.Results))
	typ = f.Results[0]
	assertEqual(t, pkgName+".Doer", typ.String())

	assertEqual(
		t,
		`Package: example (github.com/lucassabreu/go-container/scan/test)
	Funcs:
		NewIDo() (github.com/lucassabreu/go-container/scan/test.Doer)
		NewTheyDo(func(string)) (github.com/lucassabreu/go-container/scan/test.TheyDo)

	Structs:
		TheyDo{
			ToDo func(string),
		}
`,
		pkg.String(),
	)
}

func TestImportPackageShouldFailWhenPackageNotExists(t *testing.T) {
	pkgName := "github.com/lucassabreu/go-container/scan/test1"
	_, err := scan.ImportPackage(pkgName)
	if err == nil {
		t.Errorf("Should fail, no error received")
	}

	if !strings.HasPrefix(err.Error(), "cannot find package") {
		t.Errorf("Should've not found the package, got: %v", err)
	}
}

func assertEqual(t *testing.T, expected interface{}, value interface{}) {
	if expected != value {
		t.Fatalf("expected %v, got %v", expected, value)
	}
}
