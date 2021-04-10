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
	Short: "Initialize a new schema repository",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		repo, err := repository.NewRepository(namespace, output, root)
		cobra.CheckErr(err)
		pterm.Info.Printf("Created repository for %s\n", repo.Namespace)
	},
}

func init() {
	repoCmd.AddCommand(repoInitCmd)
	repoCmd.PersistentFlags().StringVarP(&root, "root", "r", "repository", "directory to store the repository")

	repoCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "schemas.you.com", "repository namespace")

	repoCmd.PersistentFlags().StringVarP(&output, "output", "o", "_build", "directory where build output will be written")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
