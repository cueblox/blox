package blox

import (
	// import go:embed
	_ "embed"
	"io/ioutil"
	"path/filepath"
	"strings"

	"cuelang.org/go/cuego"
	"github.com/devrel-blox/blox/cueutils"
	"github.com/goccy/go-yaml"
)

// ArticleCue is the cue for an Article
//go:embed article.cue
var ArticleCue string

// ArticleTemplate is the template for an article
//go:embed article.md
var ArticleTemplate string

// CategoryCue is the Cue for a category
//go:embed category.cue
var CategoryCue string

// CategoryTemplate is the template for a Category
//go:embed category.md
var CategoryTemplate string

// ProfileCue is the cue for a profile
//go:embed profile.cue
var ProfileCue string

// ProfileTemplate is the template for a profile
//go:embed profile.md
var ProfileTemplate string

// PageCue is the cue template for a page
//go:embed page.cue
var PageCue string

// PageTemplate is the page template
//go:embed page.md
var PageTemplate string

// FromYAML converts converted YAML content into output-ready map
func FromYAML(path string, modelName string, cue string) (map[string]interface{}, error) {
	var model = make(map[string]interface{})

	cuego.DefaultContext = &cuego.Context{}

	err := cuego.Constrain(&model, cue)
	if err != nil {
		return nil, cueutils.UsefulError(err)
	}

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, cueutils.UsefulError(err)
	}

	err = yaml.Unmarshal(bytes, &model)
	if err != nil {
		return nil, cueutils.UsefulError(err)
	}

	err = cuego.Complete(&model)
	if err != nil {
		return nil, cueutils.UsefulError(err)
	}

	ext := filepath.Ext(path)
	slug := strings.Replace(filepath.Base(path), ext, "", -1)

	model["id"] = slug

	return model, nil
}
