package builder

import (
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Resource mock Terraform resource with metadata required for file generation
type Resource struct {
	Name     string
	FullName string
	Schema   SchemaMap // Overrrides schema.Resource.Schema
}

type resourceBuilder struct {
	message  protoreflect.MessageDescriptor
	resource *Resource
}

// NewResource builds and initializes an empty Resource
func NewResource() *Resource {
	return &Resource{
		Schema: make(SchemaMap),
	}
}

// BuildResourceFromMessage creates new resource from message
func BuildResourceFromMessage(message *protoreflect.MessageDescriptor, resourceMap ResourceMap) *Resource {
	resource := NewResource()

	builder := resourceBuilder{message: *message, resource: resource}

	builder.setName()
	builder.setFullName()

	BuildSchemaMapFromMessage(message, resource, resourceMap)

	return resource
}

func (b *resourceBuilder) setName() {
	b.resource.Name = string(b.message.Name())
}

func (b *resourceBuilder) setFullName() {
	b.resource.FullName = string(b.message.FullName())
}
