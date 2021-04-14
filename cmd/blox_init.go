package cmd

import (
	// import go:embed
	_ "embed"
	"errors"
	"os"

	"github.com/cueblox/blox/internal/hosting"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var (
	sourceDir  string
	buildDir   string
	staticDir  string
	schemaDir  string
	skipConfig bool
	extension  string
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create folders and configuration to maintain content with the blox toolset",
	Long: `Create a group of folders to store your content. 

If provided, the folders will be created under the "base" directory. 
If "base" is set to an empty string, the source, destination, and template
folders will be created in the root of the current directory.

The "source" directory will store your un-processed content, 
typically Markdown files.

The "destination" directory is where the blox tools will put 
content after it has been validated and processed.

The "template" directory is where you can store templates for
each content type with pre-filled values.
`,
	Run: func(cmd *cobra.Command, args []string) {
		err := createDirectories()
		cobra.CheckErr(err)
		pterm.Info.Println("Initialized folder structures.")

	},
}

func createDirectories() error {
	err := os.MkdirAll(sourceDir, 0755)
	if err != nil {
		return errors.New("creating source directory")
	}

	err = os.MkdirAll(buildDir, 0755)
	if err != nil {
		return errors.New("creating destination directory")
	}

	err = os.MkdirAll(staticDir, 0755)
	if err != nil {
		return errors.New("creating static directory")
	}
	err = os.MkdirAll(schemaDir, 0755)
	if err != nil {
		return errors.New("creating schema directory")
	}
	hosting.CreateFileWithContents("blox.cue", bloxcue)
	return nil
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringVarP(&sourceDir, "source", "s", "content", "where pre-processed content will be stored (source markdown)")
	initCmd.Flags().StringVarP(&buildDir, "destination", "d", "_build", "where post-processed content will be stored (output json)")
	initCmd.Flags().StringVarP(&staticDir, "static", "k", "static", "where static files will be stored")
	initCmd.Flags().StringVarP(&schemaDir, "schema", "a", "schema", "where schema definitions will be stored")
	initCmd.Flags().StringVarP(&extension, "extension", "e", ".md", "default file extension for new content")
	initCmd.Flags().BoolVarP(&skipConfig, "skip", "c", false, "don't write a configuration file")

}

//go:embed blox.cue
var bloxcue string
