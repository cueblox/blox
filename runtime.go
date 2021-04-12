package blox

import (
	"cuelang.org/go/cue"
)

type Runtime struct {
	CueRuntime *cue.Runtime
	Database   cue.Value
}

// NewRuntime creates a new runtime engine
func NewRuntime() (*Runtime, error) {
	var cueRuntime cue.Runtime
	cueInstance, err := cueRuntime.Compile("", "")

	if nil != err {
		return &Runtime{}, err
	}

	runtime := &Runtime{
		CueRuntime: &cueRuntime,
		Database:   cueInstance.Value(),
	}

	return runtime, nil
}

// NewRuntimeWithBase creates a new runtime engine
// with the cue provided in `base` as the initial cue values
func NewRuntimeWithBase(base string) (*Runtime, error) {
	var cueRuntime cue.Runtime
	cueInstance, err := cueRuntime.Compile("", base)

	if nil != err {
		return &Runtime{}, err
	}

	runtime := &Runtime{
		CueRuntime: &cueRuntime,
		Database:   cueInstance.Value(),
	}

	return runtime, nil
}
