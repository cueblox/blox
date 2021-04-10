package cuedb

import (
	"fmt"
	"strings"

	"cuelang.org/go/cue"
	"github.com/pterm/pterm"
)

const (
	dataPathRoot   = "data"
	dataSetField   = "_dataset"
	schemaPathRoot = "schema"
	schemaField    = "_schema"
)

type Runtime struct {
	cueRuntime *cue.Runtime
	database   cue.Value
	dataSets   map[string]DataSet
}

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

// NewRuntime setups a new database for DataSets to be registered,
// and data inserted.
func NewRuntime() (Runtime, error) {
	var cueRuntime cue.Runtime
	cueInstance, err := cueRuntime.Compile("", "")

	if nil != err {
		return Runtime{}, err
	}

	runtime := Runtime{
		cueRuntime: &cueRuntime,
		database:   cueInstance.Value(),
		dataSets:   make(map[string]DataSet),
	}

	return runtime, nil
}

func (r *Runtime) CountDataSets() int {
	return len(r.dataSets)
}

func (r *Runtime) extractSchemaMetadata(schema cue.Value) (SchemaMetadata, error) {
	cueInstance, err := r.cueRuntime.Compile("", SchemaMetadataCue)
	if err != nil {
		pterm.Debug.Println("Failed to compile Cue")
		return SchemaMetadata{}, err
	}

	// Ensure we have all the required fields
	versionedSchema := schema.Unify(cueInstance.Value())
	if err = versionedSchema.Validate(); err != nil {
		return SchemaMetadata{}, err
	}

	schemaMetadata := SchemaMetadata{}
	schemaValue := versionedSchema.LookupPath(cue.MakePath(cue.Hid(schemaField, "_")))

	err = schemaValue.Decode(&schemaMetadata)
	if err != nil {
		return SchemaMetadata{}, err
	}

	return schemaMetadata, nil
}

func (r *Runtime) extractDataSetMetadata(schema cue.Value) (DataSetMetadata, error) {
	cueInstance, err := r.cueRuntime.Compile("", DataSetMetadataCue)
	if err != nil {
		pterm.Debug.Println("Failed to compile Cue")
		return DataSetMetadata{}, err
	}

	// Ensure we have all the required fields
	dataSetMetadataCueVal := schema.Unify(cueInstance.Value())
	if err = dataSetMetadataCueVal.Validate(); err != nil {
		return DataSetMetadata{}, err
	}

	dataSetMetadata := DataSetMetadata{}
	schemaValue := dataSetMetadataCueVal.LookupPath(cue.MakePath(cue.Hid(dataSetField, "_")))

	err = schemaValue.Decode(&dataSetMetadata)
	if err != nil {
		return DataSetMetadata{}, err
	}

	return dataSetMetadata, nil
}

func (r *Runtime) GetDataSets() map[string]DataSet {
	return r.dataSets
}

func (r *Runtime) GetDataSet(name string) (DataSet, error) {
	cueName := name
	if !strings.HasPrefix(cueName, "#") {
		cueName = fmt.Sprintf("#%s", name)
	}

	if dataSet, ok := r.dataSets[cueName]; ok {
		return dataSet, nil
	}

	return DataSet{}, fmt.Errorf("Couldn't find DataSet with name %s", name)
}

func (d *DataSet) ID() string {
	return strings.ToLower(d.name)
}

func (d *DataSet) GetDefinitionPath() cue.Path {
	return cue.ParsePath(d.cuePath.String() + "." + d.name)
}

// GetInlinePath returns an inline cue.Path that can be used within a Cue document
// like "some: key: #Here"
func (d *DataSet) GetInlinePath() string {
	inlinePath := ""

	for _, seg := range d.cuePath.Selectors() {
		inlinePath = fmt.Sprintf("%s: %s", inlinePath, seg)
	}

	return strings.TrimPrefix(inlinePath, ": ")
}

// AddDataMap
func (d *DataSet) GetDataMapCue() string {
	return fmt.Sprintf(`{
		%s: %s: _
		%s: %s: [ID=string]: %s.%s
	}`,
		d.GetInlinePath(), d.name,
		dataPathRoot, d.metadata.Plural, d.cuePath.String(), d.name,
	)
}

func (r *Runtime) RegisterSchema(cueString string) error {
	cueInstance, err := r.cueRuntime.Compile("", cueString)
	if nil != err {
		return err
	}

	cueValue := cueInstance.Value()

	schemaMetadata, err := r.extractSchemaMetadata(cueValue)
	if err != nil {
		return err
	}

	schemaPath := cue.ParsePath(fmt.Sprintf(`%s."%s"."%s"`, schemaPathRoot, schemaMetadata.Namespace, schemaMetadata.Name))

	// First, Unify whatever schemas the users want. We'll
	// do our best to extract whatever information from
	// it we require
	r.database = r.database.FillPath(schemaPath, cueValue)

	// Find DataSets and register
	fields, err := cueValue.Fields(cue.Definitions(true))
	if err != nil {
		return err
	}

	for fields.Next() {
		if !fields.IsDefinition() {
			// Only Definitions can be registered as DataSets
			continue
		}

		dataSetMetadata, err := r.extractDataSetMetadata(fields.Value())
		if err != nil {
			return err
		}

		dataset := DataSet{
			schemaMetadata: schemaMetadata,
			name:           fields.Label(),
			metadata:       dataSetMetadata,
			cuePath:        schemaPath,
		}

		if _, ok := r.dataSets[dataset.name]; ok {
			return fmt.Errorf("DataSet with name '%s' already registered", fields.Label())
		}

		// Compile our BaseModel
		r.database = r.database.FillPath(dataset.GetDefinitionPath(), fields.Value())
		if err = r.database.Validate(); err != nil {
			return err
		}

		r.dataSets[dataset.name] = dataset

		// TODO: Insert DataKey
		inst, err := r.cueRuntime.Compile("", dataset.GetDataMapCue())
		if err != nil {
			return err
		}

		r.database = r.database.FillPath(cue.Path{}, inst.Value())

		if err := r.database.Validate(); nil != err {
			return err
		}
	}

	return nil
}

func (d *DataSet) CueDataPath() cue.Path {
	return cue.ParsePath(fmt.Sprintf("%s.%s", dataPathRoot, d.metadata.Plural))
}

func (r *Runtime) Insert(dataSet DataSet, record map[string]interface{}) error {
	r.database = r.database.FillPath(dataSet.CueDataPath(), record)

	err := r.database.Validate()
	if nil != err {
		return err
	}

	return nil
}
