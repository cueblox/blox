package hosting

import (
	"sync"
)

var (
	providersMu sync.RWMutex
	providers   = make(map[string]Provider)
)

// Provider is an interface that all hosting
// providers will implement
type Provider interface {
	Name() string
	Description() string
	Install() error
}

// Register stores a provider in the map
func Register(name string, provider Provider) {
	providersMu.Lock()
	defer providersMu.Unlock()
	if provider == nil {
		panic("hosting: Register provider is nil")
	}
	if _, dup := providers[name]; dup {
		panic("hosting: Register called twice for provider " + name)
	}
	providers[name] = provider
}

// Providers returns the list of registered providers
func Providers() []Provider {
	providersMu.RLock()
	defer providersMu.RUnlock()
	list := make([]Provider, 0, len(providers))
	for name := range providers {
		list = append(list, providers[name])
	}

	return list
}

// GetProvider returns a single provider
// by name
func GetProvider(name string) Provider {
	providersMu.RLock()
	defer providersMu.RUnlock()

	if p, ok := providers[name]; ok {
		return p
	}
	return nil
}
