package renderer

import (
	"bytes"
	"log"
	"path"
	"runtime"
	"text/template"

	"github.com/nategadzhi/protoc-gen-tfschema/builder"
	"github.com/nategadzhi/protoc-gen-tfschema/config"
)

type viewData struct {
	Version       string
	ProtocVersion string
	PackageName   string
}

// Render renders file
func Render(resource *builder.ResourceMap, version string) *bytes.Buffer {
	var buf bytes.Buffer

	_, filename, _, ok := runtime.Caller(1)

	if !ok {
		log.Fatal("Can't get path to runtime file")
	}

	filepath := path.Join(path.Dir(filename), "/renderer/tpl/file.gotpl")

	tpl, err := template.ParseFiles(filepath)
	if err != nil {
		log.Fatal(err)
	}

	data := viewData{
		PackageName:   *config.PackageName,
		ProtocVersion: config.ProtocVersion,
		Version:       version,
	}

	err = tpl.ExecuteTemplate(&buf, "file.gotpl", data)
	if err != nil {
		log.Fatal(err)
	}

	return &buf
}
