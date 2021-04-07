package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/cueblox/blox/internal/blox"
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
		cobra.CheckErr(database.RegisterTables(blox.ProfileCue))

		// cobra.CheckErr(convertModels(&database))
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

		fmt.Println("I should write this to a file")
		fmt.Println(string(jso))

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

				slug := strings.Replace(filepath.Base(path), ext, "", -1)

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
