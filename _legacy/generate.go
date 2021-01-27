// kind2tf is a mapping from protobuf primitive types to terraform types
var kind2tf = map[protoreflect.Kind]string{
	protoreflect.BoolKind: "TypeBool",

	protoreflect.StringKind: "TypeString",
	protoreflect.BytesKind:  "TypeString",
	protoreflect.EnumKind:   "TypeString",

	protoreflect.Int32Kind:   "TypeInt",
	protoreflect.Int64Kind:   "TypeInt",
	protoreflect.Fixed32Kind: "TypeInt",
	protoreflect.Fixed64Kind: "TypeInt",
	protoreflect.Uint32Kind:  "TypeInt",
	protoreflect.Uint64Kind:  "TypeInt",
	protoreflect.DoubleKind:  "TypeFloat",
}

// generateFile generates go code implementing Terraform schema for
// resources provided in a single proto file
func (g *Generator) generateFile(file *protogen.File) (string, error) {

	// // start an output string by generating the package name.
	// generated := g.generatePackageName(file)

	// // filter out a list of the messages to generate.
	// // This might include all messages in the proto file if types option is omitted,
	// // Or a list of the selected types.
	// messages := g.filterRequiredMessages(file.Desc.Messages())

	// // Append the generated message schema for each of the messages
	// for _, message := range messages {
	// 	generated += g.generateMessageSchema(message)
	// }

	// return generated, nil
	return "", nil
}

// TODO: Generator functions are using a bunch of `fmt.Sprintf` and that doesn't look nice.
// Refactor this into templates?

// generateInlineMessageSchema renders a schema for a complex type field inline in a message schema
func (g *Generator) generateInlineMessageSchema(message protoreflect.MessageDescriptor) string {

	fields := []string{}
	for _, field := range g.fieldsToList(message.Fields()) {
		fields = append(fields, g.generateFieldSchema(field))
	}

	return fmt.Sprintf(`&schema.Resource{
		// nested type name: %s
		Schema: map[string]*schema.Schema{
			%s
		},
	}`, message.FullName(), strings.Join(fields, ""))
}

// generateMessage generates a terraform schema go code for one MessageDescriptor (type)
// and all of it's internals and nested types.
func (g *Generator) generateMessageSchema(message protoreflect.MessageDescriptor) string {

	// generate each field schema,
	// then join them together and
	// wrap them with the schema function
	fields := []string{}
	for _, field := range g.fieldsToList(message.Fields()) {
		fields = append(fields, g.generateFieldSchema(field))
	}

	return fmt.Sprintf(`
	// proto type fullname: %s
	func schema%s() map[string]*schema.Schema {
		return map[string]*schema.Schema{
			%s
		}
	}
	`, message.FullName(), message.Name(), strings.Join(fields, ""))
}

// generateFieldSchema generates a schema for a single field in a proto message
// TODO: refactor this to use switch, and maybe extract into several functions.
func (g *Generator) generateFieldSchema(field protoreflect.FieldDescriptor) string {

	if field.IsMap() {
		// map with string keys and primitive type values
		if field.MapValue().Kind() != protoreflect.MessageKind {
			return fmt.Sprintf(`"%s": {
				Type: schema.TypeMap,
				%s,
				Elem: &schema.Schema{
					Type: schema.%s,
				},
			},
			`, field.Name(), requiredOrOptional(field), kindToTerraform(field.MapValue().Kind()))
		}

		// map with string keys and complex type values
		if field.MapValue().Kind() != protoreflect.MessageKind {
			return fmt.Sprintf(`"%s": {
				Type: schema.TypeMap,
				%s,
				Elem: %s
			},
			`, field.Name(), requiredOrOptional(field), g.generateInlineMessageSchema(field.MapValue().Message()))
		}
	}

	if field.IsList() {
		// List of complex type objects
		if field.Kind() == protoreflect.MessageKind {
			return fmt.Sprintf(`"%s": {
				Type: schema.TypeList,
				%s,
				Elem: %s,
			},
			`, field.Name(), requiredOrOptional(field), g.generateInlineMessageSchema(field.Message()))
		}

		// List of primitive type objects
		return fmt.Sprintf(`"%s": {
			Type: schema.TypeList,
			%s,
			Elem: &schema.Schema{Type: schema.%s},
		},
		`, field.Name(), requiredOrOptional(field), kindToTerraform(field.Kind()))
	}

	if field.Kind() == protoreflect.MessageKind {

		// check if the field has one of a known mappable types
		// and generate an inline schema for types we don't recognize
		switch field.Message().FullName() {
		case "google.protobuf.Timestamp":
			return fmt.Sprintf(`"%s": {
				Type: schema.TypeString,
				%s,
				ValidateFunc: validation.ValidateRFC3339TimeString,
			},
			`, field.Name(), requiredOrOptional(field))
		default:

			// fields of custom time are represented as lists in TF schema
			return fmt.Sprintf(`"%s": {
				Type: schema.TypeList,
				%s,
				MaxItems:1,
				Elem: %s,
				},
			`, field.Name(), requiredOrOptional(field), g.generateInlineMessageSchema(field.Message()))
		}

	}

	return fmt.Sprintf(`"%s": {
		Type: schema.%s,
		%s,
	},
	`, field.Name(), kindToTerraform(field.Kind()), requiredOrOptional(field))

}

// fieldsToList flattens a protoreflect.FieldDescriptors
// into a []protoreflect.FieldDescriptor array
func (g *Generator) fieldsToList(fields protoreflect.FieldDescriptors) []protoreflect.FieldDescriptor {
	list := []protoreflect.FieldDescriptor{}
	for i := 0; i < fields.Len(); i++ {
		list = append(list, fields.Get(i))
	}
	return list
}

// helpers
// TODO: extract helpers into a separate file
//

// kindToTerraform returns terraform schema equivalent
// of a proto field Kind
func kindToTerraform(kind protoreflect.Kind) string {
	if tfType, ok := kind2tf[kind]; ok {
		return tfType
	}

	log.Warnf("Unspecified field type: %s", kind)
	return fmt.Sprintf("%v", kind)
}

func requiredOrOptional(field protoreflect.FieldDescriptor) string {
	if field.Cardinality() == protoreflect.Required {
		return "Required: true"
	}
	return "Optional: true"
}
