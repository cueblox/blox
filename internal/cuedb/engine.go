package cuedb

import (
	"fmt"
	"sort"
	"strings"

	"cuelang.org/go/cue"
	"github.com/cueblox/blox"
	"github.com/hashicorp/go-multierror"
	"github.com/heimdalr/dag"
	"github.com/pterm/pterm"
)

const (
	dataPathRoot     = "data"
	dataSetField     = "_dataset"
	schemataPathRoot = "schemata"
	schemaField      = "_schema"
)

type Engine struct {
	// embedded runtime database
	*blox.Runtime
	dataSets map[string]DataSet
}

// RecordBaseCue is the "Base" configuration that blox
// expects to exist, but doesn't enforce in user-land.
// We'll inject this Cue into each DataSet definition.
const RecordBaseCue = `{
	id: string
}`

func NewEngine() (*Engine, error) {
	r, err := blox.NewRuntime()
	if err != nil {
		return nil, err
	}

	runtime := &Engine{
		Runtime:  r,
		dataSets: make(map[string]DataSet),
	}

	return runtime, nil
}

func (r *Engine) CountDataSets() int {
	return len(r.dataSets)
}

func (r *Engine) extractSchemaMetadata(schema cue.Value) (SchemaMetadata, error) {
	cueInstance, err := r.CueRuntime.Compile("", SchemaMetadataCue)
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

func (r *Engine) extractDataSetMetadata(schema cue.Value) (DataSetMetadata, error) {
	cueInstance, err := r.CueRuntime.Compile("", DataSetMetadataCue)
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

func (r *Engine) GetDataSets() map[string]DataSet {
	return r.dataSets
}

type DagNode struct {
	Name string
}

func (d *DagNode) ID() string {
	return d.Name
}

func (r *Engine) GetDataSetsDAG() *dag.DAG {
	graph := dag.NewDAG()

	graph.AddVertex(&DagNode{Name: "root"})
	datasets := r.GetDataSets()

	keys := make([]string, 0, len(datasets))
	for k := range datasets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		dataSet := datasets[k]

		// AddVertex returns a string ID. We don't need to worry
		// about it, because the method checks for an ID() method
		// on the struct, which we have.
		d := DagNode{Name: dataSet.ID()}
		_, err := graph.AddVertex(&d)
		if err != nil {
			pterm.Warning.Printf("failed to add vertex: %v", err)
			continue
		}

		graph.AddEdge("root", d.ID())
	}

	for _, k := range keys {
		dataSet := datasets[k]
		for _, relationship := range dataSet.relationships {
			edge, _ := r.GetDataSet(relationship)
			graph.AddEdge(dataSet.ID(), edge.ID())
		}
	}

	return graph
}

func (r *Engine) GetDataSet(name string) (DataSet, error) {
	cueName := strings.ToLower(name)
	if !strings.HasPrefix(cueName, "#") {
		cueName = fmt.Sprintf("#%s", strings.ToLower(name))
	}

	if dataSet, ok := r.dataSets[cueName]; ok {
		return dataSet, nil
	}

	return DataSet{}, fmt.Errorf("couldn't find DataSet with name %s", name)
}

func (d *DataSet) ID() string {
	return strings.ToLower(d.name)
}

func (d *DataSet) String() string {
	return d.ID()
}

func (d *DataSet) GetDataDirectory() string {
	return d.metadata.Plural
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

func (r *Engine) RegisterSchema(cueString string) error {
	cueInstance, err := r.CueRuntime.Compile("", cueString)
	if nil != err {
		return err
	}

	cueValue := cueInstance.Value()

	schemaMetadata, err := r.extractSchemaMetadata(cueValue)
	if err != nil {
		return err
	}
	schemaPath := cue.ParsePath(fmt.Sprintf(`%s."%s"."%s"`, schemataPathRoot, schemaMetadata.Namespace, schemaMetadata.Name))

	// First, Unify whatever schemas the users want. We'll
	// do our best to extract whatever information from
	// it we require
	r.Database = r.Database.FillPath(schemaPath, cueValue)

	// Find DataSets and register
	fields, err := cueValue.Fields(cue.Definitions(true))
	if err != nil {
		return err
	}

	// Base Record Constraints
	baseRecordInstance, err := r.CueRuntime.Compile("", RecordBaseCue)
	if err != nil {
		return err
	}
	baseRecordValue := baseRecordInstance.Value()

	for fields.Next() {
		pterm.Debug.Println("\t\t\tNext field")
		if !fields.IsDefinition() {
			// Only Definitions can be registered as DataSets
			continue
		}

		dataSetMetadata, err := r.extractDataSetMetadata(fields.Value())
		if err != nil {

			pterm.Debug.Printf("\t\t\t%v\n", err)
			// No dataset metadata, skip
			continue
		}
		pterm.Debug.Printf("\t\t\t%s\n", fields.Value())

		// Find relationships
		relationships, err := getDataSetRelationships(fields.Label(), fields.Value())
		if err != nil {
			return err
		}

		dataSet := DataSet{
			schemaMetadata: schemaMetadata,
			schema:         fields.Value(),
			name:           fields.Label(),
			relationships:  relationships,
			metadata:       dataSetMetadata,
			cuePath:        schemaPath,
		}

		if _, ok := r.dataSets[dataSet.name]; ok {
			return fmt.Errorf("DataSet with name '%s' already registered", fields.Label())
		}

		// Compile our BaseModel
		r.Database = r.Database.FillPath(dataSet.GetDefinitionPath(), baseRecordValue)
		if err = r.Database.Validate(); err != nil {
			return err
		}

		r.dataSets[strings.ToLower(dataSet.name)] = dataSet

		inst, err := r.CueRuntime.Compile("", dataSet.GetDataMapCue())
		if err != nil {
			return err
		}

		r.Database = r.Database.FillPath(cue.Path{}, inst.Value())

		if err := r.Database.Validate(); nil != err {
			return err
		}
	}

	return nil
}

func (r *Engine) Insert(dataSet DataSet, record map[string]interface{}) error {
	r.Database = r.Database.FillPath(dataSet.CueDataPath(), record)

	err := r.Database.Validate()
	if nil != err {
		return err
	}

	return nil
}

func (r *Engine) GetAllData(dataSetName string) cue.Value {
	d, err := r.GetDataSet(dataSetName)

	if err != nil {
		return cue.Value{}
	}

	return r.Database.LookupPath(d.CueDataPath())
}

// MarshalJSON returns the database encoded in JSON format
func (r *Engine) MarshalJSON() ([]byte, error) {
	return r.Database.LookupPath(cue.ParsePath(dataPathRoot)).MarshalJSON()
}

func getDataSetRelationships(label string, schema cue.Value) ([]string, error) {
	fields, err := schema.Fields(cue.All())
	if err != nil {
		return []string{}, err
	}

	relationships := []string{}

	for fields.Next() {
		relationship := fields.Value().Attribute("relationship")
		if err = relationship.Err(); err == nil {
			relationships = append(relationships, relationship.Contents())
			continue
		}

		if strings.HasSuffix(fields.Label(), "_id") {
			relationships = append(relationships, strings.TrimSuffix(fields.Label(), "_id"))
		}
	}

	return relationships, nil
}

// ReferentialIntegrity checks the relationships between
// the records in the content database
func (r *Engine) ReferentialIntegrity() error {
	for _, dataSet := range r.GetDataSets() {
		// Walk each field and look for _id labels
		val := r.Database.LookupPath(dataSet.GetDefinitionPath())

		fields, err := val.Fields(cue.All())
		if err != nil {
			return err
		}

		for fields.Next() {
			relationship := fields.Value().Attribute("relationship")

			// If err is nil, that means we successfully found a relationship.
			// Lets build the integrity schema then continue
			if err = relationship.Err(); err == nil {
				foreignTable, err := r.GetDataSet(strings.ToLower(relationship.Contents()))
				if err != nil {
					return err
				}

				optional := ""
				if fields.IsOptional() {
					optional = "?"
				}

				inst, err := r.CueRuntime.Compile("", fmt.Sprintf("{data: _\n%s: %s: %s%s: or([ for k, _ in data.%s {k}])}", dataSet.GetInlinePath(), dataSet.name, fields.Label(), optional, foreignTable.GetDataDirectory()))
				if err != nil {
					return err
				}

				r.Database = r.Database.FillPath(cue.Path{}, inst.Value())
				continue
			}

			if strings.HasSuffix(fields.Label(), "_id") {
				foreignTable, err := r.GetDataSet(strings.TrimSuffix(fields.Label(), "_id"))
				if err != nil {
					return err
				}

				optional := ""
				if fields.IsOptional() {
					optional = "?"
				}

				inst, err := r.CueRuntime.Compile("", fmt.Sprintf("{data: _\n%s: %s: %s%s: or([ for k, _ in data.%s {k}])}", dataSet.GetInlinePath(), dataSet.name, fields.Label(), optional, foreignTable.GetDataDirectory()))
				if err != nil {
					return err
				}

				r.Database = r.Database.FillPath(cue.Path{}, inst.Value())
			}
		}
	}

	err := r.Database.Validate()
	if err != nil {
		return multierror.Prefix(err, "Referential Integrity Failed")
	}

	return nil
}

func (r *Engine) GetOutput() (cue.Value, error) {
	if len(r.GetDataSets()) == 0 {
		return cue.Value{}, fmt.Errorf("No DataSets to generate output")
	}

	for _, dataSet := range r.GetDataSets() {
		inst, err := r.CueRuntime.Compile("", fmt.Sprintf("{%s: %s: _\noutput: %s: [ for key, val in %s.%s {val}]}", dataPathRoot, dataSet.metadata.Plural, dataSet.metadata.Plural, dataPathRoot, dataSet.metadata.Plural))
		if err != nil {
			return cue.Value{}, err
		}

		r.Database = r.Database.FillPath(cue.Path{}, inst.Value())
	}

	return r.Database.LookupPath(cue.ParsePath("output")), nil
}
