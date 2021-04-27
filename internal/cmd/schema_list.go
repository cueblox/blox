package cmd

import (
	"github.com/cueblox/blox/internal/repository"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type schemaListCmd struct {
	cmd *cobra.Command
}

func newSchemaListCmd() *schemaListCmd {
	root := &schemaListCmd{}
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List Schemata",
		Long:  `List schemata that are published via this repository`,
		Args:  cobra.NoArgs,
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
			_ = pterm.DefaultTable.WithHasHeader().WithData(td).Render()
		},
	}

	root.cmd = cmd
	return root
}
