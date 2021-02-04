package renderer

import (
	"bytes"
	"errors"
	"path"
	"runtime"
	"text/template"

	"github.com/nategadzhi/protoc-gen-tfschema/builder"
	"github.com/nategadzhi/protoc-gen-tfschema/config"
	"github.com/stoewer/go-strcase"
)

const (
	templateFilename = "schema.tpl"
	templatesDir     = "/renderer/tpl"
)

type viewData struct {
	Version       string
	ProtocVersion string
	PackageName   string
	Resources     *builder.ResourceMap
}

// Render renders file
func Render(resources *builder.ResourceMap, version string) (*bytes.Buffer, error) {
	var buf bytes.Buffer

	_, filename, _, ok := runtime.Caller(1)

	if !ok {
		return nil, errors.New("Can't get path to runtime file")
	}

	filepath := path.Join(path.Dir(filename), templatesDir, templateFilename)

	tpl, err := template.New("test").Funcs(helpers()).ParseFiles(filepath)

	if err != nil {
		return nil, err
	}

	data := viewData{
		PackageName:   *config.PackageName,
		ProtocVersion: config.ProtocVersion,
		Version:       version,
		Resources:     resources,
	}

	err = tpl.ExecuteTemplate(&buf, templateFilename, data)
	if err != nil {
		return nil, err
	}

	return &buf, nil
}

func helpers() template.FuncMap {
	return template.FuncMap{
		// IsSchema returns true if an element is a schema
		"IsSchema": func(i interface{}) bool {
			_, ok := i.(*builder.Schema)
			return ok
		},
		"SnakeCase": strcase.SnakeCase,
	}
}
