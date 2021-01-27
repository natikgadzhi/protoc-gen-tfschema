package builder

type schemaMap map[string]*Schema

type SchemaMap struct {
	schemaMap schemaMap
}

func NewSchemaMap() SchemaMap {
	return SchemaMap{schemaMap: make(schemaMap)}
}

// func AddSchema(name string) {
// 	return Schema
// }
