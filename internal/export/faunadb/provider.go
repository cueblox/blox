package faunadb

import (
	"errors"
	"fmt"
	"hash/adler32"
	"os"
	"strings"

	"github.com/cueblox/blox/internal/export"
	f "github.com/fauna/faunadb-go/v4/faunadb"
	"github.com/pterm/pterm"
)

func init() {
	export.Register("faunadb", &FaunaProvider{})
}

// FaunaProvider satisfies the export.ExportProvider interface
// returning a FaunaClient which implements export.EngineProvider
type FaunaProvider struct{}

func (f *FaunaProvider) Initialize() (export.EngineProvider, error) {
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

// FaunaClient satisfies the export.EngineProvider interface
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
	data, err := export.MakeMap(bb)
	if err != nil {
		return err
	}

	// make sure we have a client
	c.ensureConnect()
	tables := export.GetTables(data)
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
	if !strings.Contains(err.Error(), "Collection already exists") {
		pterm.Error.Println(err)
		return err
	}
	_, err = c.client.Query(
		f.CreateIndex(
			f.Obj{
				"name":   name + "-index",
				"source": f.Collection(name),
			}))

	if !strings.Contains(err.Error(), "Index already exists") {
		pterm.Error.Println(err)
		return err
	}
	return nil
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
		values, ok := record.(map[string]interface{})
		if !ok {
			return errors.New("can't convert data")
		}
		id, ok := values["id"].(string)
		if !ok {
			return errors.New("can't assert ID")
		}
		_, err := c.client.Query(
			f.Create(
				f.Ref(f.Collection(table), hashIdentity(id)),
				f.Obj{"data": record},
			),
		)
		if err != nil {
			if strings.Contains(err.Error(), "instance already exists") {
				_, e := c.client.Query(
					f.Update(
						f.Ref(f.Collection(table), hashIdentity(id)),
						f.Obj{"data": record},
					),
				)
				if e != nil {
					return e
				}
			} else {
				return err
			}

		}
	}
	return nil
}
func hashIdentity(id string) uint32 {

	h := adler32.New()
	h.Write([]byte(id))
	return h.Sum32()

}