package cmd

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/cueblox/blox/blox"
	"github.com/cueblox/blox/config"
	"github.com/cueblox/blox/cuedb"
	"github.com/cueblox/blox/encoding/markdown"
	"github.com/hashicorp/go-multierror"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// convertCmd represents the convert command
var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert markdown to yaml",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		cobra.CheckErr(err)

		database, err := cuedb.NewDatabase(cfg)
		cobra.CheckErr(err)

		// Load Schemas!
		cobra.CheckErr(database.RegisterTables(blox.ProfileCue))

		cobra.CheckErr(convertModels(&database))
	},
}

func convertModels(db *cuedb.Database) error {
	var errors error
	pterm.Info.Println("Converting Markdown files...")

	for _, table := range db.GetTables() {
		err := filepath.Walk(db.SourcePath(table),
			func(fpath string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if info.IsDir() {
					return err
				}

				ext := filepath.Ext(fpath)
				slug := strings.Replace(filepath.Base(fpath), ext, "", -1)

				if ext != ".md" && ext != ".mdx" {
					return nil
				}

				f, err := os.Open(fpath)
				if err != nil {
					errors = multierror.Append(errors, multierror.Prefix(err, fpath))
					return nil
				}

				bb, err := os.ReadFile(fpath)
				if err != nil {
					errors = multierror.Append(errors, multierror.Prefix(err, fpath))
					return nil
				}
				f.Close()

				md, err := markdown.ToYAML(string(bb))
				if err != nil {
					errors = multierror.Append(errors, multierror.Prefix(err, fpath))
					return nil
				}

				err = os.MkdirAll(path.Join(db.DestinationPath(table)), 0755)
				if err != nil {
					errors = multierror.Append(errors, multierror.Prefix(err, fpath))
					return nil
				}

				mdf, err := os.Create(fmt.Sprintf("%s.yaml", path.Join(db.DestinationPath(table), slug)))
				if err != nil {
					errors = multierror.Append(errors, multierror.Prefix(err, fpath))
					return nil
				}

				_, err = mdf.WriteString(md)
				if err != nil {
					errors = multierror.Append(errors, multierror.Prefix(err, fpath))
					return nil
				}
				mdf.Close()

				return nil
			},
		)

		if err != nil {
			errors = multierror.Append(err)
		}
	}

	if errors != nil {
		pterm.Error.Println("Conversions failed")
	} else {
		pterm.Success.Println("Conversions complete")
	}

	return errors
}

func init() {
	rootCmd.AddCommand(convertCmd)
}
