package cmd

import (
	"github.com/cueblox/blox/schema"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var name string

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
		repo, err := schema.GetRepository()
		cobra.CheckErr(err)
		pterm.Info.Printf("Using repository: %s\n", repo.Namespace)
		schema := args[0]
		pterm.Info.Printf("Adding schema: %s\n", schema)
		err = repo.AddSchema(schema)
		cobra.CheckErr(err)
		pterm.Info.Printf("Schema %s created\n", schema)
	},
}

func init() {
	schemaCmd.AddCommand(schemaNewCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
