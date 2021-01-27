package builder

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Schema extends Terraform Schema with metadata
type Schema struct {
	FullName string

	*schema.Schema
}

func NewSchema(fullName string) *Schema {
	return &Schema{FullName: fullName}
}

func (s *Schema) SetType(t schema.ValueType) *Schema {
	s.Type = t
	return s
}

func (s *Schema) SetRequired(r bool) *Schema {
	s.Required = r
	return s
}
