package cmd

import (
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Schema Version Management",
}

func init() {
	schemaCmd.AddCommand(versionCmd)
}
