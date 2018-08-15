package generate

import (
	"bytes"
	"errors"
	"fmt"
	"text/template"
)

var tpl *template.Template

func init() {
	tpl = template.New("container")
	var err error
	tpl, err = tpl.Parse(`
package {{.PackageName}}

import (
	{{ range .Packages -}}
	{{ .UniqueName }} {{ .Path }}
	{{- end }}
)

// {{.ContainerDocs}}
type {{.ContainerName}} struct {
	parametersBag {{.GoContainerPackageName}}.ParametersBag

	{{ range .Services }}

	{{ end }}
}
`)
	if err != nil {
		panic(err)
	}
}

// IsCompiled indicates if the container was compiled
func (cg ContainerGenerator) IsCompiled() bool {
	return cg.buffer.Len() > 0
}

// Compile generates the container contentes and closes it for edition
func (cg ContainerGenerator) Compile() error {
	if cg.IsCompiled() {
		return nil
	}

	if len(cg.packages) == 0 {
		return errors.New("no package was registered")
	}

	goContainerPackageAlias := "container"
	for i := 0; contains(cg.importedPackageNames, goContainerPackageAlias); i++ {
		goContainerPackageAlias = fmt.Sprintf("container%d", i)
	}
	cg.RegisterPackage("github.com/lucassabreu/go-container", &goContainerPackageAlias)

	if len(cg.ContainerName) == 0 {
		cg.ContainerName = (DefaultContainerPackage)
	}

	if len(cg.ContainerDocs) == 0 {
		cg.ContainerDocs = DefaultContainerDocs
	}
	if len(cg.ContainerName) == 0 {
		cg.ContainerName = DefaultContainerName
	}

	b := bytes.Buffer{}
	err := tpl.Execute(&b, map[string]interface{}{
		"ContainerPackage":       cg.ContainerPackage,
		"ContainerDocs":          cg.ContainerDocs,
		"ContainerName":          cg.ContainerName,
		"GoCotainerPackageAlias": goContainerPackageAlias,
		"Packages":               cg.packages,
		"Services":               cg.services,
	})
	if err != nil {
		return err
	}

	cg.buffer = b
	return nil
}

func (cg ContainerGenerator) Read(p []byte) (int, error) {
	if err := cg.Compile(); err != nil {
		return 0, err
	}

	return cg.buffer.Read(p)
}

func (cg ContainerGenerator) String() string {
	if err := cg.Compile(); err != nil {
		return ""
	}

	return cg.buffer.String()
}
