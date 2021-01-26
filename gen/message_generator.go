package gen

import (
	"google.golang.org/protobuf/reflect/protoreflect"
)

type messageGenerator struct {
	message protoreflect.MessageDescriptor
	*builder
}

func (m messageGenerator) generate() {
	m.Printf("// proto type full name: %s", m.message.FullName())
	m.Printf("func schema%s() map[string]*schema.Schema {", m.message.Name())
	m.Printf("    return map[string]*schema.Schema{")

	fields := m.message.Fields()
	for i := 0; i < fields.Len(); i++ {
		m.genField(fields.Get(i))
	}

	m.Printf("    }")
	m.Printf("}")
	m.CrLf()
}

func (m *messageGenerator) genField(field protoreflect.FieldDescriptor) {
	m.Printf(`"%s": {`, field.Name()) // Field name

	m.printfFieldCardinality(field) // Required or not
	m.printfFieldType(field)        // Type with some additional traits

	m.Printf("},")
}

func (m *messageGenerator) printfFieldCardinality(field protoreflect.FieldDescriptor) {
	if field.Cardinality() == protoreflect.Required {
		m.Printf("    Required: true,")
	} else {
		m.Printf("    Optional: true,")
	}
}

func (m *messageGenerator) printfFieldTypeDecl(t string) {
	m.Printf("    Type: %s,", t)
}

func (m *messageGenerator) printfFieldType(field protoreflect.FieldDescriptor) {
	if field.IsMap() {
		m.printfFieldTypeDecl("schema.TypeMap")
	} else if field.IsList() {
		m.printfFieldTypeDecl("schema.TypeList")
	} else {
		switch field.Kind() {
		case protoreflect.MessageKind:
			m.printfFieldTypeDecl("schema.TypeList")
		case protoreflect.BoolKind:
			m.printfFieldTypeDecl("schema.TypeBool")
		case protoreflect.StringKind, protoreflect.BytesKind, protoreflect.EnumKind:
			m.printfFieldTypeDecl("schema.TypeString")
		case protoreflect.Int32Kind, protoreflect.Int64Kind,
			protoreflect.Fixed32Kind, protoreflect.Fixed64Kind,
			protoreflect.Uint32Kind, protoreflect.Uint64Kind:
			m.printfFieldTypeDecl("schema.TypeInt")
		case protoreflect.DoubleKind:
			m.printfFieldTypeDecl("schema.TypeFloat")
		default:
			// TODO
		}
	}
}

// // generateFieldSchema generates a schema for a single field in a proto message
// // TODO: refactor this to use switch, and maybe extract into several functions.
// func (g *Generator) generateFieldSchema(field protoreflect.FieldDescriptor) string {
// 	if field.IsMap() {
// 		// map with string keys and primitive type values
// 		if field.MapValue().Kind() == protoreflect.MessageKind {
// 			return fmt.Sprintf(`"%s": {
// 				Type: schema.TypeMap,
// 				%s,
// 				Elem: &schema.Schema{
// 					Type: schema.%s,
// 				},
// 			},
// 			`, field.Name(), requiredOrOptional(field), kindToTerraform(field.MapValue().Kind()))
// 		}

// 		// map with string keys and complex type values
// 		if field.MapValue().Kind() != protoreflect.MessageKind {
// 			return fmt.Sprintf(`"%s": {
// 				Type: schema.TypeMap,
// 				%s,
// 				Elem: %s
// 			},
// 			`, field.Name(), requiredOrOptional(field), g.generateInlineMessageSchema(field.MapValue().Message()))
// 		}
// 	}

// 	if field.IsList() {
// 		// List of complex type objects
// 		if field.Kind() == protoreflect.MessageKind {
// 			return fmt.Sprintf(`"%s": {
// 				Type: schema.TypeList,
// 				%s,
// 				Elem: %s,
// 			},
// 			`, field.Name(), requiredOrOptional(field), g.generateInlineMessageSchema(field.Message()))
// 		}

// 		// List of primitive type objects
// 		return fmt.Sprintf(`"%s": {
// 			Type: schema.TypeList,
// 			%s,
// 			Elem: &schema.Schema{Type: schema.%s},
// 		},
// 		`, field.Name(), requiredOrOptional(field), kindToTerraform(field.Kind()))
// 	}

// 	if field.Kind() == protoreflect.MessageKind {

// 		// check if the field has one of a known mappable types
// 		// and generate an inline schema for types we don't recognize
// 		switch field.Message().FullName() {
// 		case "google.protobuf.Timestamp":
// 			return fmt.Sprintf(`"%s": {
// 				Type: schema.TypeString,
// 				%s,
// 				ValidateFunc: validation.ValidateRFC3339TimeString,
// 			},
// 			`, field.Name(), requiredOrOptional(field))
// 		default:

// 			// fields of custom time are represented as lists in TF schema
// 			return fmt.Sprintf(`"%s": {
// 				Type: schema.TypeList,
// 				%s,
// 				MaxItems:1,
// 				Elem: %s,
// 				},
// 			`, field.Name(), requiredOrOptional(field), g.generateInlineMessageSchema(field.Message()))
// 		}

// 	}

// 	return fmt.Sprintf(`"%s": {
// 		Type: schema.%s,
// 		%s,
// 	},
// 	`, field.Name(), kindToTerraform(field.Kind()), requiredOrOptional(field))

// }
