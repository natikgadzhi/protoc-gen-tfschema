package reader

import (
	"github.com/nategadzhi/protoc-gen-tfschema/builder"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// ReadMessage reads and parses resource into schema map
func ReadMessage(message *protoreflect.MessageDescriptor, resource *builder.Resource) {
	fields := (*message).Fields()

	for i := 0; i < fields.Len(); i++ {
		field := fields.Get(i)
		name := string(field.Name())
		fullName := string(field.FullName())

		schema := resource.AddSchema(name, fullName)
		schema.SetType(getType(message))
	}
}
