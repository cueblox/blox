package cuedb

import (
	"fmt"
	"path"
	"strings"

	"cuelang.org/go/cue"
	"github.com/cueblox/blox/config"
	"github.com/hashicorp/go-multierror"
)

// Database is the "world" struct. We can "insert" records
// into it and know immediately if they're valid or not.
type Database struct {
	runtime *cue.Runtime
	config  *config.BloxConfig
	db      cue.Value
	tables  map[string]Table
}

// NewDatabase creates a "world" struct to store
// records
func NewDatabase(cfg *config.BloxConfig) (Database, error) {
	var cueRuntime cue.Runtime
	cueInstance, err := cueRuntime.Compile("", "")

	if nil != err {
		return Database{}, err
	}

	return Database{
		runtime: &cueRuntime,
		config:  cfg,
		db:      cueInstance.Value(),
		tables:  make(map[string]Table),
	}, nil
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
		plural, err := GetV1Model(fields.Value())
		if nil != err {
			continue
		}

		table := Table{
			schemaNamespace: metadata.Namespace,
			schemaName:      metadata.Name,
			name:            fields.Label(),
			plural:          plural,
		}

		if _, ok := d.tables[table.ID()]; ok {
			return fmt.Errorf("Table with name '%s' already registered", fields.Label())
		}

		err = d.db.Validate()
		if nil != err {
			return err
		}

		inst, err := d.runtime.Compile("", fmt.Sprintf("{%s: _\n%s: [ID=string]: %s}", fields.Label(), plural, fields.Label()))
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

// MarshalJSON returns the database encoded in JSON format
func (d *Database) MarshalJSON() ([]byte, error) {
	return d.db.MarshalJSON()
}

// Table represents a schema record
type Table struct {
	name   string
	plural string // The directory where we find the records

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
	return t.plural
}

// CuePath returns the plural form
// of the table's name
func (t *Table) CuePath() cue.Path {
	return cue.ParsePath(t.plural)
}

func (d *Database) SourcePath(t Table) string {
	return path.Join(d.config.SourceDir, t.Directory())
}

func (d *Database) DestinationPath(t Table) string {
	return path.Join(d.config.BuildDir, t.Directory())
}

func (d *Database) StaticPath(t Table) string {
	return path.Join(d.config.StaticDir, t.Directory())
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
		// fmt.Println("Finding Def: ", table.name)
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

				inst, err := d.runtime.Compile("", fmt.Sprintf("{%s: _\n%s: %s: or([ for k, _ in %s {k}])}", foreignTable.plural, table.name, fields.Label(), foreignTable.plural))
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
