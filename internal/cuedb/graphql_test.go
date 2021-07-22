package cuedb

import (
	"testing"

	"cuelang.org/go/cue"
	"github.com/graphql-go/graphql"
	"github.com/stretchr/testify/assert"
)

func TestGraphqlGeneration(t *testing.T) {
	type test struct {
		cueLiteral string
		expected   map[string]*graphql.Field
	}

	tests := []test{
		{cueLiteral: "t1: string", expected: map[string]*graphql.Field{
			"t1": {Type: &graphql.NonNull{OfType: graphql.String}},
		}},
		{cueLiteral: "t2: int", expected: map[string]*graphql.Field{
			"t2": {Type: &graphql.NonNull{OfType: graphql.Int}},
		}},
		{cueLiteral: "t3: [...int]", expected: map[string]*graphql.Field{
			"t3": {Type: &graphql.List{OfType: graphql.Int}},
		}},
		{cueLiteral: "t4: [...string]", expected: map[string]*graphql.Field{
			"t4": {Type: &graphql.List{OfType: graphql.String}},
		}},
		{cueLiteral: "t5: string, t6?: int", expected: map[string]*graphql.Field{
			"t5": {Type: &graphql.NonNull{OfType: graphql.String}},
			"t6": {Type: graphql.Int},
		}},
		{cueLiteral: "{ t7: string, t8: { t9: string, t10: string} }", expected: map[string]*graphql.Field{
			"t7": {Type: &graphql.NonNull{OfType: graphql.String}},
			"t8": {Type: graphql.NewObject(graphql.ObjectConfig{
				Name: "t8",
				Fields: map[string]*graphql.Field{
					"t9":  {Type: &graphql.NonNull{OfType: graphql.String}},
					"t10": {Type: &graphql.NonNull{OfType: graphql.String}},
				},
			})},
		}},
		{cueLiteral: "{ t11: string, t12: { t13: string, t14: string, t15: { t16: int } } }", expected: map[string]*graphql.Field{
			"t11": {Type: &graphql.NonNull{OfType: graphql.String}},
			"t12": {Type: graphql.NewObject(graphql.ObjectConfig{
				Name: "t12",
				Fields: map[string]*graphql.Field{
					"t13": {Type: &graphql.NonNull{OfType: graphql.String}},
					"t14": {Type: &graphql.NonNull{OfType: graphql.String}},
					"t15": {Type: graphql.NewObject(graphql.ObjectConfig{
						Name: "t15",
						Fields: map[string]*graphql.Field{
							"t16": {Type: &graphql.NonNull{OfType: graphql.Int}},
						},
					})},
				},
			})},
		}},
	}

	var cueRuntime cue.Runtime

	for _, tc := range tests {
		cueInstance, err := cueRuntime.Compile("", tc.cueLiteral)
		assert.Equal(t, nil, err)

		cueValue := cueInstance.Value()

		graphqlObjects := make(map[string]GraphQlObjectGlue)

		graphQlObject, err := CueValueToGraphQlField(graphqlObjects, cueValue)
		assert.Equal(t, nil, err)
		assert.EqualValues(t, tc.expected, graphQlObject)
	}
}
