package cmd

import (
	"github.com/spf13/cobra"
)

type schemaVersionCmd struct {
	cmd *cobra.Command
}

func newSchemaVersionCmd() *schemaVersionCmd {
	root := &schemaVersionCmd{}
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Schema Version Management",
	}
	cmd.AddCommand(
		newSchemaVersionAddCmd().cmd,
	)
	root.cmd = cmd
	return root
}
