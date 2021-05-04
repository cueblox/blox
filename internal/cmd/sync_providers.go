package cmd

import (
	"github.com/cueblox/blox/internal/sync"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	_ "github.com/cueblox/blox/internal/sync/faunadb"
)

type syncProvidersCmd struct {
	cmd *cobra.Command
}

func newSyncProvidersCmd() *syncProvidersCmd {
	root := &syncProvidersCmd{}
	cmd := &cobra.Command{
		Use:   "providers",
		Short: "List available synchronization providers",
		Long:  `List registered synchronization providers`,
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
	for _, s := range sync.Providers() {
		line := []string{s}
		td = append(td, line)

	}
	_ = pterm.DefaultTable.WithHasHeader().WithData(td).Render()
	return nil
}
