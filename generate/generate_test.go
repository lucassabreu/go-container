package generate_test

import (
	"testing"

	"github.com/lucassabreu/go-container/def"

	"github.com/lucassabreu/go-container/generate"
	yaml "gopkg.in/yaml.v2"
)

var testCircularReferenteTable = map[string]string{
	"Service Dependency not found": "services:\n Dependent:\n  factory: func\n  arguments: [@Dependency] } }",
}

func TestCircularReference(t *testing.T) {
	var cDef def.Container
	for expectedErr, yamlStr := range testCircularReferenteTable {
		yaml.Unmarshal([]byte(yamlStr), &cDef)
		t.Errorf("%#v", len(cDef.Services))
		_, err := generate.Generate(cDef)

		if err == nil {
			t.Errorf("expected error %s, got no error", expectedErr)
			continue
		}

		if err.Error() != expectedErr {
			t.Errorf("expected error %s, got %s", expectedErr, err.Error())
			continue
		}
	}
}
