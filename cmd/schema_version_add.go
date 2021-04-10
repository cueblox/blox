package cmd

import (
	"github.com/cueblox/blox/internal/repository"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add [schema name]",
	Short: "Add a new version to a schema",
	Args:  cobra.ExactArgs(1),

	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		repo, err := repository.GetRepository()
		cobra.CheckErr(err)
		pterm.Info.Printf("Using repository: %s\n", repo.Namespace)
		schema := args[0]
		pterm.Info.Printf("Creating schema: %s\n", schema)
		err = repo.AddVersion(schema)
		cobra.CheckErr(err)
		pterm.Success.Printf("Schema %s created\n", schema)

	},
}

func init() {
	versionCmd.AddCommand(addCmd)
}
