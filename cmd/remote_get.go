package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/cueblox/blox/internal/cuedb"
	"github.com/cueblox/blox/internal/repository"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var remoteGetCmd = &cobra.Command{
	Use:   "get <repository> <schema name> <version>",
	Short: "Add a remote schema to your repository",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(3),

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
						fmt.Println(v.Name, v.Schema, v.Definition)
					}
				}
			}
		}
		engine, err := cuedb.NewEngine()
		cobra.CheckErr(err)

		schemaDir := engine.Config.GetStringOr("schema_dir", "schema")

		err = os.MkdirAll(schemaDir, 0755)
		cobra.CheckErr(err)
		filename := fmt.Sprintf("%s_%s.cue", schemaName, version)
		filePath := path.Join(schemaDir, filename)
		err = os.WriteFile(filePath, []byte(selectedVersion.Definition), 0755)
		cobra.CheckErr(err)

	},
}

func init() {
	remoteCmd.AddCommand(remoteGetCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
