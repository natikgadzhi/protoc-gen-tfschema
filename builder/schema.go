package builder

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type Schema struct {
	schema schema.Schema
}

func (s *Schema) SetType(t schema.ValueType) *Schema {
	s.schema.Type = t
	return s
}

func (s *Schema) SetRequired(r bool) *Schema {
	s.schema.Required = r
	return s
}
