package netlify

import (
	// import go:embed
	_ "embed"
	"os"
	"path"

	"github.com/devrel-blox/blox/hosting"
	"github.com/pterm/pterm"
)

func init() {
	vp := &Provider{
		internalName:        "netlify",
		internalDescription: "Netlify express api-only provider",
	}
	hosting.Register(vp.Name(), vp)
}

// Provider represents the Netlify provider
type Provider struct {
	internalName        string
	internalDescription string
}

// Name is the name
func (p *Provider) Name() string {
	return p.internalName
}

// Description is the description
func (p *Provider) Description() string {
	return p.internalDescription
}

// Install installs netlify files
func (p *Provider) Install() error {
	// make api directory
	root, err := os.Getwd()
	if err != nil {
		return err
	}
	api := path.Join(root, "api")
	err = os.MkdirAll(api, 0755)
	if err != nil {
		return err
	}
	// create index.js in api directory

	index := path.Join(api, "index.mjs")
	err = hosting.CreateFileWithContents(index, indexjs)
	if err != nil {
		return err
	}

	// create package.json
	pkg := path.Join(root, "package.json")
	err = hosting.CreateFileWithContents(pkg, packagejson)
	if err != nil {
		return err
	}
	// create netlify.toml
	vc := path.Join(root, "netlify.toml")
	err = hosting.CreateFileWithContents(vc, netlifyToml)
	if err != nil {
		return err
	}
	pterm.Info.Println("Netlify provider installed.")
	pterm.Info.Println("Run `npm install` to install dependencies.")

	return nil
}

//go:embed index.mjs
var indexjs string

//go:embed netlify.toml
var netlifyToml string

//go:embed package.json
var packagejson string
