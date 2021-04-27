package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"cuelang.org/go/cue"
	"github.com/cueblox/blox"
	"github.com/cueblox/blox/internal/cuedb"
	"github.com/cueblox/blox/internal/encoding/markdown"
	"github.com/disintegration/imaging"
	"github.com/goccy/go-yaml"
	"github.com/h2non/filetype"
	"github.com/hashicorp/go-multierror"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type bloxBuildCmd struct {
	cmd *cobra.Command
}

func newBloxBuildCmd() *bloxBuildCmd {
	root := &bloxBuildCmd{}
	cmd := &cobra.Command{
		Use:   "build",
		Short: "Validate & Build dataset",
		Long: `The build command will ensure that your dataset is correct by
	validating it against your schemata. Once validated, it will render all
	your content into a single JSON file, which can be consumed by your tooling
	of choice.
	
	Referential Integrity can be enforced with -i. This ensures that any fields
	ending with _id are valid references to identifiers within the other content type.`,
		Run: func(cmd *cobra.Command, args []string) {
			userConfig, err := ioutil.ReadFile("blox.cue")

			pterm.Debug.Printf("loading user config")

			cobra.CheckErr(err)

			engine, err := cuedb.NewEngine()

			pterm.Debug.Printf("new engine")
			cobra.CheckErr(err)

			cfg, err := blox.NewConfig(BaseConfig)

			pterm.Debug.Printf("newConfig")
			cobra.CheckErr(err)

			err = cfg.LoadConfigString(string(userConfig))
			cobra.CheckErr(err)

			// Load Schemas!
			schemataDir, err := cfg.GetString("schemata_dir")
			cobra.CheckErr(err)

			remotes, err := cfg.GetList("remotes")
			if err == nil {
				cobra.CheckErr(parseRemotes(remotes))
			}

			err = processImages(cfg)
			if err != nil {
				cobra.CheckErr(err)
			}
			pterm.Debug.Printf("\t\tUsing schemata from: %s\n", schemataDir)

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
			cobra.CheckErr(err)

			pterm.Debug.Println("\t\tBuilding DataSets")
			cobra.CheckErr(buildDataSets(engine, cfg))

			if referentialIntegrity {
				pterm.Info.Println("Verifying Referential Integrity")
				cobra.CheckErr(engine.ReferentialIntegrity())
				pterm.Success.Println("Referential Integrity OK")
			}

			pterm.Debug.Println("Building output data blox")
			output, err := engine.GetOutput()
			cobra.CheckErr(err)

			pterm.Debug.Println("Rendering data blox to JSON")
			jso, err := output.MarshalJSON()
			cobra.CheckErr(err)

			buildDir, err := cfg.GetString("build_dir")
			cobra.CheckErr(err)
			cobra.CheckErr(os.MkdirAll(buildDir, 0o755))

			filename := "data.json"
			filePath := path.Join(buildDir, filename)
			cobra.CheckErr(os.WriteFile(filePath, jso, 0o755))
			pterm.Success.Printf("Data blox written to '%s'\n", filePath)
		},
	}
	cmd.Flags().BoolVarP(&referentialIntegrity, "referential-integrity", "i", false, "Verify referential integrity")

	root.cmd = cmd
	return root
}

var referentialIntegrity bool

const DefaultConfigName = "blox.cue"

const BaseConfig = `{
    #Remote: {
        name: string
        version: string
        repository: string
    }
    build_dir:    string | *"_build"
    data_dir:     string | *"data"
    schemata_dir: string | *"schemata"
	static_dir: string | *"static"
	template_dir: string | *"templates"
	remotes: [ ...#Remote ]

}`

func buildDataSets(engine *cuedb.Engine, cfg *blox.Config) error {
	var errors error

	for _, dataSet := range engine.GetDataSets() {
		pterm.Debug.Printf("\t\tBuilding Dataset: %s\n", dataSet.ID())

		// We're using the Or variant of GetString because we know this call can't
		// fail, as the config isn't valid without.
		dataSetDirectory := fmt.Sprintf("%s/%s", cfg.GetStringOr("data_dir", ""), dataSet.GetDataDirectory())

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

				err = engine.Insert(dataSet, record)
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

	if errors != nil {
		pterm.Error.Println("Validation Failed")
		return errors
	}

	pterm.Success.Println("Validation Complete")
	return nil
}

func parseRemotes(value cue.Value) error {
	iter, err := value.List()
	if err != nil {
		return err
	}
	for iter.Next() {
		val := iter.Value()
		name, err := val.FieldByName("name", false)
		if err != nil {
			return err
		}
		n, err := name.Value.String()
		if err != nil {
			return err
		}
		version, err := val.FieldByName("version", false)
		if err != nil {
			return err
		}
		v, err := version.Value.String()
		if err != nil {
			return err
		}
		repository, err := val.FieldByName("repository", false)
		if err != nil {
			return err
		}
		r, err := repository.Value.String()
		if err != nil {
			return err
		}
		err = ensureRemote(n, v, r)
		if err != nil {
			return err
		}
	}
	return nil
}

// processImages scans the static dir for images
// when it finds an image, it reads the image metadata
// and saves a corresponding YAML file describing the image
// in the 'images' data directory.
func processImages(cfg *blox.Config) error {
	staticDir, err := cfg.GetString("static_dir")
	if err != nil {
		pterm.Info.Printf("no static directory present, skipping image linking")
	}
	cobra.CheckErr(err)
	pterm.Info.Printf("processing images in %s\n", staticDir)
	fi, err := os.Stat(staticDir)
	cobra.CheckErr(err)
	if !fi.IsDir() {
		return errors.New("given static directory is not a directory")
	}
	imagesDirectory := filepath.Join(staticDir, "images")
	err = filepath.Walk(imagesDirectory,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			pterm.Debug.Printf("\t\tProcessing %s\n", path)
			if !info.IsDir() {
				buf, err := ioutil.ReadFile(path)
				if err != nil {
					return err
				}
				if filetype.IsImage(buf) {

					src, err := imaging.Open(path)
					if err != nil {
						return err
					}

					relpath, err := filepath.Rel(staticDir, path)
					if err != nil {
						return err
					}
					pterm.Debug.Printf("\t\tFile is an image: %s\n", relpath)
					kind, err := filetype.Match(buf)
					if err != nil {
						return err
					}
					pterm.Debug.Printf("\t\tFile type: %s. MIME: %s\n", kind.Extension, kind.MIME.Value)
					if err != nil {
						return err
					}

					bi := &BloxImage{
						FileName: relpath,
						Height:   src.Bounds().Dy(),
						Width:    src.Bounds().Dx(),
					}
					bytes, err := yaml.Marshal(bi)
					if err != nil {
						return err
					}
					dataDir, err := cfg.GetString("data_dir")
					if err != nil {
						return err
					}

					ext := strings.TrimPrefix(filepath.Ext(relpath), ".")
					slug := strings.TrimSuffix(relpath, "."+ext)

					outputPath := filepath.Join(dataDir, slug+".yaml")
					err = os.MkdirAll(filepath.Dir(outputPath), 0o755)
					if err != nil {
						pterm.Error.Println(err)
						return err
					}
					// only write the yaml file if it doesn't exist.
					// don't overwrite existing records.
					_, err = os.Stat(outputPath)
					if err != nil && errors.Is(err, os.ErrNotExist) {
						err = os.WriteFile(outputPath, bytes, 0o755)
						if err != nil {
							pterm.Error.Println(err)
							return err
						}
					}
				} else {
					pterm.Debug.Printf("File is not an image: %s\n", path)
				}
			}

			return nil
		})
	cobra.CheckErr(err)
	return nil
}

type BloxImage struct {
	FileName string `yaml:"file_name"`
	Height   int    `yaml:"height"`
	Width    int    `yaml:"width"`
}
