package builder

import (
	"github.com/nategadzhi/protoc-gen-tfschema/config"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// ResourceMap string -> Resource map
type ResourceMap map[string]*Resource

// BuildResourceMapFromFile returns SchemaMap parsed from the file
func BuildResourceMapFromFile(file *protogen.File) *ResourceMap {
	m := make(ResourceMap)
	messages := file.Desc.Messages()

	for i := 0; i < messages.Len(); i++ {
		message := messages.Get(i)

		if isMessageRequired(message) {
			resource := BuildResourceFromMessage(&message)
			m[resource.Name] = resource
		}
	}

	return &m
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
