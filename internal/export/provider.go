package export

import (
	"encoding/json"
	"fmt"
	"sort"
	"sync"

	"github.com/pterm/pterm"
)

var (
	providersMu sync.RWMutex
	providers   = make(map[string]ExportProvider)
)

type ExportProvider interface {
	Initialize() (EngineProvider, error)
	Name() string
}

type EngineProvider interface {
	Sync([]byte) error
	Help() string
}

func Register(name string, provider ExportProvider) {
	providersMu.Lock()
	defer providersMu.Unlock()
	if provider == nil {
		panic("sync: Register provider is nil")
	}
	if _, dup := providers[name]; dup {
		panic("sync: Register called twice for provider " + name)
	}
	providers[name] = provider
}

// Providers returns a sorted list of the names of the registered providers.
func Providers() []string {
	providersMu.RLock()
	defer providersMu.RUnlock()
	list := make([]string, 0, len(providers))
	for name := range providers {
		list = append(list, name)
	}
	sort.Strings(list)
	return list
}

func Open(providerName string) (*ExportEngine, error) {
	providersMu.RLock()
	provideri, ok := providers[providerName]
	providersMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("export: unknown provider %q (forgotten import?)", providerName)
	}

	engine, err := provideri.Initialize()
	if err != nil {
		return nil, err
	}
	pterm.Info.Printf("Created Export Provider for %s\n", provideri.Name())
	return &ExportEngine{p: engine}, nil
}

type ExportEngine struct {
	p EngineProvider
}

func (e *ExportEngine) Synchronize(bb []byte) error {
	return e.p.Sync(bb)
}

func (e *ExportEngine) Help() string {
	return e.p.Help()
}

func MakeMap(bb []byte) (map[string]interface{}, error) {
	var data = make(map[string]interface{})
	err := json.Unmarshal(bb, &data)
	return data, err

}

func GetTables(m map[string]interface{}) []string {

	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
