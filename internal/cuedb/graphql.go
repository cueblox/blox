package cuedb

import (
	"fmt"
	"strings"

	"cuelang.org/go/cue"
	"github.com/graphql-go/graphql"
)

type GraphQlObjectGlue struct {
	Object   *graphql.Object
	Engine   *Engine
	Resolver func(p graphql.ResolveParams) (interface{}, error)
}

func CueValueToGraphQlField(existingObjects map[string]GraphQlObjectGlue, cueValue cue.Value) (graphql.Fields, error) {
	fields, err := cueValue.Fields(cue.All())
	if err != nil {
		return nil, err
	}

	graphQlFields := make(map[string]*graphql.Field)

	for fields.Next() {
		if strings.HasPrefix(fields.Label(), "_") {
			continue
		}

		if fields.Selector().IsDefinition() {
			continue
		}

		relationshipLabel := fields.Label()
		relationship := fields.Value().Attribute("relationship")

		switch fields.Value().IncompleteKind() {
		case cue.StructKind:
			subFields, err := CueValueToGraphQlField(existingObjects, fields.Value())
			if err != nil {
				return nil, err
			}

			graphQlFields[fields.Label()] = &graphql.Field{
				Type: graphql.NewObject(graphql.ObjectConfig{
					Fields: subFields,
				}),
			}

		case cue.ListKind:
			listOf := fields.Value().LookupPath(cue.MakePath(cue.AnyIndex))

			kind, err := CueValueToGraphQlType(listOf)
			if err == nil {
				if err = relationship.Err(); err == nil {
					graphQlFields[fields.Label()] = &graphql.Field{
						Type: graphql.NewList(existingObjects[relationship.Contents()].Object),
						Resolve: func(p graphql.ResolveParams) (interface{}, error) {
							data := existingObjects[relationship.Contents()].Engine.GetAllData(fmt.Sprintf("#%s", relationship.Contents()))

							records := make(map[string]interface{})
							if err = data.Decode(&records); err != nil {
								return nil, err
							}

							source, ok := p.Source.(map[string]interface{})

							if !ok {
								return nil, nil
							}

							searchIds := source[relationshipLabel].([]interface{})
							returnRecords := []interface{}{}

							for recordID, record := range records {
								for _, searchID := range searchIds {
									if string(recordID) == searchID.(string) {
										returnRecords = append(returnRecords, record)
									}
								}
							}

							return returnRecords, nil
						},
					}
				} else {
					graphQlFields[fields.Label()] = &graphql.Field{
						Type: &graphql.List{
							OfType: kind,
						},
					}
				}
				continue
			}

			// List of non-scalar types
			subFields, err := CueValueToGraphQlField(existingObjects, listOf.Value())

			// No error, so we know this is a simple value or struct
			if err == nil {
				graphQlFields[fields.Label()] = &graphql.Field{
					Type: &graphql.List{OfType: graphql.NewObject(graphql.ObjectConfig{
						Name:   fields.Label(),
						Fields: subFields,
					})},
				}
				continue
			}

			// Error, probably a disjunction or other complex value
			// Not handled, yet.
			// switch listOf.IncompleteKind() {
			// case cue.StructKind:
			// 	fields, err := listOf.Fields()
			// 	fmt.Println(err)
			// 	fmt.Println(fields)
			// }

			return nil, err

		case cue.BoolKind, cue.FloatKind, cue.IntKind, cue.NumberKind, cue.StringKind:
			kind, err := CueValueToGraphQlType(fields.Value())
			if err != nil {
				return nil, err
			}

			if err = relationship.Err(); err == nil {
				graphQlFields[fields.Label()] = &graphql.Field{
					Type: existingObjects[relationship.Contents()].Object,
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						data := existingObjects[relationship.Contents()].Engine.GetAllData(fmt.Sprintf("#%s", relationship.Contents()))

						records := make(map[string]interface{})
						if err = data.Decode(&records); err != nil {
							return nil, err
						}

						source, ok := p.Source.(map[string]interface{})

						if !ok {
							return nil, nil
						}

						for recordID, record := range records {
							if string(recordID) == source[relationshipLabel].(string) {
								return record, nil
							}
						}

						return nil, nil
					},
				}
			} else if fields.IsOptional() {
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
