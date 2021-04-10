package repository

// Version represents a version of a schema
type Version struct {
	Namespace  string
	Name       string
	Schema     string
	Definition string
}
