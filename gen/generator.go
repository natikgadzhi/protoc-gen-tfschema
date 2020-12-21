package gen

import (
	"fmt"
	"strings"

	"github.com/gravitational/trace"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/compiler/protogen"
)

type schemaGenerator struct {
	Plugin *protogen.Plugin
}

// NewGenerator returns new generator instance.
func NewGenerator(plugin *protogen.Plugin) *schemaGenerator {
	return &schemaGenerator{
		Plugin: plugin,
	}
}

// Generate takes a list of protogen targets (proto files), and
// generates files with Go code implementing terraform schema for
// resources provided in the proto files.
func (g *schemaGenerator) Generate() ([]*protogen.GeneratedFile, error) {
	var files []*protogen.GeneratedFile

	for _, file := range g.Plugin.Files {

		generated, err := g.generateFile(file)
		if err != nil {
			return files, trace.Wrap(err)
		}

		// TODO move filename suffix to a const
		filename := file.GeneratedFilenamePrefix + ".schema.go"
		file := g.Plugin.NewGeneratedFile(filename, ".")
		_, err = file.Write([]byte(generated))
		if err != nil {
			return files, trace.Wrap(err)
		}

		files = append(files, file)
		log.Infof("Will emit %s", filename)
	}

	return files, nil
}

// generateFile generates go code implementing Terraform schema for
// resources provided in a single proto file
func (g *schemaGenerator) generateFile(file *protogen.File) (string, error) {
	generated := []string{}

	generated = append(generated, fmt.Sprintf(`package %s`, file.GoPackageName))

	for _, msg := range file.Proto.MessageType {
		generated = append(generated, fmt.Sprintf(`
		func (x %s) Foo() string {
			return "bar"
		}`, *msg.Name))
	}

	return strings.Join(generated, "\n"), nil
}
