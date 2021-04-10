package cmd

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

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

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Validate and build your data",
	Run: func(cmd *cobra.Command, args []string) {
		// This can happen at a global Cobra level, if I knew how
		config, err := cuedb.NewRuntime()
		cobra.CheckErr(err)
		cobra.CheckErr(config.LoadConfig())

		runtime, err := cuedb.NewRuntime()
		cobra.CheckErr(err)

		// Load Schemas!
		schemaDir, err := config.GetString("schema_dir")
		cobra.CheckErr(err)

		err = filepath.WalkDir(schemaDir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() {
				bb, err := os.ReadFile(path)
				if err != nil {
					return err
				}

				err = runtime.RegisterSchema(string(bb))
				if err != nil {
					return err
				}
			}
			return nil
		})
		cobra.CheckErr(err)

		cobra.CheckErr(buildModels(&config, &runtime))

		if referentialIntegrity {
			pterm.Info.Println("Checking Referential Integrity")
			err = runtime.ReferentialIntegrity()
			if err != nil {
				pterm.Error.Println(err)
			} else {
				pterm.Success.Println("Foreign Keys Validated")
			}
		}

		output, err := runtime.GetOutput()
		cobra.CheckErr(err)

		jso, err := output.MarshalJSON()
		cobra.CheckErr(err)

		buildDir, err := config.GetString("build_dir")
		cobra.CheckErr(err)
		err = os.MkdirAll(buildDir, 0755)
		cobra.CheckErr(err)
		filename := "data.json"
		filePath := path.Join(buildDir, filename)
		err = os.WriteFile(filePath, jso, 0755)
		cobra.CheckErr(err)

	},
}

func buildModels(config *cuedb.Runtime, runtime *cuedb.Runtime) error {
	var errors error

	pterm.Info.Println("Validating ...")

	for _, dataSet := range runtime.GetDataSets() {
		// We're using the Or variant of GetString because we know this call can't
		// fail, as the config isn't valid without.
		dataSetDirectory := fmt.Sprintf("%s/%s", config.GetStringOr("data_dir", ""), dataSet.GetDataDirectory())

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

				err = runtime.Insert(dataSet, record)
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
		pterm.Error.Println("Validations failed")
	} else {
		pterm.Success.Println("Validations complete")
	}

	return errors
}

func init() {
	rootCmd.AddCommand(buildCmd)

	buildCmd.Flags().BoolVarP(&referentialIntegrity, "referential-integrity", "i", false, "Enforce referential integrity")
}
