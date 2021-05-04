package faunadb

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/cueblox/blox/internal/sync"
	f "github.com/fauna/faunadb-go/v4/faunadb"
	"github.com/pterm/pterm"
)

func init() {
	sync.Register("faunadb", &FaunaProvider{})
}

// FaunaProvider satisfies the sync.SyncProvider interface
// returning a FaunaClient which implements sync.EngineProvider
type FaunaProvider struct{}

func (f *FaunaProvider) Initialize() (sync.EngineProvider, error) {
	key := os.Getenv("FAUNA_KEY")
	c := &FaunaClient{}
	if key == "" {
		return nil, errors.New("no FAUNA_KEY found.\n" + c.Help())
	}
	c.key = key
	return c, nil
}

func (f *FaunaProvider) Name() string {
	return "FaunaDB"
}

// FaunaClient satisfies the sync.EngineProvider interface
type FaunaClient struct {
	key    string
	client *f.FaunaClient
}

func (c *FaunaClient) ensureConnect() {
	if c.client == nil {
		c.client = f.NewFaunaClient(c.key)
	}
}
func (c *FaunaClient) Sync(bb []byte) error {
	data, err := sync.MakeMap(bb)
	if err != nil {
		return err
	}

	// make sure we have a client
	c.ensureConnect()
	tables := sync.GetTables(data)
	// make sure tables exist
	err = c.ensureTables(tables)
	if err != nil {
		return err
	}

	for _, table := range tables {
		tableData, ok := data[table].([]interface{})
		if !ok {
			return errors.New("unable to read table data")
		}
		err := c.syncTable(table, tableData)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *FaunaClient) Help() string {
	return `Create a database in FaunaDB, then create a database specific
access key.  Export this key as FAUNA_KEY in your environment.`
}

func (c *FaunaClient) checkOrCreateCollection(name string) error {
	_, err := c.client.Query(f.CreateCollection(f.Obj{"name": name}))
	return err

}

func (c *FaunaClient) ensureTables(tables []string) error {
	for _, table := range tables {
		err := c.checkOrCreateCollection(table)
		if err != nil {
			if !strings.Contains(err.Error(), "Collection already exists") {
				pterm.Error.Println(err)
				return err
			}
		}
	}
	return nil

}

func (c *FaunaClient) syncTable(table string, data []interface{}) error {
	for index, record := range data {
		fmt.Printf("%s: record %d\n", table, index)
		_, err := c.client.Query(
			f.Create(
				f.Collection(table),
				f.Obj{"data": record},
			),
		)
		if err != nil {
			return err
		}
	}
	return nil
}
