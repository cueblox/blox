package cuedb

import (
	"fmt"
	"strings"

	"cuelang.org/go/cue"
)

// Can't use schemaField, yet.
const SchemaMetadataCue = `{
	_schema: {
		namespace: string
		name: string
	}
}`

type SchemaMetadata struct {
	Namespace string
	Name      string
}

// Can't use dataSetField, yet.
const DataSetMetadataCue = `{
	_dataset: {
		plural: string
		supportedExtensions: [...string]
	}
}`

type DataSetMetadata struct {
	Plural              string
	SupportedExtensions []string
}

type DataSet struct {
	name           string
	schemaMetadata SchemaMetadata
	schema         cue.Value
	cuePath        cue.Path
	metadata       DataSetMetadata
	relationships  []string
}

func (d *DataSet) GetDataMapCue() string {
	return fmt.Sprintf(`{
		%s: %s: _
		%s: %s: [ID=string]: %s.%s & {id: (ID)}
	}`,
		d.GetInlinePath(), d.name,
		dataPathRoot, d.metadata.Plural, d.cuePath.String(), d.name,
	)
}

func (d *DataSet) GetPluralName() string {
	return strings.Title(d.metadata.Plural)
}

func (d *DataSet) GetExternalName() string {
	return strings.Replace(d.name, "#", "", 1)
}

func (d *DataSet) GetSchemaCue() cue.Value {
	return d.schema
}

func (d *DataSet) CueDataPath() cue.Path {
	return cue.ParsePath(fmt.Sprintf("%s.%s", dataPathRoot, d.metadata.Plural))
}

func (d *DataSet) IsSupportedExtension(ext string) bool {
	for _, val := range d.metadata.SupportedExtensions {
		if val == ext {
			return true
		}
	}

	return false
}

func (d *DataSet) GetSupportedExtensions() []string {
	return d.metadata.SupportedExtensions
}
