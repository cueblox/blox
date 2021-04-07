package schema

// Schema is a definition of a set of
// related models
type Schema struct {
	Namespace string
	Name      string
	Versions  []*Version
}
