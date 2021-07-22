package cmd

import (
	"fmt"
	"io/ioutil"

	"github.com/cueblox/blox/repository"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	// Import the blob packages we want to be able to open.
	_ "gocloud.dev/blob/azureblob"
	_ "gocloud.dev/blob/gcsblob"
	_ "gocloud.dev/blob/s3blob"
)

type bloxServeCmd struct {
	cmd *cobra.Command
}

func newBloxServeCmd() *bloxServeCmd {
	root := &bloxServeCmd{}
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Serve a GraphQL API",

		Run: func(cmd *cobra.Command, args []string) {
			userConfig, err := ioutil.ReadFile("blox.cue")

			pterm.Debug.Printf("loading user config")

			cobra.CheckErr(err)

			/*
				remotes, err := cfg.GetList("remotes")
				if err == nil {
					cobra.CheckErr(parseRemotes(remotes))
				}
				if images {
					err = processImages(cfg)
					if err != nil {
						cobra.CheckErr(err)
					}
				}
			*/

			repo, err := repository.NewService(string(userConfig), referentialIntegrity)
			cobra.CheckErr(err)

			bb, err := repo.RenderJSON()
			cobra.CheckErr(err)
			fmt.Println(string(bb))

		},
	}
	cmd.Flags().BoolVarP(&static, "static", "s", true, "Serve static files")
	cmd.Flags().StringVarP(&address, "address", "a", ":8080", "Listen address")

	root.cmd = cmd
	return root
}

var (
	static  bool
	address string
)
