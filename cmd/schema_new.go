package cmd

import (
	"github.com/cueblox/blox/internal/repository"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// schemaNewCmd represents the new command
var schemaNewCmd = &cobra.Command{
	Use:   "new [schema name]",
	Short: "Add a new schema to a repository",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo, err := repository.GetRepository()
		cobra.CheckErr(err)
		pterm.Info.Printf("Using repository: %s\n", repo.Namespace)
		schema := args[0]
		pterm.Info.Printf("Adding schema: %s\n", schema)
		err = repo.AddSchema(schema)
		cobra.CheckErr(err)
		pterm.Success.Printf("Schema %s created\n", schema)
	},
}

func init() {
	schemaCmd.AddCommand(schemaNewCmd)
}
