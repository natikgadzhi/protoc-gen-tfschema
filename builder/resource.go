package builder

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

// Resource extends Terraform resource with metadata required for file generation
type Resource struct {
	Name     string
	FullName string
	Schema   *SchemaMap // Overrrides schema.Resource.Schema

	*schema.Resource
}

// NewResource creates resource
func NewResource(name string, fullName string) *Resource {
	return &Resource{Name: name, FullName: name, Schema: NewSchemaMap()}
}

// AddSchema adds new schema to Schema
func (r *Resource) AddSchema(name string, fullName string) *Schema {
	return r.AddSchema(name, fullName)
}
