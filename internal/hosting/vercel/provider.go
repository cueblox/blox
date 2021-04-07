package vercel

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
		internalName:        "vercel",
		internalDescription: "Vercel express api-only provider",
	}
	hosting.Register(vp.Name(), vp)
}

// Provider represents vercel
type Provider struct {
	internalName        string
	internalDescription string
}

// Name returns name
func (p *Provider) Name() string {
	return p.internalName
}

// Description returns description
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
	api := path.Join(root, "api")
	err = os.MkdirAll(api, 0755)
	if err != nil {
		return err
	}
	// create index.js in api directory

	index := path.Join(api, "index.js")
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
	// create vercel.json
	vc := path.Join(root, "vercel.json")
	err = hosting.CreateFileWithContents(vc, verceljson)
	if err != nil {
		return err
	}
	pterm.Info.Println("Vercel provider installed.")
	pterm.Info.Println("Run `npm install` to install dependencies.")

	return nil
}

//go:embed index.js
var indexjs string

//go:embed vercel.json
var verceljson string

//go:embed package.json
var packagejson string
