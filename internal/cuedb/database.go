package cuedb

import (
	_ "embed"
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"cuelang.org/go/cue"
	"github.com/hashicorp/go-multierror"
)

//go:embed config.cue
var BaseConfig string

// Database is the "world" struct. We can "insert" records
// into it and know immediately if they're valid or not.
type Database struct {
	runtime *cue.Runtime
	config  *cue.Value
	db      cue.Value
	tables  map[string]Table
}

// NewDatabase creates a "world" struct to store
// records
func NewDatabase() (Database, error) {
	var cueRuntime cue.Runtime
	cueInstance, err := cueRuntime.Compile("", "")

	if nil != err {
		return Database{}, err
	}

	database := Database{
		runtime: &cueRuntime,
		db:      cueInstance.Value(),
		tables:  make(map[string]Table),
	}

	configStr, err := readConfigFile()
	if err != nil {
		return database, err
	}

	err = database.loadConfig(configStr)
	if nil != err {
		return Database{}, err
	}

	return database, nil
}

func newDatabaseWithConfig(config string) (Database, error) {
	var cueRuntime cue.Runtime
	cueInstance, err := cueRuntime.Compile("", "")

	if nil != err {
		return Database{}, err
	}

	database := Database{
		runtime: &cueRuntime,
		db:      cueInstance.Value(),
		tables:  make(map[string]Table),
	}

	err = database.loadConfig(config)
	if nil != err {
		return Database{}, err
	}

	return database, nil
}

func readConfigFile() (string, error) {
	config, err := ioutil.ReadFile("blox.cue")
	if err != nil {
		return "", err
	}

	return string(config), nil
}

func (d *Database) loadConfig(config string) error {
	configInstance, err := d.runtime.Compile("", BaseConfig)
	if err != nil {
		return err
	}

	configValue := configInstance.Value()

	localConfigInstance, err := d.runtime.Compile("", config)
	if err != nil {
		return err
	}

	mergedConfig := configValue.Unify(localConfigInstance.Value())
	if err = mergedConfig.Validate(cue.Concrete(true)); err != nil {
		return err
	}

	d.config = &mergedConfig

	return nil
}

func (d *Database) GetConfigString(key string) (string, error) {
	value, err := d.config.LookupField(key)
	if err != nil {
		return "", err
	}

	str, err := value.Value.String()
	if err != nil {
		return "", err
	}

	return str, nil
}

const BaseModelCue = `{
	id: string
}`

// RegisterTables ensures that the cueString schema is a valid schema
// and parses the Cue to find Models within. Each Model is registered
// as a Table, provided the name is available.
func (d *Database) RegisterTables(cueString string) error {
	cueInstance, err := d.runtime.Compile("", cueString)
	if nil != err {
		return err
	}

	cueValue := cueInstance.Value()

	// We only have a V1 :)
	metadata, err := GetSchemaV1Metadata(cueValue)
	if err != nil {
		return err
	}

	// First, Unify whatever schemas the users want. We'll
	// do our best to extract whatever information from
	// it we require
	schemaPath := cue.ParsePath(fmt.Sprintf(`schema."%s"."%s"`, metadata.Namespace, metadata.Name))

	d.db = d.db.FillPath(schemaPath, cueValue)

	// Find Models and register as a table
	fields, err := cueValue.Fields(cue.Definitions(true))
	if err != nil {
		return err
	}

	baseModelInstance, err := d.runtime.Compile("", BaseModelCue)
	if nil != err {
		return err
	}
	baseModelValue := baseModelInstance.Value()

	for fields.Next() {
		if !fields.IsDefinition() {
			// Only Definitions can be registered as tables
			continue
		}

		// We have a Definition, does it define a model?
		modelV1Metadata, err := GetV1Model(fields.Value())
		if nil != err {
			continue
		}

		table := Table{
			schemaNamespace: metadata.Namespace,
			schemaName:      metadata.Name,
			name:            fields.Label(),
			metadata:        modelV1Metadata,
			cuePath:         schemaPath,
		}

		if _, ok := d.tables[table.ID()]; ok {
			return fmt.Errorf("Table with name '%s' already registered", fields.Label())
		}

		// Compile our BaseModel
		mergedModel := fields.Value().Unify(baseModelValue)
		d.db = d.db.FillPath(table.GetDefPath(), mergedModel)

		if err = d.db.Validate(); err != nil {
			return err
		}

		d.tables[table.ID()] = table

		inst, err := d.runtime.Compile("", table.DataKey())
		if err != nil {
			return err
		}

		d.db = d.db.FillPath(cue.Path{}, inst.Value())

		if err := d.db.Validate(); nil != err {
			return err
		}

	}

	return nil
}

func (t *Table) DataKey() string {
	return fmt.Sprintf(`{
		%s: %s: _
		data: %s: [ID=string]: %s.%s

}`,
		t.InlinePath(), t.name,
		t.metadata.Plural, t.cuePath.String(), t.name,
	)
}

// GetTables returns the tables in the Database
func (d *Database) GetTables() map[string]Table {
	return d.tables
}

// GetTable returns a single table in the database
func (d *Database) GetTable(name string) (Table, error) {
	if table, ok := d.tables[name]; ok {
		return table, nil
	}

	return Table{}, fmt.Errorf("Table '%s' doesn't exist in database", name)
}

func (d *Database) GetTableDataDir(table Table) string {
	dataDir, err := d.GetConfigString("data_dir")
	if err != nil {
		// Config is already validated at this stage, should
		// never happen
		panic("Unexpected error fetching data_dir")
	}

	return path.Join(dataDir, table.Directory())
}

func (d *Database) DumpAll() {
	fmt.Println(d.config)
	fmt.Println(d.db)
}

// Table represents a schema record
type Table struct {
	name     string
	cuePath  cue.Path
	metadata ModelV1Metadata

	// Which schema registered this table?
	schemaNamespace string
	schemaName      string
}

// ID returns the table's name
func (t *Table) ID() string {
	return strings.ToLower(t.name)
}

// Directory returns the plural form
// of the table name
func (t *Table) Directory() string {
	return t.metadata.Plural
}

func (t *Table) IsSupportedExtension(ext string) bool {
	for _, val := range t.metadata.SupportedExtensions {
		if val == ext {
			return true
		}
	}

	return false
}

func (t *Table) GetSupportedExtensions() []string {
	return t.metadata.SupportedExtensions
}

func (t *Table) GetDefPath() cue.Path {
	return cue.ParsePath(t.cuePath.String() + "." + t.name)
}

// CuePath returns the plural form
// of the table's name
func (t *Table) CuePath() cue.Path {
	return t.cuePath
}

func (t *Table) InlinePath() string {
	inlinePath := ""

	for _, seg := range t.cuePath.Selectors() {
		inlinePath = fmt.Sprintf("%s: %s", inlinePath, seg)
	}

	return strings.TrimPrefix(inlinePath, ": ")
}

func (t *Table) CueDataPath() cue.Path {
	return cue.ParsePath(fmt.Sprintf("data.%s", t.metadata.Plural))
}

// Insert adds a record
func (d *Database) Insert(table Table, record map[string]interface{}) error {
	d.db = d.db.FillPath(table.CueDataPath(), record)
	err := d.db.Validate()
	if nil != err {
		return err
	}

	return nil
}

// MarshalJSON returns the database encoded in JSON format
func (d *Database) MarshalJSON() ([]byte, error) {
	data := d.db.LookupPath(cue.ParsePath("data"))
	return data.MarshalJSON()
}

// ReferentialIntegrity checks the relationships between
// the records in the content database
func (d *Database) ReferentialIntegrity() error {
	for _, table := range d.GetTables() {
		// Walk each field and look for _id labels
		val := d.db.LookupPath(table.GetDefPath())

		fields, err := val.Fields(cue.All())
		if err != nil {
			return err
		}

		for fields.Next() {
			if strings.HasSuffix(fields.Label(), "_id") {
				foreignTable, err := d.GetTable(fmt.Sprintf("#%s", strings.TrimSuffix(fields.Label(), "_id")))
				if err != nil {
					return err
				}

				inst, err := d.runtime.Compile("", fmt.Sprintf("{data: %s: _\n%s: %s: %s: or([ for k, _ in data.%s {k}])}", foreignTable.metadata.Plural, table.InlinePath(), table.name, fields.Label(), foreignTable.metadata.Plural))
				if err != nil {
					return err
				}

				d.db = d.db.FillPath(cue.Path{}, inst.Value())
			}
		}
	}

	err := d.db.Validate()
	if err != nil {
		return multierror.Prefix(err, "Referential Integrity Failed")
	}

	return nil
}

func (d *Database) GetOutput() cue.Value {
	for _, table := range d.GetTables() {
		inst, err := d.runtime.Compile("", fmt.Sprintf("{data: %s: _\noutput: %s: [ for key, val in data.%s {val & {id: key} }]}", table.metadata.Plural, table.metadata.Plural, table.metadata.Plural))
		if err != nil {
			return d.db
		}
		d.db = d.db.FillPath(cue.Path{}, inst.Value())
	}

	return d.db.LookupPath(cue.ParsePath("output"))
}
