package cmd

import (
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"

	"cuelang.org/go/cue"
	"github.com/cueblox/blox"
	"github.com/cueblox/blox/internal/cuedb"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type exportSchemaCmd struct {
	cmd *cobra.Command
}

func newExportSchemaCmd() *exportSchemaCmd {
	root := &exportSchemaCmd{}
	cmd := &cobra.Command{
		Use:   "schema",
		Short: "Export Schemata",
		Long:  `Export schemata that are published via this repository`,
		Args:  cobra.NoArgs,
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

			/*remotes, err := cfg.GetList("remotes")
			if err == nil {
				cobra.CheckErr(parseRemotes(remotes))
			}
			*/

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

			for _, dataSet := range engine.GetDataSets() {
				// Walk each field and look for _id labels
				val := engine.Database.LookupPath(dataSet.GetDefinitionPath())

				fields, err := val.Fields(cue.All())
				if err != nil {
					cobra.CheckErr(err)
				}

				for fields.Next() {
				}
			}

			pterm.Println("Export")
		},
	}

	root.cmd = cmd
	return root
}
