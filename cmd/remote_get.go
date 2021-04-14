package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/cueblox/blox"
	"github.com/cueblox/blox/internal/repository"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var remoteGetCmd = &cobra.Command{
	Use:   "get <repository> <schema name> <version>",
	Short: "Add a remote schema to your repository",
	Args:  cobra.ExactArgs(3),

	Run: func(cmd *cobra.Command, args []string) {
		manifest := fmt.Sprintf("https://%s/manifest.json", args[0])
		res, err := http.Get(manifest)
		cobra.CheckErr(err)

		var repos repository.Repository
		json.NewDecoder(res.Body).Decode(&repos)

		schemaName := args[1]
		version := args[2]

		var selectedVersion *repository.Version
		for _, s := range repos.Schemas {
			if s.Name == schemaName {
				for _, v := range s.Versions {
					if v.Name == version {
						selectedVersion = v
						pterm.Debug.Println(v.Name, v.Schema, v.Definition)
					}
				}
			}
		}

		cfg, err := blox.NewConfig(BaseConfig)
		schemataDir := cfg.GetStringOr("schemata_dir", "schemata")
		cobra.CheckErr(os.MkdirAll(schemataDir, 0755))

		filename := fmt.Sprintf("%s_%s.cue", schemaName, version)
		filePath := path.Join(schemataDir, filename)
		cobra.CheckErr(os.WriteFile(filePath, []byte(selectedVersion.Definition), 0755))

	},
}

func init() {
	remoteCmd.AddCommand(remoteGetCmd)
}
