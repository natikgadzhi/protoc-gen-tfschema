package builder

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Resource extends Terraform resource with metadata required for file generation
type Resource struct {
	Name     string
	FullName string
	Schema   SchemaMap // Overrrides schema.Resource.Schema

	schema.Resource
}

// NewResource builds and initializes an empty Resource
func NewResource() *Resource {
	return &Resource{
		Schema: make(SchemaMap),
	}
}

type resourceBuilder struct {
	message  protoreflect.MessageDescriptor
	resource *Resource
}

func (b *resourceBuilder) setName() {
	b.resource.Name = string(b.message.Name())
}

func (b *resourceBuilder) setFullName() {
	b.resource.FullName = string(b.message.FullName())
}

// BuildResourceFromMessage creates new resource from message
func BuildResourceFromMessage(message *protoreflect.MessageDescriptor) *Resource {
	resource := NewResource()

	builder := resourceBuilder{message: *message, resource: resource}

	builder.setName()
	builder.setFullName()

	return resource
}
