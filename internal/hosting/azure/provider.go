package azure

import "github.com/cueblox/blox/internal/hosting"

func init() {
	vp := &Provider{
		internalName:        "azure",
		internalDescription: "Azure static web apps provider",
	}
	hosting.Register(vp.Name(), vp)
}

// Provider represents Azure
type Provider struct {
	internalName        string
	internalDescription string
}

// Name is the name of the provider
func (p *Provider) Name() string {
	return p.internalName
}

// Description is the description of the provider
func (p *Provider) Description() string {
	return p.internalDescription
}

// Install installs the provider
func (p *Provider) Install() error {
	return nil
}
