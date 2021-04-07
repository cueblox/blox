package blox

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/devrel-blox/blox/config"
	"github.com/spf13/cobra"
)

// Models is used by various commands to determine
// how to perform certain actions based on arguments
// and flags provided. All new types must be
// represented in this slice.
var Models = []Model{
	{
		ID:         "profile",
		Name:       "Profile",
		Folder:     "profiles",
		ForeignKey: "profile_id",
		Cue:        ProfileCue,
	},
	{
		ID:         "article",
		Name:       "Article",
		Folder:     "articles",
		ForeignKey: "article_id",
		Cue:        ArticleCue,
	},
	{
		ID:         "category",
		Name:       "Category",
		Folder:     "categories",
		ForeignKey: "category_id",
		Cue:        CategoryCue,
	},
	{
		ID:         "page",
		Name:       "Page",
		Folder:     "pages",
		ForeignKey: "page_id",
		Cue:        PageCue,
	},
}

// Model represents a content model
type Model struct {
	ID         string
	Name       string
	Folder     string
	ForeignKey string
	Cue        string
}

// GetModel finds a Model definition and returns
// it to the caller.
func GetModel(id string) (Model, error) {
	for _, m := range Models {
		if m.ID == id {
			return m, nil
		}
	}
	return Model{}, errors.New("model not found")
}

// StaticContentPath returns the path where
// images and other static content path will
// be stored
func (m Model) StaticContentPath() string {
	cfg, err := config.Load()
	cobra.CheckErr(err)
	return path.Join(cfg.Base, cfg.Static)
}

// SourceContentPath returns the path where
// content source files (markdown/yaml) are stored
func (m Model) SourceContentPath() string {
	cfg, err := config.Load()
	cobra.CheckErr(err)
	return path.Join(cfg.Base, cfg.Source, m.Folder)
}

// DestinationContentPath returns the path where
// output YAML is stored
func (m Model) DestinationContentPath() string {
	cfg, err := config.Load()
	cobra.CheckErr(err)
	return path.Join(cfg.Base, cfg.Destination, m.Folder)
}

// SourceFilePath returns the path where a specific
// content file should be stored, based on its slug
func (m Model) SourceFilePath(slug string) string {
	cfg, err := config.Load()
	cobra.CheckErr(err)
	fileName := slug + cfg.DefaultExtension

	return path.Join(cfg.Base, cfg.Source, m.Folder, fileName)
}

// TemplatePath returns the path where templates are stored
func (m Model) TemplatePath() string {
	cfg, err := config.Load()
	cobra.CheckErr(err)
	return path.Join(cfg.Base, cfg.Templates, m.Folder)
}

// TemplateFilePath returns the path for a template file
func (m Model) TemplateFilePath(slug string) string {
	cfg, err := config.Load()
	cobra.CheckErr(err)
	fileName := slug + cfg.DefaultExtension

	return path.Join(cfg.Base, cfg.Templates, m.Folder, fileName)
}

// DestinationFilePath returns the path where a converted
// model will be stored as YAML.
func (m Model) DestinationFilePath(slug string) string {
	cfg, err := config.Load()
	cobra.CheckErr(err)
	fileName := slug + ".yaml"

	return path.Join(cfg.Base, cfg.Destination, m.Folder, fileName)
}

// New creates a new content model
func (m Model) New(slug string, destination string) error {
	err := os.MkdirAll(destination, 0744)
	if err != nil {
		return err
	}

	// check for user installed templates first
	cfg, err := config.Load()
	cobra.CheckErr(err)

	templatePath := path.Join(cfg.Base, cfg.Templates, m.Folder, m.ID+cfg.DefaultExtension)

	joined := path.Join(destination, slug)
	// check to see if we're creating the templates
	// from the `init` command
	var bb []byte
	if templatePath == joined {
		bb, err = m.defaultTemplate()
		cobra.CheckErr(err)
	} else {
		bb, err = os.ReadFile(templatePath)
		cobra.CheckErr(err)
	}

	// create the destination file
	f, err := os.Create(joined)
	if err != nil {
		return err
	}
	defer f.Close()

	f.Write(bb)

	return nil
}

func (m Model) defaultTemplate() ([]byte, error) {

	switch m.ID {
	case "article":
		return []byte(ArticleTemplate), nil
	case "category":
		return []byte(CategoryTemplate), nil
	case "profile":
		return []byte(ProfileTemplate), nil
	case "page":
		return []byte(PageTemplate), nil
	default:
		return []byte{}, fmt.Errorf("generator doesn't support %s yet", m.ID)
	}
}

// BaseModel defines fields used by all drb
// models
type BaseModel struct {
	ID      string `json:"id"`
	Body    string `json:"body"`
	BodyRaw string `json:"body_raw"`
}
