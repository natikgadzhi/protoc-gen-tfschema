package reducer

import (
	"github.com/nategadzhi/protoc-gen-tfschema/builder"
)

// Reduce removes nested resources with a single field
func Reduce(m builder.ResourceMap) {

	for name, value := range m {
		if len(value.Schema) == 1 {
			keys := value.Schema.Keys()
			replacementValue := value.Schema[keys[0]]

			replace(m, name, replacementValue)
			delete(m, name)
		}
	}
}

func replace(m builder.ResourceMap, replacementName string, replacementValue *builder.Schema) {
	for _, resource := range m {
		for name, schema := range resource.Schema {
			if schema.Elem != nil && schema.Elem == replacementName {
				resource.Schema[name] = replacementValue
			}
		}
	}
}
