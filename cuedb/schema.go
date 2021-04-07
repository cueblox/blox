package cuedb

import (
	"errors"
	"fmt"

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

// SchemaV1MetadataCue is the cue representation of
// SchemaV1Metadata
const SchemaV1MetadataCue = `
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
	cueInstance, err := cueRuntime.Compile("schemav1", SchemaV1MetadataCue)
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
type ModelV1Metadata struct {
	Plural              string
	SupportedExtensions []string
}

const ModelV1MetadataCue = `
{
	_model: {
		plural: string
		supportedExtensions: [...string]
	}
}
`

// GetV1Model returns the V1_MODEL information from a cue.Value
func GetV1Model(schema cue.Value) (ModelV1Metadata, error) {
	var cueRuntime cue.Runtime
	modelV1Metadata := ModelV1Metadata{}

	cueInstance, err := cueRuntime.Compile("", ModelV1MetadataCue)
	if err != nil {
		return modelV1Metadata, err
	}

	fields, err := schema.Fields(cue.Hidden(true))
	for fields.Next() {
		if fields.Label() == "_model" {
			modelSchema := fields.Value()
			modelSchema.Unify(cueInstance.Value())

			if err = modelSchema.Validate(); err != nil {
				return modelV1Metadata, err
			}

			err = modelSchema.Decode(&modelV1Metadata)
			if err != nil {
				return modelV1Metadata, nil
			}

			return modelV1Metadata, nil
		}
	}

	return modelV1Metadata, fmt.Errorf("Couldn't find _schema")
}
