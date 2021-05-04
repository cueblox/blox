package cmd

import (
	"github.com/cueblox/blox/internal/sync"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	_ "github.com/cueblox/blox/internal/sync/faunadb"
)

type syncCmd struct {
	cmd *cobra.Command
}

func newSyncCmd() *syncCmd {
	root := &syncCmd{}
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Synchronize blox dataset",
		Long: `The sync command allows you to synchronize your blox dataset with a
remote datastore.`,
		Run: func(cmd *cobra.Command, args []string) {
			pterm.Info.Println("Synchronization started")
			cobra.CheckErr(synchronizeDataset())
		},
	}
	cmd.AddCommand(
		newSyncProvidersCmd().cmd,
	)
	root.cmd = cmd
	return root
}

func synchronizeDataset() error {
	engine, err := sync.Open("faunadb")
	cobra.CheckErr(err)
	err = engine.Synchronize()
	if err != nil {
		pterm.Error.Println(engine.Help())
		return err
	}
	return nil
}
