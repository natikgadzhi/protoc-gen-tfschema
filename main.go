package main

import (
	"io/ioutil"
	"os"

	"github.com/gravitational/trace"
	log "github.com/sirupsen/logrus"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/nategadzhi/protoc-gen-tfschema/gen"
)

var request = pluginpb.CodeGeneratorRequest{}

func main() {
	// 1. Read request from stdin
	log.Info("Reading CodeGeneratorRequest from stdin.")
	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fatal(err)
	}
	proto.Unmarshal(input, &request)

	// 2. Initialize plugin
	log.Info("Initializing protoc-gen-tfschema plugin.")
	opts := protogen.Options{}
	plugin, err := opts.New(&request)
	if err != nil {
		fatal(err)
	}

	// 3. Run plugin, generate schema files
	log.Info("Generating schema files...")
	generator := gen.NewGenerator(plugin)
	generated, err := generator.Generate()
	if err != nil {
		log.Errorf("Error generating schemas: %v", err)
		generator.Plugin.Error(err)
	}
	log.Infof("Done, generated %d files", len(generated))

	// 4. Put response back to stdout
	buf, err := proto.Marshal(generator.Plugin.Response())
	if err != nil {
		fatal(err)
	}
	if _, err := os.Stdout.Write(buf); err != nil {
		fatal(err)
	}
}

func fatal(err error) {
	log.Fatal(trace.Wrap(err).Error())
}
