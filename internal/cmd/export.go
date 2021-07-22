package cmd

import (
	"github.com/spf13/cobra"
)

func newExportCmd() *exportCmd {
	root := &exportCmd{}
	cmd := &cobra.Command{
		Use:   "export",
		Short: "ALPHA: Export Schema", SilenceUsage: true,
		Args: cobra.NoArgs,
	}

	root.cmd = cmd
	cmd.AddCommand(
		newExportSchemaCmd().cmd,
	)

	return root
}

type exportCmd struct {
	cmd *cobra.Command
}
