package cuedb

import (
	"cuelang.org/go/cue"
	"github.com/graphql-go/graphql"
)

func CueValueToGraphQlField(cueValue cue.Value) (map[string]*graphql.Field, error) {
	fields, err := cueValue.Fields()
	if err != nil {
		return nil, err
	}

	graphQlFields := make(map[string]*graphql.Field)

	for fields.Next() {
		switch fields.Value().IncompleteKind() {
		case cue.StructKind:
			subFields, err := CueValueToGraphQlField(fields.Value())
			if err != nil {
				return nil, err
			}

			graphQlFields[fields.Label()] = &graphql.Field{
				Type: graphql.NewObject(graphql.ObjectConfig{
					Name:   fields.Label(),
					Fields: subFields,
				})}

		case cue.ListKind:
			// listValues, err := fields.Value().
			// if err != nil {
			// 	return nil, err
			// }

			// fmt.Printf("11: %v\n\n", listValues)
			// fmt.Printf("11: %v\n\n", listValues.)
			// // Useless
			// fmt.Printf("12: %v\n\n", listValues.Value())
			// fmt.Printf("13: %v\n\n", listValues.Label())
			// fmt.Printf("14: %v\n\n", listValues.IsOptional())

			// for listValues.Next() {
			// 	fmt.Printf("\n\n2.x: %v\n", listValues.Value().IncompleteKind())
			// }

			// graphQlFields[fields.Label()] = graphql.Field{Type: &graphql.List{
			// 	OfType: graphql.String,
			// }}

		case cue.BoolKind:
			graphQlFields[fields.Label()] = &graphql.Field{Type: graphql.Boolean}

		case cue.IntKind:
			graphQlFields[fields.Label()] = &graphql.Field{Type: graphql.Int}

		case cue.StringKind:
			graphQlFields[fields.Label()] = &graphql.Field{Type: graphql.String}
		}
	}

	return graphQlFields, nil
}
