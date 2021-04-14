package cmd

import (
	// import go:embed
	_ "embed"
	"io/ioutil"
	"os"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var (
	//go:embed blox.cue
	bloxcue     string
	dataDir     string
	buildDir    string
	schemataDir string
	skipConfig  bool
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create folders and configuration to maintain content with the blox toolset",
	Long: `Create a group of folders to store your content. A directory for your data,
schemata, and build output will be created.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := createDirectories()
		cobra.CheckErr(err)
		pterm.Info.Println("Initialized folder structures.")

	},
}

func createDirectories() error {
	pterm.Debug.Printf("Creating directory for data at '%s'\n", dataDir)
	err := os.MkdirAll(dataDir, 0755)
	if err != nil {
		return err
	}

	pterm.Debug.Printf("Creating directory for schemata at '%s'\n", schemataDir)
	err = os.MkdirAll(schemataDir, 0755)
	if err != nil {
		return err
	}

	pterm.Debug.Printf("Creating directory for build output at '%s'\n", buildDir)
	err = os.MkdirAll(buildDir, 0755)
	if err != nil {
		return err
	}

	pterm.Debug.Println("Creating 'blox.cue' configuration file")
	return ioutil.WriteFile("blox.cue", []byte(bloxcue), 0644)
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringVarP(&dataDir, "data", "d", "data", "where pre-processed content will be stored (source markdown or yaml)")
	initCmd.Flags().StringVarP(&buildDir, "build", "b", "_build", "where post-processed content will be stored (output json)")
	initCmd.Flags().StringVarP(&schemataDir, "schemata", "s", "schemata", "where the schemata will be stored")
	initCmd.Flags().BoolVarP(&skipConfig, "skip", "c", false, "don't write a configuration file")

}
