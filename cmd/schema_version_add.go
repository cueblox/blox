package cmd

import (
	"github.com/cueblox/blox/internal/repository"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add [schema name]",
	Short: "Create a New Version of a Schema",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo, err := repository.GetRepository()
		cobra.CheckErr(err)

		pterm.Info.Printf("Using repository: %s\n", repo.Namespace)
		schema := args[0]

		pterm.Info.Printf("Creating schema: %s\n", schema)
		cobra.CheckErr(repo.AddVersion(schema))
		pterm.Success.Printf("Schema %s created\n", schema)
	},
}

func init() {
	versionCmd.AddCommand(addCmd)
}
