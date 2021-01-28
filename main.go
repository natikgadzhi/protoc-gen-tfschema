package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gravitational/trace"
	"github.com/nategadzhi/protoc-gen-tfschema/builder"
	"github.com/nategadzhi/protoc-gen-tfschema/config"
	"github.com/nategadzhi/protoc-gen-tfschema/renderer"
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

	log.Infof("Command-line arguments: %s", request.GetParameter())
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

	// Save ProtocVersion
	if v := plugin.Request.GetCompilerVersion(); v != nil {
		config.ProtocVersion = fmt.Sprintf("v%v.%v.%v", v.GetMajor(), v.GetMinor(), v.GetPatch())
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

		// Build resource tree
		resources := builder.BuildResourceMapFromFile(file)

		// Render final template
		result, err := renderer.Render(resources, Version)
		if err != nil {
			log.Errorf("Error rendering template: %v", err)
			plugin.Error(err)
		}

		_, err = out.Write(result.Bytes())
		if err != nil {
			log.Errorf("Error generating schemas: %v", err)
			plugin.Error(err)
		}

		log.Infof("%s prepared", filename)
		numFilesWritten++
	}

	log.Infof("%v files generated", numFilesWritten)
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
