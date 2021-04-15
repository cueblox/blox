package cmd

import (
	// import go:embed
	_ "embed"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var (
	//go:embed blox.cue
	bloxcue     string
	dataDir     string
	buildDir    string
	schemataDir string
	starter     string
	skipConfig  bool
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create folders and configuration to maintain content with the blox toolset",
	Long: `Create a group of folders to store your content. A directory for your data,
schemata, and build output will be created.`,
	Run: func(cmd *cobra.Command, args []string) {

		if starter != "" {
			cobra.CheckErr(installStarter(starter))
			pterm.Info.Println("Starter initialized.")
			return
		}
		// not a starter
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

	initCmd.Flags().StringVarP(&starter, "starter", "t", "", "use a pre-defined starter in the CURRENT directory")

}

func installStarter(starter string) error {

	pterm.Info.Printf("Installing starter %s\n", starter)
	// kinda hacky, look for things that make a url or domain name
	// if it's internal, it'll just be one word
	internal := !strings.ContainsAny(starter, "/,.")
	var repo string
	if internal {
		repo = fmt.Sprintf("https://github.com/cueblox/starter-%s", starter)
	} else {
		repo = starter
	}

	// git init in the existing directory
	cmd := exec.Command("git", "init")
	pterm.Info.Println("git init...")
	err := cmd.Run()
	if err != nil {
		pterm.Info.Printf("git init error: %s\n", err)
		return err
	}

	// add the starter as a remote
	cmd = exec.Command("git", "remote", "add", "origin", repo)
	pterm.Info.Println("git remote add...")
	err = cmd.Run()
	if err != nil {
		pterm.Info.Printf("git remote error: %s\n", err)
		return err
	}
	// git fetch on the remote
	cmd = exec.Command("git", "fetch")
	pterm.Info.Println("git fetch...")
	err = cmd.Run()
	if err != nil {
		pterm.Info.Printf("git fetch error: %s\n", err)
		return err
	}

	// checkout main
	cmd = exec.Command("git", "checkout", "origin/main", "-ft")
	pterm.Info.Println("git checkout...")
	err = cmd.Run()
	if err != nil {
		pterm.Info.Printf("git checkout error: %s\n", err)
		return err
	}
	here, err := os.Getwd()
	if err != nil {
		return err
	}

	// now remove all traces of the git checkout
	err = os.RemoveAll(path.Join(here, ".git"))

	if err != nil {
		return err
	}
	pterm.Info.Println(repo)
	return nil
}
