package cuedb

import (
	"testing"
)

func TestRegisterTable(t *testing.T) {
	db, err := NewDatabase()
	if nil != err {
		t.FailNow()
	}

	tableOne := `{
		_schema: {
			version: "v1"
			name: "test"
			namespace: "testns"
		}

		#TableOne: {
			_model: {
				plural: "tables"
			}

			name: string
		}

		#TableTwo: {
			name: string
		}

		#TableThree: {
			_model: {
				plural: "threes"
			}

			name: string
		}
	}`

	err = db.RegisterTables(tableOne)
	if nil != err {
		t.Fatalf("Failed to register table '%s': %v", tableOne, err)
	}

	table, err := db.GetTable("#tableone")
	if nil != err {
		t.Fatalf("Failed to GetTable: %s", err)
	}

	if "#TableOne" != table.name {
		t.Fatalf("TableOne wasn't registered with correct id. Expected '#TableOne' got '%s'", table.name)
	}

	if "test" != table.schemaName {
		t.Fatal("Schema Name not registered correctly")
	}

	if "testns" != table.schemaNamespace {
		t.Fatal("Schema Namespace not registered correctly")
	}

	_, err = db.GetTable("#tabletwo")
	if nil == err {
		t.Fatal("Definition without model was registered: TableTwo")
	}

	_, err = db.GetTable("#tablethree")
	if nil != err {
		t.Fatalf("Second model in Cue wasn't registered: %s", err)
	}

	// Make sure a table can't be registered twice
	err = db.RegisterTables(tableOne)
	if nil == err {
		t.Fatalf("Was able to register table when ID was already registered: %s", tableOne)
	}
}
