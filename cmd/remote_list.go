package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cueblox/blox/internal/repository"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var remoteListCmd = &cobra.Command{
	Use:   "list <remote repository URL>",
	Short: "List schemas and versions available in a remote repository",
	Args:  cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		manifest := fmt.Sprintf("https://%s/manifest.json", args[0])
		res, err := http.Get(manifest)
		cobra.CheckErr(err)

		var repos repository.Repository
		json.NewDecoder(res.Body).Decode(&repos)

		// TODO extract and reuse with schema_list.go
		var td pterm.TableData
		header := []string{"Namespace", "Schema", "Version"}
		td = append(td, header)
		for _, s := range repos.Schemas {
			for _, v := range s.Versions {
				line := []string{repos.Namespace, s.Name, v.Name}
				td = append(td, line)
			}
		}

		pterm.DefaultTable.WithHasHeader().WithData(td).Render()
	},
}

func init() {
	remoteCmd.AddCommand(remoteListCmd)
}
