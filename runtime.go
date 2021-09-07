package blox

import (
	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
)

type Runtime struct {
	CueContext *cue.Context
	Database   cue.Value
}

// NewRuntime creates a new runtime engine
func NewRuntime() (*Runtime, error) {
	cueContext := cuecontext.New()
	cueValue := cueContext.CompileString("")

	if cueValue.Err() != nil {
		return nil, cueValue.Err()
	}

	runtime := &Runtime{
		CueContext: cueContext,
		Database:   cueValue,
	}

	return runtime, nil
}

// NewRuntimeWithBase creates a new runtime engine
// with the cue provided in `base` as the initial cue values
func NewRuntimeWithBase(base string) (*Runtime, error) {
	cueContext := cuecontext.New()
	cueValue := cueContext.CompileString(base)

	if cueValue.Err() != nil {
		return nil, cueValue.Err()
	}

	runtime := &Runtime{
		CueContext: cueContext,
		Database:   cueValue,
	}

	return runtime, nil
}
