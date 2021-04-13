package cueutils

import (
	"fmt"
	"strconv"
	"strings"

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

func CreateFromTemplate(valueOut cue.Value, valueIn cue.Value, interactive bool) (cue.Value, error) {
	fieldIterator, err := valueIn.Fields(cue.Optional(true))
	if err != nil {
		return valueOut, err
	}

	for fieldIterator.Next() {
		fieldValue := fieldIterator.Value()

		if cue.StructKind == fieldValue.IncompleteKind() {
			valueOut, err = CreateFromTemplate(valueOut, fieldValue, interactive)
			if err != nil {
				return valueOut, err
			}
			continue
		}

		comments := fieldValue.Doc()

		templateAttribute := fieldValue.Attribute("template")
		if err = templateAttribute.Err(); err != nil {
			// For now, we just skip
			continue
		}

		templateValue := strings.TrimPrefix(templateAttribute.Contents(), `"`)
		templateValue = strings.TrimSuffix(templateValue, `"`)

		switch fieldValue.IncompleteKind() {
		case cue.StringKind:
			if interactive {
				valueOut = valueOut.FillPath(fieldValue.Path(), scanString(fieldIterator.Label(), comments))
			} else {
				valueOut = valueOut.FillPath(fieldValue.Path(), templateValue)
			}

		case cue.IntKind:
			var i int
			if interactive {
				i = scanInt(fieldIterator.Label(), comments)
			} else {
				i, err = strconv.Atoi(templateValue)
				if err != nil {
					return valueOut, err
				}
			}
			valueOut = valueOut.FillPath(fieldValue.Path(), i)

		case cue.BoolKind:
			var b bool
			if interactive {
				b = scanBool(fieldIterator.Label(), comments)
			} else {
				b, err = strconv.ParseBool(templateValue)
				if err != nil {
					return valueOut, err
				}
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

	fmt.Println(valueOut)
	return valueOut, nil
}

func scanString(field string, comments []*ast.CommentGroup) string {
	var str string

	for _, comment := range comments {
		pterm.Description.Println(comment.Text())
	}

	pterm.BgCyan.Printf("%s: ", field)

	_, err := fmt.Scanln(&str)
	if err != nil {
		return ""
	}

	return str
}

func scanInt(field string, comments []*ast.CommentGroup) int {
	var i int

	for _, comment := range comments {
		pterm.Description.Println(comment.Text())
	}

	pterm.BgCyan.Printf("%s: ", field)

	_, err := fmt.Scanln(&i)
	if err != nil {
		return 0
	}

	return i
}

func scanBool(field string, comments []*ast.CommentGroup) bool {
	var b bool

	for _, comment := range comments {
		pterm.Description.Println(comment.Text())
	}

	pterm.BgCyan.Printf("%s: ", field)

	_, err := fmt.Scanln(&b)
	if err != nil {
		return false
	}

	return b
}
