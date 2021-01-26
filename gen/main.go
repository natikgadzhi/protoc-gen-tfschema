package gen

import (
	"fmt"
	"strings"

	"github.com/gravitational/trace"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// kind2tf is a mapping from protobuf primitive types to terraform types
var kind2tf = map[protoreflect.Kind]string{
	protoreflect.BoolKind: "TypeBool",

	protoreflect.StringKind: "TypeString",
	protoreflect.BytesKind:  "TypeString",
	protoreflect.EnumKind:   "TypeString",

	protoreflect.Int32Kind:   "TypeInt",
	protoreflect.Int64Kind:   "TypeInt",
	protoreflect.Fixed32Kind: "TypeInt",
	protoreflect.Fixed64Kind: "TypeInt",
	protoreflect.Uint32Kind:  "TypeInt",
	protoreflect.Uint64Kind:  "TypeInt",
	protoreflect.DoubleKind:  "TypeFloat",
}

// Generator is the Terraform Schema generator
type Generator struct {
	Plugin        *protogen.Plugin // Plugin
	packageName   string           // Package name
	types         []string         // Types to export
	protocVersion string           // protoc version
}

// NewGenerator returns new generator instance.
func NewGenerator(plugin *protogen.Plugin) *Generator {
	g := &Generator{Plugin: plugin}

	// TODO: Move to this file once it becomes shorter
	g.parseArgs()
	g.fillVersion()

	log.Infof("Command line options: %s", g.Plugin.Request.GetParameter())
	log.Infof("Package name: %s", g.packageName)
	log.Infof("Types: %s", g.types)
	log.Infof("Protoc version: %s", g.protocVersion)

	return g
}

const (
	defaultPackageName = "tfschema"
)

// parseOptions fills configuration from command line args
func (g *Generator) parseArgs() {
	var dict map[string]string = make(map[string]string)

	fields := strings.Fields(*g.Plugin.Request.Parameter)
	for _, value := range fields {
		parts := strings.Split(value, "=")
		if len(parts) > 1 {
			key := strings.ToLower(parts[0])
			value := parts[1]
			dict[key] = strings.Trim(value, " ")
		}
	}

	g.packageName = dict["pkgname"]
	if g.packageName == "" {
		g.packageName = defaultPackageName
	}

	if dict["types"] != "" {
		g.types = strings.Split(dict["types"], ",")
	}
}

// fillVersion saves protoc version number as string
func (g *Generator) fillVersion() {
	if v := g.Plugin.Request.GetCompilerVersion(); v != nil {
		g.protocVersion = fmt.Sprintf("v%v.%v.%v", v.GetMajor(), v.GetMinor(), v.GetPatch())
	}
}

// Generate is an entry point for the generator.
// main() invokes Generate after initializing the plugin and expects it to
// return a bunch of generated files and errors.
//
// The errors will be rendered to STDERR, and should be runtime / execution errors.
// Proto-specific and generation errors should be reported in g.Plugin.Error,
// that way protoc will intercept them. main() handles grabbing those errors from
// g.PLugin.
//
// Generate() builds the map of all of the available messages first, and then
// generates schemas for required messages.
func (g *Generator) Generate() ([]*protogen.GeneratedFile, error) {
	var files []*protogen.GeneratedFile

	for _, file := range g.Plugin.Files {
		if !file.Generate {
			continue
		}

		// Output strings builder
		result := builder{}

		// Generates individual file content
		fileGenerator{
			generator: g,
			file:      file,
			builder:   &result,
		}.generate()

		// generate the tfschema source based on the proto
		// generated, err := g.generateFile(file)
		// if err != nil {
		// 	// main wraps the call to Generate and puts emits errors to the protoc binary.
		// 	// if there's a runtime error not related to the generator itself, we log it to stdout.
		// 	// but if there'a generator error, we just return it.
		// 	return files, trace.Wrap(err)
		// }

		// generate the file name and write the file
		filename := file.GeneratedFilenamePrefix + ".tfschema.go"
		out := g.Plugin.NewGeneratedFile(filename, ".")

		// _, err = out.Write([]byte(generated))
		_, err := out.Write(result.Bytes())
		if err != nil {
			log.Errorf("Couldn't write %s, error: %v", filename, err)
			return files, trace.Wrap(err)
		}

		files = append(files, out)
		log.Infof("Will emit %s", filename)
	}

	return files, nil
}

// generateFile generates go code implementing Terraform schema for
// resources provided in a single proto file
func (g *Generator) generateFile(file *protogen.File) (string, error) {

	// // start an output string by generating the package name.
	// generated := g.generatePackageName(file)

	// // filter out a list of the messages to generate.
	// // This might include all messages in the proto file if types option is omitted,
	// // Or a list of the selected types.
	// messages := g.filterRequiredMessages(file.Desc.Messages())

	// // Append the generated message schema for each of the messages
	// for _, message := range messages {
	// 	generated += g.generateMessageSchema(message)
	// }

	// return generated, nil
	return "", nil
}

// TODO: Generator functions are using a bunch of `fmt.Sprintf` and that doesn't look nice.
// Refactor this into templates?

// generateInlineMessageSchema renders a schema for a complex type field inline in a message schema
func (g *Generator) generateInlineMessageSchema(message protoreflect.MessageDescriptor) string {

	fields := []string{}
	for _, field := range g.fieldsToList(message.Fields()) {
		fields = append(fields, g.generateFieldSchema(field))
	}

	return fmt.Sprintf(`&schema.Resource{
		// nested type name: %s
		Schema: map[string]*schema.Schema{
			%s
		},
	}`, message.FullName(), strings.Join(fields, ""))
}

// generateMessage generates a terraform schema go code for one MessageDescriptor (type)
// and all of it's internals and nested types.
func (g *Generator) generateMessageSchema(message protoreflect.MessageDescriptor) string {

	// generate each field schema,
	// then join them together and
	// wrap them with the schema function
	fields := []string{}
	for _, field := range g.fieldsToList(message.Fields()) {
		fields = append(fields, g.generateFieldSchema(field))
	}

	return fmt.Sprintf(`
	// proto type fullname: %s
	func schema%s() map[string]*schema.Schema {
		return map[string]*schema.Schema{
			%s
		}
	}
	`, message.FullName(), message.Name(), strings.Join(fields, ""))
}

// generateFieldSchema generates a schema for a single field in a proto message
// TODO: refactor this to use switch, and maybe extract into several functions.
func (g *Generator) generateFieldSchema(field protoreflect.FieldDescriptor) string {

	if field.IsMap() {
		// map with string keys and primitive type values
		if field.MapValue().Kind() != protoreflect.MessageKind {
			return fmt.Sprintf(`"%s": {
				Type: schema.TypeMap,
				%s,
				Elem: &schema.Schema{
					Type: schema.%s,
				},
			},
			`, field.Name(), requiredOrOptional(field), kindToTerraform(field.MapValue().Kind()))
		}

		// map with string keys and complex type values
		if field.MapValue().Kind() != protoreflect.MessageKind {
			return fmt.Sprintf(`"%s": {
				Type: schema.TypeMap,
				%s,
				Elem: %s
			},
			`, field.Name(), requiredOrOptional(field), g.generateInlineMessageSchema(field.MapValue().Message()))
		}
	}

	if field.IsList() {
		// List of complex type objects
		if field.Kind() == protoreflect.MessageKind {
			return fmt.Sprintf(`"%s": {
				Type: schema.TypeList,
				%s,
				Elem: %s,
			},
			`, field.Name(), requiredOrOptional(field), g.generateInlineMessageSchema(field.Message()))
		}

		// List of primitive type objects
		return fmt.Sprintf(`"%s": {
			Type: schema.TypeList,
			%s,
			Elem: &schema.Schema{Type: schema.%s},
		},
		`, field.Name(), requiredOrOptional(field), kindToTerraform(field.Kind()))
	}

	if field.Kind() == protoreflect.MessageKind {

		// check if the field has one of a known mappable types
		// and generate an inline schema for types we don't recognize
		switch field.Message().FullName() {
		case "google.protobuf.Timestamp":
			return fmt.Sprintf(`"%s": {
				Type: schema.TypeString,
				%s,
				ValidateFunc: validation.ValidateRFC3339TimeString,
			},
			`, field.Name(), requiredOrOptional(field))
		default:

			// fields of custom time are represented as lists in TF schema
			return fmt.Sprintf(`"%s": {
				Type: schema.TypeList,
				%s,
				MaxItems:1,
				Elem: %s,
				},
			`, field.Name(), requiredOrOptional(field), g.generateInlineMessageSchema(field.Message()))
		}

	}

	return fmt.Sprintf(`"%s": {
		Type: schema.%s,
		%s,
	},
	`, field.Name(), kindToTerraform(field.Kind()), requiredOrOptional(field))

}

// fieldsToList flattens a protoreflect.FieldDescriptors
// into a []protoreflect.FieldDescriptor array
func (g *Generator) fieldsToList(fields protoreflect.FieldDescriptors) []protoreflect.FieldDescriptor {
	list := []protoreflect.FieldDescriptor{}
	for i := 0; i < fields.Len(); i++ {
		list = append(list, fields.Get(i))
	}
	return list
}

// helpers
// TODO: extract helpers into a separate file
//

// kindToTerraform returns terraform schema equivalent
// of a proto field Kind
func kindToTerraform(kind protoreflect.Kind) string {
	if tfType, ok := kind2tf[kind]; ok {
		return tfType
	}

	log.Warnf("Unspecified field type: %s", kind)
	return fmt.Sprintf("%v", kind)
}

func requiredOrOptional(field protoreflect.FieldDescriptor) string {
	if field.Cardinality() == protoreflect.Required {
		return "Required: true"
	}
	return "Optional: true"
}
