package cueutils

import (
	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/errors"
	"github.com/hashicorp/go-multierror"
)

// UsefulError returns an error that is concatenated from
// multiple cue errors
func UsefulError(err error) error {
	var usefulError error
	for _, err := range errors.Errors(err) {
		usefulError = multierror.Append(usefulError, err)
	}
	return usefulError
}

// GetAcceptedValues returns the values constraints
// for a cue node
func GetAcceptedValues(node ast.Node) ([]string, error) {
	switch v := node.(type) {
	case *ast.Ident:
		return []string{v.Name}, nil

	case *ast.ListLit:
		return []string{"list"}, nil
	}

	return []string{"None"}, nil
}
