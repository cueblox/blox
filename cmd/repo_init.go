package cmd

import (
	"github.com/cueblox/blox/internal/repository"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var root string
var namespace string
var output string

// initCmd represents the init command
var repoInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a New Schema Repository",
	Long: `Initializing a new schema repository creates the
configuration required to published your schemata.`,
	Run: func(cmd *cobra.Command, args []string) {
		repo, err := repository.NewRepository(namespace, output, root)
		cobra.CheckErr(err)
		pterm.Success.Printf("Created repository for %s\n", repo.Namespace)
	},
}

func init() {
	repoCmd.AddCommand(repoInitCmd)
	repoCmd.PersistentFlags().StringVarP(&root, "root", "r", "repository", "directory to store the repository, relative to current directory")
	repoCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "schemas.you.com", "repository namespace")
	repoCmd.PersistentFlags().StringVarP(&output, "output", "o", "_build", "directory where build output will be written")

}
