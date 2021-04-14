package cmd

import (
	"github.com/cueblox/blox/internal/repository"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// repoBuildCmd represents the build command
var repoBuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build a Schema Repository",
	Long: `In order to consume your schema repository with the Blox CLI, you
need to build a manifest file and publish. This command provides the build output
that can be deployed to any static file hosting, or even GitHub raw content links.`,
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
