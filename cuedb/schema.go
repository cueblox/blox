package cuedb

import (
	"errors"

	"cuelang.org/go/cue"
)

// GET_VERSION_CUE is a constant?
// that is added to cue models
var GET_VERSION_CUE = `
{
	_schema: {
		version: string
	}
}
`

// GetSchemaVersion will attempt to pull a "version" out of the
// schema's metadata, returning an error if it can't
func GetSchemaVersion(schema cue.Value) (string, error) {
	var cueRuntime cue.Runtime
	cueInstance, err := cueRuntime.Compile("validateVersionCue", GET_VERSION_CUE)
	if err != nil {
		return "", err
	}
	versionedSchema := schema.Unify(cueInstance.Value())
	if err = versionedSchema.Validate(); err != nil {
		return "", err
	}

	fields, err := versionedSchema.Fields(cue.All())
	if err != nil {
		return "", err
	}

	for fields.Next() {
		if fields.Label() == "_schema" {
			schemaValue := fields.Value()

			versionField, err := schemaValue.LookupField("version")
			if err != nil {
				return "", err
			}

			stringVersion, err := versionField.Value.String()
			if err != nil {
				return "", err
			}

			return stringVersion, nil
		}
	}

	return "", nil
}

// SchemaV1Metadata stores information about the schema
// that is shared
type SchemaV1Metadata struct {
	Namespace string
	Name      string
}

// SCHEMA_V1_METADATA is the cue representation of
// SchemaV1Metadata
const SCHEMA_V1_METADATA = `
{
	_schema: {
		namespace: string
		name: string
	}
}
`

// GetSchemaV1Metadata returns the metadata for a cue.Value
func GetSchemaV1Metadata(schema cue.Value) (SchemaV1Metadata, error) {
	var cueRuntime cue.Runtime
	cueInstance, err := cueRuntime.Compile("schemav1", SCHEMA_V1_METADATA)
	if err != nil {
		return SchemaV1Metadata{}, err
	}
	versionedSchema := schema.Unify(cueInstance.Value())
	if err = versionedSchema.Validate(); err != nil {
		return SchemaV1Metadata{}, err
	}

	fields, err := versionedSchema.Fields(cue.All())
	if err != nil {
		return SchemaV1Metadata{}, err
	}

	schemaV1 := SchemaV1Metadata{}

	for fields.Next() {
		if fields.Label() == "_schema" {
			schemaValue := fields.Value()

			err := schemaValue.Decode(&schemaV1)
			if err != nil {
				return SchemaV1Metadata{}, err
			}

			return schemaV1, nil
		}
	}

	return SchemaV1Metadata{}, errors.New("couldn't get SchemaV1Metadata")
}

// V1_MODEL is the cue representation of
// Model metadata
const V1_MODEL = `
{
	_model: {
		plural: string
	}
}
`

// GetV1Model returns the V1_MODEL information from a cue.Value
func GetV1Model(schema cue.Value) (string, error) {
	var cueRuntime cue.Runtime
	cueInstance, err := cueRuntime.Compile("", V1_MODEL)
	if err != nil {
		return "", err
	}
	modelSchema := schema.Unify(cueInstance.Value())
	if err = modelSchema.Validate(); err != nil {
		return "", err
	}

	fields, err := modelSchema.Fields(cue.All())
	if err != nil {
		return "", err
	}

	for fields.Next() {
		if fields.Label() == "_model" {
			schemaValue := fields.Value()

			pluralField, err := schemaValue.LookupField("plural")
			if err != nil {
				return "", err
			}

			plural, err := pluralField.Value.String()
			if err != nil {
				return "", err
			}

			return plural, nil
		}
	}

	return "", nil
}
