package content

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/cueblox/blox"
	"github.com/graphql-go/graphql"

	"github.com/cueblox/blox/internal/cuedb"
	"github.com/cueblox/blox/internal/encoding/markdown"
	"github.com/cueblox/blox/internal/repository"
	"github.com/goccy/go-yaml"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-plugin"
	"github.com/pterm/pterm"
)

type Service struct {
	engine    *cuedb.Engine
	Cfg       *blox.Config
	rawConfig string
	ri        bool
	schema    *graphql.Schema
	built     bool
}

var prePluginMap map[string]plugin.Plugin
var postPluginMap map[string]plugin.Plugin

func NewService(bloxConfig string, referentialIntegrity bool) (*Service, error) {
	cfg, err := blox.NewConfig(BaseConfig)
	if err != nil {
		return nil, err
	}

	err = cfg.LoadConfigString(bloxConfig)
	if err != nil {
		return nil, err
	}

	engine, err := cuedb.NewEngine()
	if err != nil {
		return nil, err
	}

	schemataDir, err := cfg.GetString("schemata_dir")
	if err != nil {
		return nil, err
	}
	err = filepath.WalkDir(schemataDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			bb, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			pterm.Debug.Printf("\t\tAttempting to register schema: %s\n", path)
			err = engine.RegisterSchema(string(bb))
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	prePluginMap = make(map[string]plugin.Plugin)
	postPluginMap = make(map[string]plugin.Plugin)
	return &Service{
		engine:    engine,
		Cfg:       cfg,
		rawConfig: bloxConfig,
		ri:        referentialIntegrity,
	}, nil
}

const BaseConfig = `{
	#Remote: {
	    name: string
	    version: string
	    repository: string
	}
	#Plugin: {
		name: string
		executable: string
	}
	build_dir:    string | *"_build"
	data_dir:     string | *"data"
	schemata_dir: string | *"schemata"
	static_dir: string | *"static"
	template_dir: string | *"templates"
	remotes: [ ...#Remote ]
	prebuild: [...#Plugin]
	postbuild: [...#Plugin]
}`

func (s *Service) build() error {
	var errors error

	err := s.parseRemotes()
	if err != nil {
		return err
	}

	err = s.runPrePlugins()
	if err != nil {
		return err
	}

	for _, dataSet := range s.engine.GetDataSets() {
		pterm.Debug.Printf("\t\tBuilding Dataset: %s\n", dataSet.ID())

		// We're using the Or variant of GetString because we know this call can't
		// fail, as the config isn't valid without.
		dataSetDirectory := fmt.Sprintf("%s/%s", s.Cfg.GetStringOr("data_dir", ""), dataSet.GetDataDirectory())

		err := os.MkdirAll(dataSetDirectory, 0o755)
		if err != nil {
			errors = multierror.Append(err)
			continue
		}

		err = filepath.Walk(dataSetDirectory,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if info.IsDir() {
					return err
				}

				ext := strings.TrimPrefix(filepath.Ext(path), ".")

				if !dataSet.IsSupportedExtension(ext) {
					return nil
				}
				relative, err := filepath.Rel(dataSetDirectory, path)
				if err != nil {
					return err
				}
				slug := strings.TrimSuffix(relative, "."+ext)
				pterm.Debug.Println(slug)
				bytes, err := ioutil.ReadFile(path)
				if err != nil {
					return multierror.Append(err)
				}

				// Loaders to get to YAML
				// We should offer various, simple for now with markdown
				mdStr := ""
				if ext == "md" || ext == "mdx" {
					mdStr, err = markdown.ToYAML(string(bytes))
					if err != nil {
						return err
					}

					bytes = []byte(mdStr)
				}

				istruct := make(map[string]interface{})

				err = yaml.Unmarshal(bytes, &istruct)
				if err != nil {
					return multierror.Append(err)
				}

				record := make(map[string]interface{})
				record[slug] = istruct

				err = s.engine.Insert(dataSet, record)
				if err != nil {
					return multierror.Append(err)
				}

				return err
			},
		)

		if err != nil {
			errors = multierror.Append(err)
		}
	}

	if s.ri {
		err := s.engine.ReferentialIntegrity()
		if err != nil {
			errors = multierror.Append(err)
		}
	}

	if errors != nil {
		pterm.Error.Println("Validation Failed")
		return errors
	}

	pterm.Success.Println("Validation Complete")
	s.built = true
	return nil
}

func (s *Service) parseRemotes() error {
	remotes, err := s.Cfg.GetList("remotes")
	if err != nil {
		return err
	}
	iter, err := remotes.List()
	if err != nil {
		return err
	}
	for iter.Next() {
		val := iter.Value()

		//nolint
		name, err := val.FieldByName("name", false)
		if err != nil {
			return err
		}
		n, err := name.Value.String()
		if err != nil {
			return err
		}
		//nolint
		version, err := val.FieldByName("version", false)
		if err != nil {
			return err
		}
		v, err := version.Value.String()
		if err != nil {
			return err
		}
		//nolint
		repository, err := val.FieldByName("repository", false)
		if err != nil {
			return err
		}
		r, err := repository.Value.String()
		if err != nil {
			return err
		}
		err = s.ensureRemote(n, v, r)
		if err != nil {
			return err
		}
	}
	return nil
}

// ensureRemote downloads a remote schema at a specific
// version if it doesn't exist locally
func (s *Service) ensureRemote(name, version, repo string) error {
	// Load Schemas!
	schemataDir, err := s.Cfg.GetString("schemata_dir")
	if err != nil {
		return err
	}

	schemaFilePath := path.Join(schemataDir, fmt.Sprintf("%s_%s.cue", name, version))
	_, err = os.Stat(schemaFilePath)
	if os.IsNotExist(err) {
		pterm.Info.Printf("Schema does not exist locally: %s_%s.cue\n", name, version)
		manifest := fmt.Sprintf("https://%s/manifest.json", repo)
		res, err := http.Get(manifest)
		if err != nil {
			return err
		}

		var repos repository.Repository
		err = json.NewDecoder(res.Body).Decode(&repos)
		if err != nil {
			return err
		}

		var selectedVersion *repository.Version
		for _, s := range repos.Schemas {
			if s.Name == name {
				for _, v := range s.Versions {
					if v.Name == version {
						selectedVersion = v
						pterm.Debug.Println(v.Name, v.Schema, v.Definition)
					}
				}
			}
		}

		// make schemata directory
		err = os.MkdirAll(schemataDir, 0o755)
		if err != nil {
			return err
		}

		// TODO: don't overwrite each time
		filename := fmt.Sprintf("%s_%s.cue", name, version)
		filePath := path.Join(schemataDir, filename)
		err = os.WriteFile(filePath, []byte(selectedVersion.Definition), 0o755)
		if err != nil {
			return err
		}
		pterm.Info.Printf("Schema downloaded: %s_%s.cue\n", name, version)
		return nil
	}
	pterm.Info.Println("Schema already exists locally, skipping download.")
	return nil
}
