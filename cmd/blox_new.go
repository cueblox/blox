package cmd

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"cuelang.org/go/encoding/yaml"
	"github.com/cueblox/blox"
	"github.com/cueblox/blox/internal/cuedb"
	"github.com/cueblox/blox/internal/cueutils"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var (
	dataSetName string
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new content file for the target dataset",
	Long: `This command will allow you to create new content based on the
template attributes within the schemata. By providing a dataset name and ID(slug)
for the new content, you can quickly scaffold new content with ease.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		userConfig, err := ioutil.ReadFile("blox.cue")
		cobra.CheckErr(err)

		engine, err := cuedb.NewEngine()
		cobra.CheckErr(err)

		cfg, err := blox.NewConfig(BaseConfig)
		cobra.CheckErr(err)

		err = cfg.LoadConfigString(string(userConfig))
		cobra.CheckErr(err)

		// Load Schemas!
		schemataDir, err := cfg.GetString("schemata_dir")
		pterm.Debug.Printf("\t\tSchemata Directory: %s\n", schemataDir)
		cobra.CheckErr(err)

		err = filepath.WalkDir(schemataDir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() {
				bb, err := os.ReadFile(path)
				if err != nil {
					return err
				}
				pterm.Debug.Printf("\t\tLoading Schema: %s\n", path)

				err = engine.RegisterSchema(string(bb))
				if err != nil {
					return err
				}
			}
			return nil
		})
		cobra.CheckErr(err)

		dataSet, err := engine.GetDataSet(dataSetName)
		if err != nil {
			pterm.Error.Printf("Couldn't find dataset '%s'\n", dataSetName)
			pterm.Info.Println("The following DataSets are available:")
			dataSets := engine.GetDataSets()
			for _, dataSet := range dataSets {
				pterm.Info.Printf("\t%s\n", strings.TrimPrefix(dataSet.ID(), "#"))
			}
			return
		}

		templateInstance, err := engine.CueRuntime.Compile("", "")
		cobra.CheckErr(err)

		dsp := dataSet.GetDefinitionPath()
		dsv := engine.Runtime.Database.LookupPath(dsp)

		templateValue, err := cueutils.CreateFromTemplate(templateInstance.Value(), dsv)
		cobra.CheckErr(err)
		templateValue = templateValue.LookupPath(dataSet.GetDefinitionPath())

		dataSetDirectory := fmt.Sprintf("%s/%s", cfg.GetStringOr("data_dir", ""), dataSet.GetDataDirectory())

		slug := args[0]
		pterm.Info.Printf("Creating new %s at %s/%s.yaml\n", dataSet.ID(), dataSetDirectory, slug)

		err = os.MkdirAll(dataSetDirectory, 0755)
		cobra.CheckErr(err)

		bytes, err := yaml.Encode(templateValue)
		cobra.CheckErr(err)

		err = ioutil.WriteFile(fmt.Sprintf("%s/%s.yaml", dataSetDirectory, slug), bytes, 0644)
		cobra.CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(newCmd)

	newCmd.Flags().StringVar(&dataSetName, "dataset", "", "Which DataSet to create content for?")
	cobra.CheckErr(newCmd.MarkFlagRequired("dataset"))
}
