package cmd

import (
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
	Short: "Convert, Validate, & Build Your JSON Blox",
	Run: func(cmd *cobra.Command, args []string) {
		database, err := cuedb.NewDatabase()
		cobra.CheckErr(err)

		// Load Schemas!
		schemaDir, err := database.GetConfigString("schema_dir")
		cobra.CheckErr(err)

		err = filepath.WalkDir(schemaDir, func(path string, d fs.DirEntry, err error) error {
			if !d.IsDir() {
				bb, err := os.ReadFile(path)
				if err != nil {
					return err
				}

				err = database.RegisterTables(string(bb))
				if err != nil {
					return err
				}
			}
			return nil
		})
		cobra.CheckErr(err)

		cobra.CheckErr(buildModels(&database))

		if referentialIntegrity {
			pterm.Info.Println("Checking Referential Integrity")
			err = database.ReferentialIntegrity()
			if err != nil {
				pterm.Error.Println(err)
			} else {
				pterm.Success.Println("Foreign Keys Validated")
			}
		}

		jso, err := database.MarshalJSON()
		cobra.CheckErr(err)

		buildDir, err := database.GetConfigString("build_dir")
		cobra.CheckErr(err)
		err = os.MkdirAll(buildDir, 0755)
		cobra.CheckErr(err)
		filename := "data.json"
		filePath := path.Join(buildDir, filename)
		err = os.WriteFile(filePath, jso, 0755)
		cobra.CheckErr(err)

	},
}

func buildModels(db *cuedb.Database) error {
	var errors error

	pterm.Info.Println("Validating ...")

	for _, table := range db.GetTables() {
		err := os.MkdirAll(db.GetTableDataDir(table), 0755)
		if err != nil {
			errors = multierror.Append(err)
			continue
		}

		err = filepath.Walk(db.GetTableDataDir(table),
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if info.IsDir() {
					return err
				}

				ext := strings.TrimPrefix(filepath.Ext(path), ".")

				if !table.IsSupportedExtension(ext) {
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

				err = db.Insert(table, record)
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
