package cueutils

import (
	"errors"
	"testing"

	"cuelang.org/go/cue"
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

func TestCreateFromTemplate(t *testing.T) {
	cueWithTemplateAttributes := `{
		// NameComment
		name: string @template(Random Name) //NameInlineComments
		age: int @template(21)
		happy: bool @template(true)
		scottish: bool @template(false)
		movies: [...string] @template(The Matrix, Top Gun)
		social: {
			network: string @template(twitter)
			name: string @template(rawkode)
		}
	}`

	var cueRuntime cue.Runtime
	cueInstance, err := cueRuntime.Compile("", cueWithTemplateAttributes)
	assert.Equal(t, nil, err)

	cueValue := cueInstance.Value()

	cueTemplate, err := CreateFromTemplate(cueValue, cueValue)
	assert.Equal(t, nil, err)

	name, err := cueTemplate.LookupPath(cue.ParsePath("name")).String()
	assert.Equal(t, nil, err)
	assert.Equal(t, "Random Name", name)

	age, err := cueTemplate.LookupPath(cue.ParsePath("age")).Int64()
	assert.Equal(t, nil, err)
	assert.Equal(t, int64(21), age)

	happy, err := cueTemplate.LookupPath(cue.ParsePath("happy")).Bool()
	assert.Equal(t, nil, err)
	assert.Equal(t, true, happy)

	scottish, err := cueTemplate.LookupPath(cue.ParsePath("scottish")).Bool()
	assert.Equal(t, nil, err)
	assert.Equal(t, false, scottish)

	social := cueTemplate.LookupPath(cue.ParsePath("social"))

	socialNetwork, err := social.LookupPath(cue.ParsePath("network")).String()
	assert.Equal(t, nil, err)
	assert.Equal(t, "twitter", socialNetwork)

	socialName, err := social.LookupPath(cue.ParsePath("name")).String()
	assert.Equal(t, nil, err)
	assert.Equal(t, "rawkode", socialName)

	movies, err := cueTemplate.LookupPath(cue.ParsePath("movies")).List()
	expectedMovies := []string{"The Matrix"}
	assert.Equal(t, nil, err)
	for _, expectedMovie := range expectedMovies {
		movies.Next()
		movieString, err := movies.Value().String()
		assert.Equal(t, nil, err)
		assert.Equal(t, expectedMovie, movieString)
	}
}
