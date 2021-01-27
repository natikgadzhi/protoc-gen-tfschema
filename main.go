package main

import (
	"io/ioutil"
	"os"

	"github.com/gravitational/trace"
	"github.com/nategadzhi/protoc-gen-tfschema/config"
	"github.com/nategadzhi/protoc-gen-tfschema/reader"
	log "github.com/sirupsen/logrus"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

var (
	request = pluginpb.CodeGeneratorRequest{}
	plugin  *protogen.Plugin
)

const (
	generatedFileSuffix = ".tfschema.go"
)

func init() {
	initRequest()
	initPlugin()
}

func main() {
	log.Info("Generating schema files...")

	generate()
	emitResponse()
}

// Parses and initializes CodeGeneratorRequest
func initRequest() {
	log.Info("Reading CodeGeneratorRequest from stdin.")

	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fatal(err)
	}
	proto.Unmarshal(input, &request)

	log.Info("Command-line arguments: %s", request.GetParameter())
}

// Parses command line options and initializes Plugin instance
func initPlugin() {
	var err error

	opts := &protogen.Options{
		ParamFunc: config.Set, // Built-in way of using go flag package to parse CLI args
	}
	plugin, err = opts.New(&request)
	if err != nil {
		fatal(err)
	}

	config.Finalize()
}

// Parses input files and writes parsed files back to plugin
func generate() {
	var numFilesWritten = 0

	for _, file := range plugin.Files {
		if !file.Generate {
			continue
		}

		// Target file name
		filename := file.GeneratedFilenamePrefix + generatedFileSuffix
		out := plugin.NewGeneratedFile(filename, ".")

		_ = reader.ReadFile(file)
		// parsed = Parser.file
		// result = Renderer.file

		//_, err := out.Write(result.Bytes())
		_, err := out.Write([]byte("ok"))
		if err != nil {
			log.Errorf("Error generating schemas: %v", err)
			plugin.Error(err)
		}

		log.Infof("%s prepared", filename)
		numFilesWritten++
	}

	log.Infof("%d files generated", numFilesWritten)
}

func emitResponse() {
	buf, err := proto.Marshal(plugin.Response())
	if err != nil {
		fatal(err)
	}
	if _, err := os.Stdout.Write(buf); err != nil {
		fatal(err)
	}
}

// Sugar
func fatal(err error) {
	log.Fatal(trace.Wrap(err).Error())
}
