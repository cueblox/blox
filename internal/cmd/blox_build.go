package cmd

import (
	"io/ioutil"

	"github.com/cueblox/blox/content"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type bloxBuildCmd struct {
	cmd *cobra.Command
}

func newBloxBuildCmd() *bloxBuildCmd {
	root := &bloxBuildCmd{}
	cmd := &cobra.Command{
		Use:   "build",
		Short: "Validate & Build dataset",
		Long: `The build command will ensure that your dataset is correct by
validating it against your schemata. Once validated, it will render all
your content into a single JSON file, which can be consumed by your tooling
of choice.
	
Referential Integrity can be enforced with -i. This ensures that any fields
ending with _id are valid references to identifiers within the other content type.
	`,
		Run: func(cmd *cobra.Command, args []string) {
			userConfig, err := ioutil.ReadFile("blox.cue")

			pterm.Debug.Printf("loading user config")

			cobra.CheckErr(err)

			repo, err := content.NewService(string(userConfig), referentialIntegrity)
			cobra.CheckErr(err)

			err = repo.RenderAndSave()
			cobra.CheckErr(err)
		},
	}
	cmd.Flags().BoolVarP(&referentialIntegrity, "referential-integrity", "i", false, "Verify referential integrity")
	root.cmd = cmd
	return root
}

var (
	referentialIntegrity bool
)

const DefaultConfigName = "blox.cue"

const BaseConfig = `{
    #Remote: {
        name: string
        version: string
        repository: string
    }
    build_dir:    string | *"_build"
    data_dir:     string | *"data"
    schemata_dir: string | *"schemata"
	static_dir: string | *"static"
	template_dir: string | *"templates"
	remotes: [ ...#Remote ]

}`
