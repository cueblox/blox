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

type Model struct {
	plural              string
	supportedExtensions []string
}

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

	err = database.LoadConfig()
	if nil != err {
		return Database{}, err
	}

	return database, nil
}

func (d *Database) LoadConfig() error {
	configInstance, err := d.runtime.Compile("", BaseConfig)
	if err != nil {
		return err
	}

	configValue := configInstance.Value()

	localConfig, err := ioutil.ReadFile("blox.cue")
	if err != nil {
		return err
	}

	localConfigInstance, err := d.runtime.Compile("", localConfig)
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

// RegisterTables ensures that the cueString schema is a valid schema
// and parses the Cue to find Models within. Each Model is registered
// as a Table, provided the name is available.
func (d *Database) RegisterTables(cueString string) error {
	cueInstance, err := d.runtime.Compile(cueString, cueString)
	if nil != err {
		return err
	}

	cueValue := cueInstance.Value()

	// First, Unify whatever schemas the users want. We'll
	// do our best to extract whatever information from
	// it we require
	d.db = d.db.FillPath(cue.Path{}, cueValue)

	// Is the Schema valid?
	_, err = GetSchemaVersion(cueValue)
	if err != nil {
		return err
	}

	// We only have a V1 :)
	metadata, err := GetSchemaV1Metadata(cueValue)
	if err != nil {
		return err
	}

	// Find Models and register as a table
	fields, err := cueValue.Fields(cue.Definitions(true))
	if err != nil {
		return err
	}

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
		}

		if _, ok := d.tables[table.ID()]; ok {
			return fmt.Errorf("Table with name '%s' already registered", fields.Label())
		}

		err = d.db.Validate()
		if nil != err {
			return err
		}

		inst, err := d.runtime.Compile("", fmt.Sprintf("{%s: _\n%s: [ID=string]: %s}", fields.Label(), modelV1Metadata.Plural, fields.Label()))
		if err != nil {
			return err
		}

		d.db = d.db.FillPath(cue.Path{}, inst.Value())

		if err := d.db.Validate(); nil != err {
			return err
		}

		d.tables[table.ID()] = table
	}

	return nil
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

// MarshalJSON returns the database encoded in JSON format
func (d *Database) MarshalJSON() ([]byte, error) {
	return d.db.MarshalJSON()
}

func (d *Database) DumpAll() {
	fmt.Println(d.config)
	fmt.Println(d.db)
}

// Table represents a schema record
type Table struct {
	name     string
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

// CuePath returns the plural form
// of the table's name
func (t *Table) CuePath() cue.Path {
	return cue.ParsePath(t.metadata.Plural)
}

// Insert adds a record
func (d *Database) Insert(table Table, record map[string]interface{}) error {
	filledValued := d.db.FillPath(table.CuePath(), record)

	err := filledValued.Validate()
	if nil != err {
		return err
	}

	d.db = d.db.Unify(filledValued)

	return nil
}

// ReferentialIntegrity checks the relationships between
// the records in the content database
func (d *Database) ReferentialIntegrity() error {
	for _, table := range d.GetTables() {
		// Walk each field and look for _id labels
		val := d.db.LookupDef(table.name)

		fields, err := val.Fields(cue.Optional(true))
		if err != nil {
			return err
		}

		for fields.Next() {
			if strings.HasSuffix(fields.Label(), "_id") {
				foreignTable, err := d.GetTable(fmt.Sprintf("#%s", strings.TrimSuffix(fields.Label(), "_id")))
				if err != nil {
					return err
				}

				inst, err := d.runtime.Compile("", fmt.Sprintf("{%s: _\n%s: %s: or([ for k, _ in %s {k}])}", foreignTable.metadata.Plural, table.name, fields.Label(), foreignTable.metadata.Plural))
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
