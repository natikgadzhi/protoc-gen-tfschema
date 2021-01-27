package reader

import (
	"github.com/nategadzhi/protoc-gen-tfschema/builder"
	"github.com/nategadzhi/protoc-gen-tfschema/config"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// ReadFile returns SchemaMap parsed from the file
func ReadFile(file *protogen.File) *builder.SchemaMap {
	messages := file.Desc.Messages()

	for i := 0; i < messages.Len(); i++ {
		message := messages.Get(i)

		if isMessageRequired(message) {

		}
	}

	return &builder.SchemaMap{}
}

// Checks if current message needs to be parsed & exported
func isMessageRequired(message protoreflect.MessageDescriptor) bool {
	if len(config.Types) == 0 {
		return true
	}

	for _, v := range config.Types {
		if v == string(message.Name()) {
			return true
		}
	}

	return false
}
