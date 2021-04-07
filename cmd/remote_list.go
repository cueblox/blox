package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cueblox/blox/schema"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var remoteListCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		manifest := fmt.Sprintf("https://%s/manifest.json", args[0])
		res, err := http.Get(manifest)
		cobra.CheckErr(err)
		var repos schema.Repository
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
