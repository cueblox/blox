package faunadb

import (
	"errors"
	"os"

	"github.com/cueblox/blox/internal/sync"
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
	key string
}

func (c *FaunaClient) Sync() error {
	return nil
}

func (c *FaunaClient) Help() string {
	return `Create a database in FaunaDB, then create a database specific
access key.  Export this key as FAUNA_KEY in your environment.`
}
