package cmd

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/arrase/cobraplugins"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string
var quiet bool
var debug bool
var Version string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "blox",
	Version: Version,
	Short:   "CueBlox is a suite of slightly opinionated tools for managing and sharing content repositories of YAML and Markdown documents.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if debug {
			pterm.EnableDebugMessages()
			return
		}

		if quiet {
			pterm.DisableOutput()
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	loadPlugins(rootCmd)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.blox.yaml)")
	rootCmd.PersistentFlags().BoolVar(&quiet, "quiet", false, "disable logging")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "enable debug logging, overrides 'quiet' flag")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".blox" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".blox")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func loadPlugins(cmd *cobra.Command) {
	pterm.Info.Println("Searching for plugins")
	home, err := homedir.Dir()
	cobra.CheckErr(err)
	pluginPath := path.Join(home, ".blox", "plugins")
	err = filepath.Walk(pluginPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {

				pterm.Info.Printf("\tloading %s\n", path)
				cmd.AddCommand(cobraplugins.GetCmd(path, "MainCmd"))
			}

			return nil
		})
	cobra.CheckErr(err)
}
