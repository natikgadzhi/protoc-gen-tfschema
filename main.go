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

// protoc invokes main and provides cli params.
// main creates a new instance of the Plugin,
// runs the generator, and sends the result back
// to stdout for protoc.
func main() {
	log.Info("Reading CodeGeneratorRequest from stdin.")
	request, err := readRequest()
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Initializing protoc-gen-schema plugin.")
	opts := protogen.Options{}
	plugin, err := opts.New(request)
	if err != nil {
		log.Fatal(err)
	}
	generator := gen.NewGenerator(plugin)

	log.Info("Generating schema files.")
	generated, err := generator.Generate()
	if err != nil {
		log.Errorf("Error generating schemas: %v", err)
		generator.Plugin.Error(err)
	}

	log.Infof("Done, generated %d files", len(generated))

	emitResponse(generator.Plugin.Response())
}

// readRequest reads the protoc-gen Code Generator request
// with a list of proto files to work through, and
// unmarshalls it from proto format.
func readRequest() (*pluginpb.CodeGeneratorRequest, error) {
	var request pluginpb.CodeGeneratorRequest
	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return nil, trace.Wrap(err)
	}

	proto.Unmarshal(input, &request)
	return &request, nil
}

// emitResponse marshalls the response to a protobuf message
// and sends it back to protoc-gen via stdout
func emitResponse(resp *pluginpb.CodeGeneratorResponse) {
	buf, err := proto.Marshal(resp)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := os.Stdout.Write(buf); err != nil {
		log.Fatal(err)
	}
}
