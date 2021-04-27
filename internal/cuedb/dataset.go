package cuedb

import (
	"fmt"

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
	cuePath        cue.Path
	metadata       DataSetMetadata
}

// AddDataMap
func (d *DataSet) GetDataMapCue() string {
	return fmt.Sprintf(`{
		%s: %s: _
		%s: %s: [ID=string]: %s.%s & {id: (ID)}
	}`,
		d.GetInlinePath(), d.name,
		dataPathRoot, d.metadata.Plural, d.cuePath.String(), d.name,
	)
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
