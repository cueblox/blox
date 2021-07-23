package cmd

import (
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/cueblox/blox/content"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
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

			repo, err := content.NewService(string(userConfig), referentialIntegrity)
			cobra.CheckErr(err)

			if static {
				staticDir, err := repo.Cfg.GetString("static_dir")
				pterm.Info.Printf("Serving static files from %s\n", staticDir)

				cobra.CheckErr(err)
				http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(filepath.Join(".", staticDir)))))
			}

			hf, err := repo.GQLHandlerFunc()
			cobra.CheckErr(err)
			http.HandleFunc("/", hf)

			h, err := repo.GQLPlaygroundHandler()
			cobra.CheckErr(err)
			http.Handle("/ui", h)

			pterm.Info.Printf("Server is running at %s\n", address)
			cobra.CheckErr(http.ListenAndServe(address, nil))
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
