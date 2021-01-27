package builder

// ResourceMap string -> Resource map
type ResourceMap map[string]*Resource

// NewResourceMap initializes schema map
func NewResourceMap() *ResourceMap {
	m := make(ResourceMap)
	return &m
}

// AddResource add new resource map entry
func (m ResourceMap) AddResource(name string, fullName string) *Resource {
	resource := NewResource(name, fullName)
	m[name] = resource
	return resource
}
