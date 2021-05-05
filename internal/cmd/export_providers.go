package cmd

import (
	"github.com/cueblox/blox/internal/export"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	_ "github.com/cueblox/blox/internal/export/faunadb"
)

type exportProvidersCmd struct {
	cmd *cobra.Command
}

func newExportProvidersCmd() *exportProvidersCmd {
	root := &exportProvidersCmd{}
	cmd := &cobra.Command{
		Use:   "providers",
		Short: "List available export providers",
		Long:  `List registered export providers`,
		Run: func(cmd *cobra.Command, args []string) {
			cobra.CheckErr(listProviders())
		},
	}

	root.cmd = cmd
	return root
}

func listProviders() error {
	var td pterm.TableData
	header := []string{"Name"}
	td = append(td, header)
	for _, s := range export.Providers() {
		line := []string{s}
		td = append(td, line)

	}
	_ = pterm.DefaultTable.WithHasHeader().WithData(td).Render()
	return nil
}
