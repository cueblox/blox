package cmd

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/cueblox/blox"
	"github.com/cueblox/blox/internal/cuedb"
	"github.com/cueblox/blox/internal/encoding/markdown"
	"github.com/goccy/go-yaml"
	"github.com/hashicorp/go-multierror"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var (
	referentialIntegrity bool
)

const DefaultConfigName = "blox.cue"

const BaseConfig = `{
    build_dir:    string | *"_build"
    data_dir:     string | *"data"
    schemata_dir: string | *"schemata"
}`

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Validate & Build",
	Long: `The build command will ensure that your content is correct by
validating it against your schemata. Once validated, it will render all
your content into a single JSON file, which can be consumed by your tooling
of choice.

Referential Integrity can be enforced with -i. This ensures that any fields
ending with _id are valid references to identifies within the other content type.`,
	Run: func(cmd *cobra.Command, args []string) {
		userConfig, err := ioutil.ReadFile("blox.cue")
		cobra.CheckErr(err)

		engine, err := cuedb.NewEngine()
		cobra.CheckErr(err)

		cfg, err := blox.NewConfig(BaseConfig)
		cobra.CheckErr(err)

		err = cfg.LoadConfigString(string(userConfig))
		cobra.CheckErr(err)

		// Load Schemas!
		schemataDir, err := cfg.GetString("schemata_dir")
		cobra.CheckErr(err)
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
		cobra.CheckErr(os.MkdirAll(buildDir, 0755))

		filename := "data.json"
		filePath := path.Join(buildDir, filename)
		cobra.CheckErr(os.WriteFile(filePath, jso, 0755))
		pterm.Success.Printf("Data blox written to '%s'\n", filePath)
	},
}

func buildDataSets(engine *cuedb.Engine, cfg *blox.Config) error {
	var errors error

	for _, dataSet := range engine.GetDataSets() {
		pterm.Debug.Printf("\t\tBuilding Dataset: %s\n", dataSet.ID())

		// We're using the Or variant of GetString because we know this call can't
		// fail, as the config isn't valid without.
		dataSetDirectory := fmt.Sprintf("%s/%s", cfg.GetStringOr("data_dir", ""), dataSet.GetDataDirectory())

		err := os.MkdirAll(dataSetDirectory, 0755)
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

				slug := strings.TrimSuffix(filepath.Base(path), "."+ext)

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

				var istruct = make(map[string]interface{})

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

func init() {
	rootCmd.AddCommand(buildCmd)
	buildCmd.Flags().BoolVarP(&referentialIntegrity, "referential-integrity", "i", false, "Verify referential integrity")
}
