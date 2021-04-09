package azure

import (
	// import go:embed
	_ "embed"
	"os"
	"path"

	"github.com/cueblox/blox/internal/hosting"
	"github.com/pterm/pterm"
)

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

// Install installs vercel hosting
func (p *Provider) Install() error {
	// make api directory
	root, err := os.Getwd()
	if err != nil {
		return err
	}

	// create azure/
	azure := path.Join(root, "azure")
	err = os.MkdirAll(azure, 0755)
	if err != nil {
		return err
	}

	// azure/api
	api := path.Join(azure, "api")
	err = os.MkdirAll(api, 0755)
	if err != nil {
		return err
	}

	// work around ASWA limition of number of files
	// in static site directory by putting the
	// api in a different root (above)
	// azure/app
	app := path.Join(azure, "app")
	err = os.MkdirAll(app, 0755)
	if err != nil {
		return err
	}

	htmlindex := path.Join(app, "index.html")
	err = hosting.CreateFileWithContents(htmlindex, indexhtml)
	if err != nil {
		return err
	}
	cssstyle := path.Join(app, "styles.css")
	err = hosting.CreateFileWithContents(cssstyle, stylescss)
	if err != nil {
		return err
	}

	// azure/api/graphql
	graphql := path.Join(api, "GraphQL")
	err = os.MkdirAll(graphql, 0755)
	if err != nil {
		return err
	}

	// create index.js in api directory

	index := path.Join(graphql, "index.js")
	err = hosting.CreateFileWithContents(index, indexjs)
	if err != nil {
		return err
	}

	// create function.json
	fjson := path.Join(graphql, "function.json")
	err = hosting.CreateFileWithContents(fjson, functionjson)
	if err != nil {
		return err
	}

	// create package.json
	pkg := path.Join(api, "package.json")
	err = hosting.CreateFileWithContents(pkg, packagejson)
	if err != nil {
		return err
	}

	// create host.json
	host := path.Join(api, "host.json")
	err = hosting.CreateFileWithContents(host, hostjson)
	if err != nil {
		return err
	}
	// create host.json
	proxies := path.Join(api, "proxies.json")
	err = hosting.CreateFileWithContents(proxies, proxiesjson)
	if err != nil {
		return err
	}

	pterm.Info.Println("Azure provider installed.")
	pterm.Info.Println("Run `npm install` to install dependencies.")

	return nil
}

//go:embed index.js
var indexjs string

//go:embed proxies.json
var proxiesjson string

//go:embed package.json
var packagejson string

//go:embed host.json
var hostjson string

//go:embed function.json
var functionjson string

//go:embed index.html
var indexhtml string

//go:embed styles.css
var stylescss string
