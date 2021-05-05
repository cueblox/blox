package cmd

import (
	"fmt"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var (
	quiet bool
	debug bool
)

func Execute(version string, exit func(int), args []string) {
	newRootCmd(version, exit).Execute(args)
}

type rootCmd struct {
	cmd  *cobra.Command
	exit func(int)
}

func (cmd *rootCmd) Execute(args []string) {
	cmd.cmd.SetArgs(args)

	if err := cmd.cmd.Execute(); err != nil {
		fmt.Println(err.Error())
		cmd.exit(1)
	}
}

func newRootCmd(version string, exit func(int)) *rootCmd {
	root := &rootCmd{
		exit: exit,
	}
	cmd := &cobra.Command{
		Use:           "blox",
		Short:         "CueBlox is a suite of slightly opinionated tools for managing and sharing content repositories of YAML and Markdown documents.",
		Version:       version,
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cobra.NoArgs,
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
	cmd.PersistentFlags().BoolVar(&quiet, "quiet", false, "disable logging")
	cmd.PersistentFlags().BoolVar(&debug, "debug", false, "enable debug logging, overrides 'quiet' flag")

	cmd.AddCommand(
		newRemoteCmd().cmd,
		newSchemaCmd().cmd,
		newRepoCmd().cmd,
		newCompletionCmd().cmd,
		newDocsCmd().cmd,
		newBloxBuildCmd().cmd,
		newBloxInitCmd().cmd,
		newBloxNewCmd().cmd,
		newBloxRenderCmd().cmd,
		newExportCmd().cmd,
	)

	root.cmd = cmd
	return root
}
