package reducer

import (
	"github.com/nategadzhi/protoc-gen-tfschema/builder"
)

// Reduce removes nested resources with a single field:
//
// Example:
//
// "allow": {
// 	 Type:     schema.TypeList,
//   Optional: true,
// 	 MaxItems: 1,
// 	 Elem: &schema.Resource{
// 		Schema: SchemaBoolValue(), // Which has the single field value bool
// 	 },
// },
//
// is replaced with
//
// "allow": {
//   Type: schema.TypeBool,
//   Optional: true,
// }
func Reduce(m builder.ResourceMap) {

	for name, value := range m {
		// Resource has a signle field
		if len(value.Schema) == 1 {
			keys := value.Schema.Keys()
			replacementValue := value.Schema[keys[0]]

			replace(m, name, replacementValue)
			delete(m, name) // Reference to this resource is no longer needed
		}
	}
}

// Scans all resources and replaces reference to nested resource with single nested field itself
func replace(m builder.ResourceMap, nameToReplace string, replacementValue *builder.Schema) {
	for _, resource := range m {
		for name, schema := range resource.Schema {
			// This schema field references nested resource in question
			if schema.Elem != nil && schema.Elem == nameToReplace {
				// Replace it with nested field definition
				resource.Schema[name] = replacementValue
			}
		}
	}
}
