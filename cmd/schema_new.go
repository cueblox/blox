package cmd

import (
	"github.com/cueblox/blox/internal/repository"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// schemaNewCmd represents the new command
var schemaNewCmd = &cobra.Command{
	Use:   "new [schema name]",
	Short: "Create a New Schema",
	Long:  `Create a new schema that can be published with the repository management commands`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo, err := repository.GetRepository()
		cobra.CheckErr(err)

		pterm.Info.Printf("Using repository: %s\n", repo.Namespace)
		schema := args[0]

		pterm.Info.Printf("Adding schema: %s\n", schema)
		cobra.CheckErr(repo.AddSchema(schema))
		pterm.Success.Printf("Schema %s created\n", schema)
	},
}

func init() {
	schemaCmd.AddCommand(schemaNewCmd)
}
