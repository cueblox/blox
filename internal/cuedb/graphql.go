package cuedb

import (
	"cuelang.org/go/cue"
	"github.com/graphql-go/graphql"
)

func CueValueToGraphQlField(cueValue cue.Value) (map[string]graphql.Field, error) {
	fields, err := cueValue.Fields()
	if err != nil {
		return nil, err
	}

	graphQlFields := make(map[string]graphql.Field)

	for fields.Next() {
		switch fields.Value().IncompleteKind() {
		case cue.StructKind:
			subFields, err := CueValueToGraphQlField(fields.Value())
			if err != nil {
				return nil, err
			}

			graphQlFields[fields.Label()] = graphql.Field{
				Type: graphql.NewObject(graphql.ObjectConfig{
					Name:   fields.Label(),
					Fields: subFields,
				})}

		case cue.IntKind:
			graphQlFields[fields.Label()] = graphql.Field{Type: graphql.Int}

		case cue.StringKind:
			graphQlFields[fields.Label()] = graphql.Field{Type: graphql.String}
		}
	}

	return graphQlFields, nil
}
