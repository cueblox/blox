package cmd

import (
	"github.com/cueblox/blox/internal/repository"
	"github.com/spf13/cobra"
)

// repoBuildCmd represents the build command
var repoBuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the schema repository",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		repo, err := repository.GetRepository()
		cobra.CheckErr(err)
		cobra.CheckErr(repo.Build())
	},
}

func init() {
	repoCmd.AddCommand(repoBuildCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// buildCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// buildCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
