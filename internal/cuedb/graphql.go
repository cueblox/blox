package cuedb

import (
	"fmt"

	"cuelang.org/go/cue"
	"github.com/graphql-go/graphql"
)

func CueValueToGraphQlField(cueValue cue.Value) (map[string]*graphql.Field, error) {
	fields, err := cueValue.Fields(cue.Optional(true))
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
			kind, err := CueValueToGraphQlType(fields.Value().LookupPath(cue.MakePath(cue.AnyIndex)))
			if err != nil {
				return nil, err
			}
			graphQlFields[fields.Label()] = &graphql.Field{
				Type: &graphql.List{
					OfType: kind,
				},
			}

		case cue.BoolKind:
		case cue.FloatKind:
		case cue.IntKind:
		case cue.NumberKind:
		case cue.StringKind:
			kind, err := CueValueToGraphQlType(fields.Value())
			if err != nil {
				return nil, err
			}
			if fields.IsOptional() {
				graphQlFields[fields.Label()] = &graphql.Field{
					Type: kind,
				}
			} else {
				graphQlFields[fields.Label()] = &graphql.Field{
					Type: &graphql.NonNull{
						OfType: kind,
					},
				}
			}
		}
	}

	return graphQlFields, nil
}

func CueValueToGraphQlType(value cue.Value) (*graphql.Scalar, error) {
	switch value.IncompleteKind() {
	case cue.BoolKind:
		return graphql.Boolean, nil
	case cue.FloatKind:
		return graphql.Float, nil
	case cue.IntKind:
		return graphql.Int, nil
	case cue.NumberKind:
		return graphql.Float, nil
	case cue.StringKind:
		return graphql.String, nil
	}

	return nil, fmt.Errorf("unhandled type: %v", value.IncompleteKind())
}
