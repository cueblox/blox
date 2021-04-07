package cmd

import (
	"path"

	"github.com/devrel-blox/blox/blox"
	"github.com/devrel-blox/blox/config"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var (
	model string
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new content file",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		model, err := blox.GetModel(model)
		cobra.CheckErr(err)

		pterm.Info.Printf("Creating new %s in %s\n", model.Name, model.SourceContentPath())

		slug := args[0]
		cfg, err := config.Load()
		cobra.CheckErr(err)
		cobra.CheckErr(model.New(slug+cfg.DefaultExtension, model.SourceContentPath()))
		pterm.Info.Printf("Your new content file is ready at %s\n", path.Join(model.SourceFilePath(slug)))
	},
}

func init() {
	rootCmd.AddCommand(newCmd)

	newCmd.Flags().StringVarP(&model, "type", "t", "article", "type of content to create")
	cobra.CheckErr(newCmd.MarkFlagRequired("type"))
	newCmd.SetUsageTemplate("drb new --type [type name] [slug]")
}
