package cuedb

import (
	"fmt"
	"testing"

	"cuelang.org/go/cue"
	"github.com/graphql-go/graphql"
	"github.com/stretchr/testify/assert"
)

func TestGraphqlGeneration(t *testing.T) {
	type test struct {
		cueLiteral string
		expected   map[string]graphql.Field
	}

	tests := []test{
		{cueLiteral: "name: string", expected: map[string]graphql.Field{
			"name": {Type: graphql.String},
		}},
		{cueLiteral: "age: int", expected: map[string]graphql.Field{
			"age": {Type: graphql.Int},
		}},
		{cueLiteral: "{ name: string, handles: { network: string, handle: string} }", expected: map[string]graphql.Field{
			"name": {Type: graphql.String},
			"handles": {Type: graphql.NewObject(graphql.ObjectConfig{
				Name: "handles",
				Fields: map[string]graphql.Field{
					"network": {Type: graphql.String},
					"handle":  {Type: graphql.String},
				},
			})},
		}}}

	var cueRuntime cue.Runtime

	for _, tc := range tests {
		cueInstance, err := cueRuntime.Compile("", tc.cueLiteral)
		assert.Equal(t, nil, err)

		cueValue := cueInstance.Value()

		graphQlObject, err := CueValueToGraphQlField(cueValue)
		fmt.Printf("%v", graphQlObject)
		assert.Equal(t, nil, err)

		assert.Equal(t, tc.expected, graphQlObject)
	}
}
