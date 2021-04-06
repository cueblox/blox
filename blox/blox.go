package blox

import (
	// import go:embed
	_ "embed"
)

// Config stores information about
// the repository location
type Config struct {
	RepositoryRoot  string `json:"repository_root"`
	Namespace       string `json:"namespace"`
	OutputDirectory string `json:"output_dir"`
}

//go:embed schema.cue
var schemaCue []byte
