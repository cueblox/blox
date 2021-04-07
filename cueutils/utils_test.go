package cueutils

import (
	"errors"
	"testing"

	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/parser"
	"github.com/stretchr/testify/assert"
)

func getAstNodeValue(cue string, label string) (ast.Node, error) {
	lAst, err := parser.ParseFile("", cue)

	if err != nil {
		return nil, err
	}

	for _, decl := range lAst.Decls {
		if field, ok := decl.(*ast.Field); ok {
			if fieldName, ok := field.Label.(*ast.Ident); ok {
				if label == fieldName.Name {
					return field.Value, nil
				}
			}
		}
	}

	return nil, errors.New("Couldn't find field with label")
}

func TestGetAcceptedValuesString(t *testing.T) {
	node, err := getAstNodeValue("field_name: string", "field_name")
	if err != nil {
		t.Fatal(err)
	}

	acceptedValues, err := GetAcceptedValues(node)

	if err != nil {
		t.Fatal(err)
	}

	assert.ElementsMatch(t, []string{"string"}, acceptedValues)
}

func TestGetAcceptedValuesInt(t *testing.T) {
	node, err := getAstNodeValue("field_name: int", "field_name")
	if err != nil {
		t.Fatal(err)
	}

	acceptedValues, err := GetAcceptedValues(node)

	if err != nil {
		t.Fatal(err)
	}

	assert.ElementsMatch(t, []string{"int"}, acceptedValues)
}

func TestGetAcceptedValuesFloat(t *testing.T) {
	node, err := getAstNodeValue("field_name: float", "field_name")
	if err != nil {
		t.Fatal(err)
	}

	acceptedValues, err := GetAcceptedValues(node)

	if err != nil {
		t.Fatal(err)
	}

	assert.ElementsMatch(t, []string{"float"}, acceptedValues)
}

func TestGetAcceptedValuesNumber(t *testing.T) {
	node, err := getAstNodeValue("field_name: number", "field_name")
	if err != nil {
		t.Fatal(err)
	}

	acceptedValues, err := GetAcceptedValues(node)

	if err != nil {
		t.Fatal(err)
	}

	assert.ElementsMatch(t, []string{"number"}, acceptedValues)
}

func TestGetAcceptedValuesBool(t *testing.T) {
	node, err := getAstNodeValue("field_name: bool", "field_name")
	if err != nil {
		t.Fatal(err)
	}

	acceptedValues, err := GetAcceptedValues(node)

	if err != nil {
		t.Fatal(err)
	}

	assert.ElementsMatch(t, []string{"bool"}, acceptedValues)
}

func TestGetAcceptedValuesList(t *testing.T) {
	node, err := getAstNodeValue("field_name: [string]", "field_name")
	if err != nil {
		t.Fatal(err)
	}

	acceptedValues, err := GetAcceptedValues(node)

	if err != nil {
		t.Fatal(err)
	}

	assert.ElementsMatch(t, []string{"list"}, acceptedValues)
}
