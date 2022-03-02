package cueutils

import (
	"strconv"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/ast"
	"github.com/pterm/pterm"
)

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
	fieldIterator, err := valueIn.Fields(cue.Optional(true))
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

		templateValue := strings.TrimPrefix(templateAttribute.Contents(), `"`)
		templateValue = strings.TrimSuffix(templateValue, `"`)

		switch fieldValue.IncompleteKind() {
		case cue.StringKind:
			valueOut = valueOut.FillPath(fieldValue.Path(), templateValue)

		case cue.IntKind:
			i, err := strconv.Atoi(templateValue)
			if err != nil {
				return valueOut, err
			}
			valueOut = valueOut.FillPath(fieldValue.Path(), i)

		case cue.BoolKind:
			b, err := strconv.ParseBool(templateValue)
			if err != nil {
				return valueOut, err
			}
			valueOut = valueOut.FillPath(fieldValue.Path(), b)

		case cue.ListKind:
			listValue := strings.Split(templateValue, ",")
			valueOut = valueOut.FillPath(fieldValue.Path(), listValue)

		default:
			// Default, just assume string and drop in the value
			pterm.Debug.Println("UNMATCHED", fieldValue.IncompleteKind())
			valueOut = valueOut.FillPath(fieldValue.Path(), templateValue)
		}
	}

	return valueOut, nil
}
