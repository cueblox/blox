package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/cueblox/blox/blox"
	"github.com/cueblox/blox/config"
	"github.com/cueblox/blox/cuedb"
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
		cfg, err := config.Load()
		cobra.CheckErr(err)

		database, err := cuedb.NewDatabase(cfg)
		cobra.CheckErr(err)

		// Load Schemas!
		cobra.CheckErr(database.RegisterTables(blox.ProfileCue))

		cobra.CheckErr(convertModels(&database))
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
		err := filepath.Walk(db.DestinationPath(table),
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if info.IsDir() {
					return err
				}

				ext := filepath.Ext(path)

				// if ext != cfg.DefaultExtension {
				// Should be SupportedExtensions?
				if ext != ".yaml" && ext != ".yml" {
					return err
				}

				slug := strings.Replace(filepath.Base(path), ext, "", -1)

				bytes, err := ioutil.ReadFile(path)
				if err != nil {
					return multierror.Append(err)
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
