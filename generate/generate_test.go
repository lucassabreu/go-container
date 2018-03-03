package generate_test

import (
	"io/ioutil"
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

	bytes, err := ioutil.ReadFile("./test-files/circular-reference.yml")
	if err != nil {
		t.Error(err)
	}

	err = yaml.Unmarshal(bytes, &tests)
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
