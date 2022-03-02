package cuedb

import (
	"testing"

	"cuelang.org/go/cue/cuecontext"
	"github.com/graphql-go/graphql"
	"github.com/stretchr/testify/assert"
)

func TestGraphqlGeneration(t *testing.T) {
	type test struct {
		cueLiteral string
		dataSet    *DataSet
		expected   graphql.Fields
	}

	tests := []test{
		{
			dataSet:    &DataSet{name: "Type"},
			cueLiteral: "t1: string", expected: graphql.Fields{
				"t1": {Type: &graphql.NonNull{OfType: graphql.String}},
			},
		},
		{
			dataSet:    &DataSet{name: "Type"},
			cueLiteral: "t2: int", expected: graphql.Fields{
				"t2": {Type: &graphql.NonNull{OfType: graphql.Int}},
			},
		},
		{
			dataSet:    &DataSet{name: "Type"},
			cueLiteral: "t3: [...int]", expected: graphql.Fields{
				"t3": {Type: &graphql.List{OfType: graphql.Int}},
			},
		},
		{
			dataSet:    &DataSet{name: "Type"},
			cueLiteral: "t4: [...string]", expected: graphql.Fields{
				"t4": {Type: &graphql.List{OfType: graphql.String}},
			},
		},
		{
			dataSet:    &DataSet{name: "Type"},
			cueLiteral: "t5: string, t6?: int", expected: graphql.Fields{
				"t5": {Type: &graphql.NonNull{OfType: graphql.String}},
				"t6": {Type: graphql.Int},
			},
		},
		{
			dataSet:    &DataSet{name: "Type"},
			cueLiteral: "{ t7: string, t8: { t9: string, t10: string} }", expected: graphql.Fields{
				"t7": {Type: &graphql.NonNull{OfType: graphql.String}},
				"t8": {Type: graphql.NewObject(graphql.ObjectConfig{
					Name: "TypeT8",
					Fields: graphql.Fields{
						"t9":  {Type: &graphql.NonNull{OfType: graphql.String}},
						"t10": {Type: &graphql.NonNull{OfType: graphql.String}},
					},
				})},
			},
		},
		{
			dataSet:    &DataSet{name: "Type"},
			cueLiteral: "{ t11: string, t12: { t13: string, t14: string, t15: { t16: int } } }", expected: graphql.Fields{
				"t11": {Type: &graphql.NonNull{OfType: graphql.String}},
				"t12": {Type: graphql.NewObject(graphql.ObjectConfig{
					Name: "TypeT12",
					Fields: graphql.Fields{
						"t13": {Type: &graphql.NonNull{OfType: graphql.String}},
						"t14": {Type: &graphql.NonNull{OfType: graphql.String}},
						"t15": {Type: graphql.NewObject(graphql.ObjectConfig{
							Name: "TypeT15",
							Fields: graphql.Fields{
								"t16": {Type: &graphql.NonNull{OfType: graphql.Int}},
							},
						})},
					},
				})},
			},
		},
		{
			dataSet:    &DataSet{name: "Type"},
			cueLiteral: "{ #Test: { t19: string}\nt17: string, t18: #Test }", expected: graphql.Fields{
				"t17": {Type: &graphql.NonNull{OfType: graphql.String}},
				"t18": {Type: graphql.NewObject(graphql.ObjectConfig{
					Name: "TypeT18",
					Fields: graphql.Fields{
						"t19": {Type: &graphql.NonNull{OfType: graphql.String}},
					},
				})},
			},
		},
		{
			dataSet:    &DataSet{name: "Type"},
			cueLiteral: "{ #Test: { t20: string}\n t21: string, t22: [ ... #Test ] }", expected: graphql.Fields{
				"t21": {Type: &graphql.NonNull{OfType: graphql.String}},
				"t22": {Type: &graphql.List{OfType: graphql.NewObject(graphql.ObjectConfig{
					Name: "TypeT22",
					Fields: graphql.Fields{
						"t20": {Type: &graphql.NonNull{OfType: graphql.String}},
					},
				})}},
			},
		},
		// WIP
		// Aim is to "flatten" disjunctions into a single struct
		// {cueLiteral: "{ #A: {t23: string}\n#B: { t24?: string}\n t25: string, t26: [ ... #A | #B ] }", expected: graphql.Fields{
		// 	"t25": {Type: &graphql.NonNull{OfType: graphql.String}},
		// 	"t26": {Type: &graphql.List{OfType: graphql.NewObject(graphql.ObjectConfig{
		// 		Fields: graphql.Fields{
		// 			"t23": {Type: &graphql.NonNull{OfType: graphql.String}},
		// 			"t24": {Type: graphql.String},
		// 		},
		// 	})}},
		// }},
	}

	cueContext := cuecontext.New()

	for _, tc := range tests {
		cueValue := cueContext.CompileString(tc.cueLiteral)
		assert.Equal(t, nil, cueValue.Err())

		graphqlObjects := make(map[string]GraphQlObjectGlue)

		graphQlObject, err := CueValueToGraphQlField(graphqlObjects, *tc.dataSet, cueValue)
		assert.Equal(t, nil, err)
		assert.EqualValues(t, tc.expected, graphQlObject)
	}
}
