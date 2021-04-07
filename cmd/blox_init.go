package cmd

import (
	"errors"
	"os"
	"path"

	"github.com/devrel-blox/blox/blox"
	"github.com/devrel-blox/blox/config"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var (
	base        string
	source      string
	destination string
	static      string
	templates   string
	skipConfig  bool
	extension   string
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create folders and configuration to maintain content with the drb toolset",
	Long: `Create a group of folders to store your content. 

If provided, the folders will be created under the "base" directory. 
If "base" is set to an empty string, the source, destination, and template
folders will be created in the root of the current directory.

The "source" directory will store your un-processed content, 
typically Markdown files.

The "destination" directory is where the drb tools will put 
content after it has been validated and processed.

The "template" directory is where you can store templates for
each content type with pre-filled values.
`,
	Run: func(cmd *cobra.Command, args []string) {
		root, err := os.Getwd()
		if err != nil {
			cmd.PrintErr("unable to get current directory")
			return
		}
		err = createDirectories(root)
		cobra.CheckErr(err)
		if !skipConfig {
			err = writeConfigFile()
			cobra.CheckErr(err)
		}
		for _, model := range blox.Models {

			model, err := blox.GetModel(model.ID)
			cobra.CheckErr(err)
			cfg, err := config.Load()
			cobra.CheckErr(err)
			cobra.CheckErr(model.New(model.ID+cfg.DefaultExtension, model.TemplatePath()))
		}
		pterm.Info.Println("Initialized folder structures.")

	},
}

func writeConfigFile() error {
	cfg := config.BloxConfig{
		Base:             base,
		Source:           source,
		Templates:        templates,
		Destination:      destination,
		Static:           static,
		DefaultExtension: extension,
	}
	f, err := os.Create("blox.yaml")
	if err != nil {
		return err
	}
	defer f.Close()
	err = cfg.Write(f)
	return err
}

func createDirectories(root string) error {
	err := os.MkdirAll(sourceDir(root), 0755)
	if err != nil {
		return errors.New("creating source directory")
	}
	err = os.MkdirAll(destinationDir(root), 0755)
	if err != nil {
		return errors.New("creating destination directory")
	}
	err = os.MkdirAll(templateDir(root), 0755)
	if err != nil {
		return errors.New("creating template directory")
	}

	err = os.MkdirAll(staticDir(root), 0755)
	if err != nil {
		return errors.New("creating dir directory")
	}
	return nil
}

func sourceDir(root string) string {
	return path.Join(root, base, source)
}
func destinationDir(root string) string {
	return path.Join(root, base, destination)
}
func templateDir(root string) string {
	return path.Join(root, base, templates)
}

func staticDir(root string) string {
	return path.Join(root, base, static)
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringVarP(&base, "base", "b", "content", "base directory for pre- and post- processed content")
	initCmd.Flags().StringVarP(&source, "source", "s", "source", "where pre-processed content will be stored (source markdown)")
	initCmd.Flags().StringVarP(&destination, "destination", "d", "out", "where post-processed content will be stored (output json)")
	initCmd.Flags().StringVarP(&static, "static", "k", "static", "where static files will be stored")
	initCmd.Flags().StringVarP(&templates, "template", "t", "templates", "where content templates will be stored")
	initCmd.Flags().StringVarP(&extension, "extension", "e", ".md", "default file extension for new content")
	initCmd.Flags().BoolVarP(&skipConfig, "skip", "c", false, "don't write a configuration file")

}
