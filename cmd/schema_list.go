package cmd

import (
	"github.com/cueblox/blox/internal/repository"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List Schemata",
	Long:  `List schemata that are published via this repository`,
	Run: func(cmd *cobra.Command, args []string) {
		repo, err := repository.GetRepository()
		cobra.CheckErr(err)

		var td pterm.TableData
		header := []string{"Namespace", "Schema", "Version"}
		td = append(td, header)
		for _, s := range repo.Schemas {
			for _, v := range s.Versions {
				line := []string{repo.Namespace, s.Name, v.Name}
				td = append(td, line)
			}
		}
		pterm.DefaultTable.WithHasHeader().WithData(td).Render()
	},
}

func init() {
	schemaCmd.AddCommand(listCmd)
}
