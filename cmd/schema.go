package cmd

import (
	"github.com/spf13/cobra"
)

var schemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "Create, Manage, and Version your Schemata",
}

func init() {
	rootCmd.AddCommand(schemaCmd)
}
