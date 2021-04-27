package cmd

import (
	"github.com/spf13/cobra"
)

type remoteCmd struct {
	cmd *cobra.Command
}

func newRemoteCmd() *remoteCmd {
	root := &remoteCmd{}
	cmd := &cobra.Command{
		Use:   "remote",
		Short: "Manage Schemata",
		Long: `Blox allows you to consume schemata from remote repositories.
The remote subcommands allow you to list the available schemata from these
repositories, as well as download a schema to your local directories.`,
	}
	cmd.AddCommand(
		newRemoteGetCmd().cmd,
		newRemoteListCmd().cmd,
	)
	root.cmd = cmd
	return root
}
