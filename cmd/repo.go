package cmd

import (
	"github.com/spf13/cobra"
)

// repoCmd represents the repo command
var repoCmd = &cobra.Command{
	Use:   "repo",
	Short: "Create & Manage Schema Repositories",
}

func init() {
	rootCmd.AddCommand(repoCmd)
}
