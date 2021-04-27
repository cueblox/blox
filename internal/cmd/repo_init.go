package cmd

import (
	"github.com/cueblox/blox/internal/repository"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var (
	reporoot  string
	namespace string
	output    string
)

type repoInitCmd struct {
	cmd *cobra.Command
}

func newRepoInitCmd() *repoInitCmd {
	root := &repoInitCmd{}
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a New Schema Repository",
		Long: `Initializing a new schema repository creates the
	configuration required to published your schemata.`,
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			repo, err := repository.NewRepository(namespace, output, reporoot)
			cobra.CheckErr(err)
			pterm.Success.Printf("Created repository for %s\n", repo.Namespace)
		},
	}

	cmd.PersistentFlags().StringVarP(&reporoot, "root", "r", "repository", "directory to store the repository, relative to current directory")
	cmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "schemas.you.com", "repository namespace")
	cmd.PersistentFlags().StringVarP(&output, "output", "o", "_build", "directory where build output will be written")
	root.cmd = cmd
	return root
}
