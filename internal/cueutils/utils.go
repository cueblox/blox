package cueutils

import (
	"strconv"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/errors"
	"github.com/hashicorp/go-multierror"
	"github.com/pterm/pterm"
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

func CreateFromTemplate(valueOut cue.Value, valueIn cue.Value) (cue.Value, error) {
	fieldIterator, err := valueIn.Fields()
	if err != nil {
		return valueOut, err
	}

	for fieldIterator.Next() {
		fieldValue := fieldIterator.Value()

		if cue.StructKind == fieldValue.IncompleteKind() {
			valueOut, err = CreateFromTemplate(valueOut, fieldValue)
			if err != nil {
				return valueOut, err
			}
			continue
		}

		templateAttribute := fieldValue.Attribute("template")
		if err = templateAttribute.Err(); err != nil {
			// For now, we just skip
			continue
		}

		templateValue := templateAttribute.Contents()

		switch fieldValue.IncompleteKind() {
		case cue.StringKind:
			valueOut = valueOut.FillPath(fieldValue.Path(), templateValue)

		case cue.IntKind:
			i, err := strconv.Atoi(templateValue)
			if err != nil {
				return valueOut, err
			}
			valueOut = valueOut.FillPath(fieldValue.Path(), i)

		default:
			// Default, just assume string and drop in the value
			pterm.Debug.Println("UNMATCHED", fieldValue.IncompleteKind())
			valueOut = valueOut.FillPath(fieldValue.Path(), templateValue)
		}
	}

	return valueOut, nil
}

// func T(cueValue cue.Value) (cue.Value, error) {
// 	fieldsIterator, err := cueValue.Fields()
// 	if err != nil {
// 		return cue.Value{}, err
// 	}

// 	var cueRuntime cue.Runtime
// 	cueInstance, err := cueRuntime.Compile("", "")
// 	if err != nil {
// 		return cueValue, err
// 	}

// 	templateValue := cueInstance.Value()

// 	for fieldsIterator.Next() {
// 		fieldLabel := fieldsIterator.Label()
// 		fieldValue := fieldsIterator.Value()

// 		fmt.Println(fieldLabel)
// 		switch fieldValue.IncompleteKind() {
// 		case cue.StringKind:
// 			// Attributes are strings, no cast needed
// 			return
// 		default:
// 			fmt.Println(fieldValue.IncompleteKind())
// 		}
// 	}

// 	// 	// Do we have a @template attribute?
// 	// 	fieldValue := fieldsIterator.Value()
// 	// 	templateAttribute := fieldValue.Attribute("template")

// 	// 	if err = templateAttribute.Err(); err != nil {
// 	// 		// For now, we just skip
// 	// 		continue
// 	// 	}

// 	// 	a, b := fieldValue.Expr()
// 	// 	switch a {
// 	// 	// I believe this means we should only have a single value
// 	// 	case cue.NoOp:
// 	// 		singleValue := b[0]
// 	// 		singleValue.
// 	// 		switch singleValue.Kind() {
// 	// 		case cue.StringKind:
// 	// 			fmt.Println("STRING")
// 	// 		default:
// 	// 			fmt.Println(singleValue.Kind())
// 	// 		}
// 	// 	case cue.OrOp:
// 	// 		fmt.Println("OR")
// 	// 	case cue.AndOp:
// 	// 		fmt.Println("AND")
// 	// 	default:
// 	// 		fmt.Println("ILLEGAL", a.Token().String())
// 	// 	}
// 	// 	fmt.Println(b)
// 	// node := fieldsIterator.Value().Source()

// 	// switch v := node.(type) {
// 	// case *ast.Ident:
// 	// 	fmt.Println("IDENT", v.Name)
// 	// default:
// 	// 	fmt.Println("NODE", v)
// 	// }

// 	// switch fieldsIterator.Value().Kind() {
// 	// case cue.StringKind:
// 	// 	fmt.Println("STRING BITCHES")
// 	// case cue.IntKind:
// 	// 	fmt.Println("INT")
// 	// default:
// 	// 	fmt.Println(fieldValue.Kind())
// 	// }

// 	// templateValue = templateValue.FillPath(cue.ParsePath(fieldsIterator.Label()), fieldValue)
// 	return templateValue, nil
// }
