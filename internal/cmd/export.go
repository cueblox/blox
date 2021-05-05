package cmd

import (
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cueblox/blox"
	"github.com/cueblox/blox/internal/cuedb"
	"github.com/cueblox/blox/internal/export"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	_ "github.com/cueblox/blox/internal/export/faunadb"
)

type exportCmd struct {
	cmd *cobra.Command
}

func newExportCmd() *exportCmd {
	root := &exportCmd{}
	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export blox dataset",
		Long:  `The sync command allows you to export your blox dataset using various providers.`,
		Run: func(cmd *cobra.Command, args []string) {
			pterm.Info.Println("Building Dataset")
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
			pterm.Info.Println("Synchronization started")
			cobra.CheckErr(synchronizeDataset(jso))
		},
	}
	cmd.AddCommand(
		newExportProvidersCmd().cmd,
	)
	cmd.Flags().BoolVarP(&referentialIntegrity, "referential-integrity", "i", false, "Verify referential integrity")

	root.cmd = cmd
	return root
}

func synchronizeDataset(jsn []byte) error {
	engine, err := export.Open("faunadb")
	cobra.CheckErr(err)
	err = engine.Synchronize(jsn)
	if err != nil {
		pterm.Error.Println(engine.Help())
		return err
	}
	return nil
}
