package cmd

import (
	"github.com/spf13/cobra"
)

func newRepoCmd() *repoCmd {
	root := &repoCmd{}
	cmd := &cobra.Command{
		Use:   "repo",
		Short: "Create & Manage Schema Repositories", SilenceUsage: true,
		Args: cobra.NoArgs,
	}

	root.cmd = cmd
	cmd.AddCommand(
		newRepoInitCmd().cmd,
		newRepoBuildCmd().cmd,
	)

	return root
}

type repoCmd struct {
	cmd *cobra.Command
}
