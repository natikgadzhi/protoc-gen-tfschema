package gen

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/gravitational/trace"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type SchemaData struct {
	PackageName string
	Messages    []MessageData
}

type MessageData struct {
	Name   string
	Fields []FieldData
}

type FieldData struct {
	Name     string
	Type     string
	Required bool
	Meta     string
}

var schemaTpl string = `
package {{ .PackageName }}

import (
	"github.com/hashicorp/terraform/helper/schema"
)

{{- range .Messages }}
func schema{{.Name}}() map[string]*schema.Schema {
	schema := map[string]*schema.Schema{
	{{- range .Fields }}
		"{{ .Name }}": {
			Type: schema.{{.Type}},
			{{- if .Required }}
			Required: true,
			{{- else }}
			Optional: true,
			{{- end }}
		},
	{{- end }}
	}
	return schema
}
{{- end }}

`

// schemaTemplate takes the data in and renders the go file
// containing terraform schema definitions for Messages specified in the data.
func schemaTemplate(data *SchemaData) (string, error) {
	log.Infof("generating template for %s", data.PackageName)
	schemaTemplate, err := template.New("schema").Parse(schemaTpl)

	if err != nil {
		return "", trace.Wrap(err)
	}

	log.Infof("writing template for %s", data.PackageName)
	var buf bytes.Buffer
	if err := schemaTemplate.Execute(&buf, data); err != nil {
		return "", err
	}

	generated := buf.String()
	return generated, nil
}

func kindToTerraform(kind protoreflect.Kind) string {

	kind2tf := map[protoreflect.Kind]string{
		protoreflect.BoolKind:   "TypeBool",
		protoreflect.StringKind: "TypeString",

		// "TYPE_STRING": "TypeString",
		// "TYPE_BYTES":  "TypeString",
		// "TYPE_ENUM":   "TypeString",

		// "TYPE_BOOL": "TypeBool",

		// "TYPE_INT64":  "TypeInt",
		// "TYPE_UINT32": "TypeInt",
	}

	tfType, ok := kind2tf[kind]
	if ok {
		return tfType
	}

	log.Warnf("Unspecified field type: %s", kind)
	return fmt.Sprintf("%v", kind)
}

func requiredProtoToSchema(t string) bool {
	switch t {
	case "LABEL_OPTIONAL":
		return false
	case "LABEL_REQUIRED":
		return true
	default:
		return true
	}
}
