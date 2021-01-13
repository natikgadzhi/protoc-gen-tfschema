package gen

import (
	"fmt"
	"strings"

	"github.com/gravitational/trace"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	// defaultPackageName is the default package name to generate
	// if `package` option is not provided in the CLI options.
	defaultPackageName = "tfschema"
)

var kind2tf = map[protoreflect.Kind]string{
	protoreflect.BoolKind:   "TypeBool",
	protoreflect.StringKind: "TypeString",
	protoreflect.BytesKind:  "TypeString",
	protoreflect.EnumKind:   "TypeString",

	protoreflect.Int32Kind:   "TypeInt",
	protoreflect.Int64Kind:   "TypeInt",
	protoreflect.DoubleKind:  "TypeFloat",
	protoreflect.Fixed32Kind: "TypeInt",
	protoreflect.Fixed64Kind: "TypeInt",
	protoreflect.Uint32Kind:  "TypeInt",
	protoreflect.Uint64Kind:  "TypeInt",
}

// Generator is the Terraform Schema generator
type Generator struct {
	Plugin  *protogen.Plugin
	options *GeneratorOptions
}

// GeneratorOptions holds all of the options for the terraform schema generator that are parsed from the command line flags
type GeneratorOptions struct {
	// Types is a list of type names in the proto to generate terraform schemas for.
	types       []string
	packageName string
}

// NewGenerator returns new generator instance.
func NewGenerator(plugin *protogen.Plugin) *Generator {
	g := &Generator{Plugin: plugin}
	g.options = parseCommandLineOptions(g.Plugin.Request.GetParameter())
	return g
}

// parseCommandLineOptions parses params passed from protoc and returns a GeneratorOptions struct
func parseCommandLineOptions(params string) *GeneratorOptions {
	// make a map of options and flags
	opts := map[string]string{}
	for _, p := range strings.Split(params, ",") {
		if i := strings.Index(p, "="); i < 0 {
			opts[p] = ""
		} else {
			opts[p[0:i]] = p[i+1:]
		}
	}

	var generateMessages []string
	if types, ok := opts["types"]; ok {
		generateMessages = strings.Split(types, " ")
	}

	packageName := opts["package"]
	if packageName == "" {
		packageName = defaultPackageName
	}

	return &GeneratorOptions{
		types:       generateMessages,
		packageName: packageName,
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

		// generate the tfschema source based on the proto
		generated, err := g.generateFile(file)
		if err != nil {
			// main wraps the call to Generate and puts emits errors to the protoc binary.
			// if there's a runtime error not related to the generator itself, we log it to stdout.
			// but if there'a generator error, we just return it.
			return files, trace.Wrap(err)
		}

		// generate the file name and write the file
		filename := file.GeneratedFilenamePrefix + ".tfschema.go"
		out := g.Plugin.NewGeneratedFile(filename, ".")
		_, err = out.Write([]byte(generated))
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

	// start an output string by generating the package name.
	generated := g.generatePackageName()

	// filter out a list of the messages to generate.
	// This might include all messages in the proto file if types option is omitted,
	// Or a list of the selected types.
	messages := g.filterRequiredMessages(file.Desc.Messages())

	// Append the generated message schema for each of the messages
	for _, message := range messages {
		generated += g.generateMessageSchema(message)
	}

	return generated, nil
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
	func schema%s() map[string]*schema.Schema {
		return map[string]*schema.Schema{
			%s
		}
	}
	`, message.Name(), strings.Join(fields, ""))
}

// generateFieldSchema generates a schema for a single field in a proto message
func (g *Generator) generateFieldSchema(field protoreflect.FieldDescriptor) string {

	if field.IsMap() {

	}

	if field.IsList() {

	}

	if field.Kind() == protoreflect.MessageKind {

	}

	return fmt.Sprintf(`"%s": {
		Type: schema.%s,
		%s,
	},
	`, field.Name(), kindToTerraform(field.Kind()), requiredOrOptional(field.HasOptionalKeyword()))
}

// generatePackageName generates package name string
// and imports string for the go source cod
func (g *Generator) generatePackageName() string {
	return fmt.Sprintf(`
	package %s

	import (
		"github.com/hashicorp/terraform/helper/schema"
	)

	`, g.options.packageName)
}

// filterRequiredMessages takes a MessageDescriptors, flattens it to a list and filters only the required messages.
// TODO: add a validation to check that all required messages are found
// in the supplied proto files, error out otherwise.
func (g *Generator) filterRequiredMessages(messages protoreflect.MessageDescriptors) []protoreflect.MessageDescriptor {

	messagesList := g.messagesToList(messages)
	if len(g.options.types) == 0 {
		return messagesList
	}

	filtered := []protoreflect.MessageDescriptor{}
	for _, m := range messagesList {
		if g.isMessageRequired(m) {
			filtered = append(filtered, m)
		}
	}
	return filtered
}

// messagesToList flattens a protoreflect.MessageDescriptors
// into a []protoreflect.MessageDescriptor
func (g *Generator) messagesToList(messages protoreflect.MessageDescriptors) []protoreflect.MessageDescriptor {
	list := []protoreflect.MessageDescriptor{}
	for i := 0; i < messages.Len(); i++ {
		list = append(list, messages.Get(i))
	}
	return list
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

// isMessageRequired returns true for Proto Type names if they're required
// to be generated by the protoc command line options.
//
// `tfschema_opt=types=Metadata` will set Metadata message to be required.
// If no option was provided, all messages are required.
func (g *Generator) isMessageRequired(message protoreflect.MessageDescriptor) bool {

	// if types are not provided in options, we should generate all types in the proto.
	if len(g.options.types) == 0 {
		return true
	}

	// if types option is provided, check if the message type is one of the provided types.
	for _, t := range g.options.types {
		if t == string(message.Name()) {
			return true
		}
	}
	return false
}

// helpers
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

func requiredOrOptional(required bool) string {
	if required {
		return "Required: true"
	}
	return "Optional: true"
}
