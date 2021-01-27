package reader

import (
	"github.com/nategadzhi/protoc-gen-tfschema/builder"
	"github.com/nategadzhi/protoc-gen-tfschema/util"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Parse returns SchemaMap parsed from the file
func ReadFile(file *protogen.File) *builder.SchemaMap {
	messages := file.Desc.Messages()

	for i := 0; i < messages.Len(); i++ {
		message := messages.Get(i)

		if isMessageParseable(message) {
		}
	}

	return &builder.SchemaMap{}
}

// Checks if current message needs to be exported (name was passed in command line argument or argument is empty)
func isMessageParseable(message protoreflect.MessageDescriptor) bool {
	if len(main.Types) == 0 {
		return true
	}

	return util.StrSliceContains(Types, string(message.Name()))
}
