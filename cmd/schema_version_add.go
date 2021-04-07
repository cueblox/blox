package cmd

import (
	"github.com/cueblox/blox/schema"
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
		repo, err := schema.GetRepository()
		cobra.CheckErr(err)
		pterm.Info.Printf("Using repository: %s\n", repo.Namespace)
		schema := args[0]
		pterm.Info.Printf("using schema: %s\n", schema)
		err = repo.AddVersion(schema)
		cobra.CheckErr(err)
	},
}

func init() {
	versionCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
