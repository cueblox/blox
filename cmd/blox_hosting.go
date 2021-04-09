package cmd

import (
	"errors"

	"github.com/cueblox/blox/internal/hosting"
	// register hosting providers
	_ "github.com/cueblox/blox/internal/hosting/azure"
	// register hosting providers
	_ "github.com/cueblox/blox/internal/hosting/netlify"
	// register hosting providers
	_ "github.com/cueblox/blox/internal/hosting/vercel"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var (
	provider string
)

// hostingCmd represents the hosting command
var hostingCmd = &cobra.Command{
	Use:   "hosting",
	Short: "Generate the necessary boiler plate to host content",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}
var cmdList = &cobra.Command{
	Use:   "list",
	Short: "List available providers",
	Long:  `List available hosting providers.`,

	Run: func(cmd *cobra.Command, args []string) {
		list := hosting.Providers()
		for _, p := range list {
			pterm.Info.Printf("%s:\t %s\n", p.Name(), p.Description())
		}
	},
}
var cmdInstall = &cobra.Command{
	Use:   "install",
	Short: "Install hosting support for a provider",
	Long:  `Install hosting support for a provider`,

	Run: func(cmd *cobra.Command, args []string) {
		p := hosting.GetProvider(provider)
		if p == nil {
			err := errors.New("unknown provider")
			cobra.CheckErr(err)
		}
		pterm.Info.Printf("Installing support for %s\n", p.Name())
		p.Install()
	},
}

func init() {
	hostingCmd.AddCommand(cmdList)
	hostingCmd.AddCommand(cmdInstall)
	rootCmd.AddCommand(hostingCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// hostingCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// hostingCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	cmdInstall.Flags().StringVarP(&provider, "provider", "p", "azure", "hosting provider to target")
	cobra.CheckErr(cmdInstall.MarkFlagRequired("provider"))
}
