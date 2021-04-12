package cmd

import (
	"github.com/cueblox/blox/internal/repository"
	"github.com/pterm/pterm"
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
		pterm.Info.Println("Building Repository")
		cobra.CheckErr(repo.Build())
		pterm.Success.Println("Build Complete")
	},
}

func init() {
	repoCmd.AddCommand(repoBuildCmd)
}
