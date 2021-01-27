package builder

// SchemaMap string -> *Schema
type SchemaMap map[string]*Schema

// NewSchemaMap initializes schema map
func NewSchemaMap() *SchemaMap {
	m := make(SchemaMap)
	return &m
}

// AddSchema add new schema entry
func (m SchemaMap) AddSchema(name string, fullName string) *Schema {
	s := NewSchema(fullName)
	m[name] = s
	return s
}
