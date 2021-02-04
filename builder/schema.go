package builder

import (
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Schema mock Terraform Schema with metadata
type Schema struct {
	Name         string
	FullName     string
	Required     bool
	Type         string
	ValidateFunc string
	MaxItems     int
	Elem         interface{}
}

type schemaBuilder struct {
	field       protoreflect.FieldDescriptor
	schema      *Schema
	resourceMap ResourceMap
}

// BuildSchemaFromField builds resource from protoreflect message
func BuildSchemaFromField(field *protoreflect.FieldDescriptor, resourceMap ResourceMap) *Schema {
	schema := &Schema{}

	b := schemaBuilder{field: *field, schema: schema, resourceMap: resourceMap}

	b.setName()
	b.setFullName()
	b.setRequired()
	b.setTypeAndElem()

	return schema
}

func (b *schemaBuilder) setName() {
	b.schema.Name = string(b.field.Name())
}

func (b *schemaBuilder) setFullName() {
	b.schema.FullName = string(b.field.FullName())
}

func (b *schemaBuilder) setRequired() {
	b.schema.Required = b.field.Cardinality() == protoreflect.Required
}

// Returns true if current field contains nested resource (and not the list of nested resources)
func (b *schemaBuilder) isNestedResource() bool {
	return b.field.Kind() == protoreflect.MessageKind && !b.field.IsList()
}

// Sets type and nested elem type
func (b *schemaBuilder) setTypeAndElem() {
	kind := b.field.Kind()

	// NOTE: How to treat MapKey()?
	if b.field.IsMap() {
		b.setMap()
	} else if b.field.IsList() {
		b.schema.Type = "schema.TypeList"
		b.setElem(b.field.Kind(), b.field.Message())
	} else if b.isNestedResource() {
		b.setNestedResourceTypeAndElem()
	} else {
		b.schema.Type = b.getTypeFromKind(kind)
	}
}

// Converts map of resources (are not supported by TerraForm) to list of resources with mock name field
func (b *schemaBuilder) setMap() {
	value := b.field.MapValue()

	// Map with embedded ordinal type is OK
	if value.Kind() != protoreflect.MessageKind {
		b.schema.Type = "schema.TypeMap"
		b.setElem(value.Kind(), value.Message())
		return
	}

	// Embedded resource has a single field: we could omit it
	resource := b.buildNestedResource(value.Message())
	b.addMockNameField(resource)
	b.schema.Type = "schema.TypeList"
	b.schema.Elem = resource.Name
}

// Adds mock name field to resource
func (b *schemaBuilder) addMockNameField(resource *Resource) {
	if resource.Schema["name"] != nil {
		logrus.Fatalf("Failed to convert map of resources to list of resources with name: name field exists")
	}

	resource.Schema["name"] = &Schema{
		Name:     "name",
		FullName: "mocks.Name",
		Type:     "schema.TypeString",
		Required: true,
	}
}

// Branch for nested resources
func (b *schemaBuilder) setNestedResourceTypeAndElem() {
	if b.isTimeStamp() {
		b.schema.Type = "schema.TypeString"
		b.schema.ValidateFunc = "validation.IsRFC3339Time"
	} else {
		// If the nested resource is another structure, than we should produce a list with the single item
		// That's the weirdo way Terraform handles such case
		b.schema.Type = "schema.TypeList"
		b.schema.MaxItems = 1
		b.setElem(b.field.Kind(), b.field.Message())
	}
}

// Sets nested element type and builds schema for it if required
func (b *schemaBuilder) setElem(kind protoreflect.Kind, message protoreflect.MessageDescriptor) {
	var elem interface{}

	if kind == protoreflect.MessageKind {
		r := b.buildNestedResource(message)
		elem = r.Name
	} else {
		s := Schema{}
		s.Type = b.getTypeFromKind(kind)
		elem = &s
	}

	b.schema.Elem = elem
}

func (b *schemaBuilder) buildNestedResource(message protoreflect.MessageDescriptor) *Resource {
	resource := BuildResourceFromMessage(&message, b.resourceMap)
	b.resourceMap[resource.Name] = resource
	return resource
}

// Converts protoc kind to Terraform type
func (b *schemaBuilder) getTypeFromKind(kind protoreflect.Kind) string {
	switch kind {
	case protoreflect.BoolKind:
		return "schema.TypeBool"
	case protoreflect.StringKind, protoreflect.BytesKind, protoreflect.EnumKind:
		return "schema.TypeString"
	case protoreflect.Int32Kind, protoreflect.Int64Kind,
		protoreflect.Fixed32Kind, protoreflect.Fixed64Kind,
		protoreflect.Uint32Kind, protoreflect.Uint64Kind:
		return "schema.TypeInt"
	case protoreflect.DoubleKind:
		return "schema.TypeFloat"
	}

	// Should not have happened
	logrus.Fatalf("Unknown protoreflect kind %#v", kind)

	return "schema.Unknown"
}

// Returns true if a field represents TimeStamp
func (b *schemaBuilder) isTimeStamp() bool {
	return string(b.field.Message().FullName()) == "google.protobuf.Timestamp"
}
