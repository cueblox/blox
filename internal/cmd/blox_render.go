package cmd

import (
	"encoding/json"
	"errors"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	tpl "text/template"

	"github.com/cueblox/blox"
	"github.com/cueblox/blox/internal/cuedb"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type bloxRenderCmd struct {
	cmd *cobra.Command
}

func newBloxRenderCmd() *bloxRenderCmd {
	root := &bloxRenderCmd{}
	cmd := &cobra.Command{
		Use:   "render",
		Short: "Render templates with compiled data",
		Long: `Render templates with compiled data. 
	Use the 'with' parameter to restrict the data set to a single content type.
	Use the 'each' parameter to execute the template once for each item.`,
		Run: func(cmd *cobra.Command, args []string) {
			// begin cut and paste from blox_build
			userConfig, err := ioutil.ReadFile("blox.cue")

			pterm.Debug.Printf("loading user config")

			cobra.CheckErr(err)

			engine, err := cuedb.NewEngine()

			pterm.Debug.Printf("new engine")
			cobra.CheckErr(err)

			cfg, err := blox.NewConfig(BaseConfig)

			pterm.Debug.Printf("newConfig")
			cobra.CheckErr(err)

			err = cfg.LoadConfigString(string(userConfig))
			cobra.CheckErr(err)

			// Load Schemas!
			schemataDir, err := cfg.GetString("schemata_dir")
			cobra.CheckErr(err)

			remotes, err := cfg.GetList("remotes")
			if err == nil {
				cobra.CheckErr(parseRemotes(remotes))
			}

			pterm.Debug.Printf("\t\tUsing schemata from: %s\n", schemataDir)

			err = filepath.WalkDir(schemataDir, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}

				if !d.IsDir() {
					bb, err := os.ReadFile(path)
					if err != nil {
						return err
					}

					pterm.Debug.Printf("\t\tAttempting to register schema: %s\n", path)
					err = engine.RegisterSchema(string(bb))
					if err != nil {
						return err
					}
				}

				return nil
			})
			cobra.CheckErr(err)

			pterm.Debug.Println("\t\tBuilding DataSets")
			cobra.CheckErr(buildDataSets(engine, cfg))

			if referentialIntegrity {
				pterm.Info.Println("Verifying Referential Integrity")
				cobra.CheckErr(engine.ReferentialIntegrity())
				pterm.Success.Println("Referential Integrity OK")
			}

			pterm.Debug.Println("Building output data blox")
			output, err := engine.GetOutput()
			cobra.CheckErr(err)

			pterm.Debug.Println("Rendering data blox to JSON")
			jso, err := output.MarshalJSON()
			cobra.CheckErr(err)

			buildDir, err := cfg.GetString("build_dir")
			cobra.CheckErr(err)
			cobra.CheckErr(os.MkdirAll(buildDir, 0o755))

			filename := "data.json"
			filePath := path.Join(buildDir, filename)
			cobra.CheckErr(os.WriteFile(filePath, jso, 0o755))
			pterm.Success.Printf("Data blox written to '%s'\n", filePath)

			// end cut and paste from blox_build

			// Load Schemas!
			templateDir, err := cfg.GetString("template_dir")
			cobra.CheckErr(err)
			pterm.Info.Printf("Using templates from %s\n", templateDir)
			if template == "" {
				pterm.Error.Println("template name required")
				return
			}
			// Files are provided as a slice of strings.
			tplPath := path.Join(templateDir, template)
			paths := []string{
				tplPath,
			}
			_, err = os.Stat(tplPath)
			cobra.CheckErr(err)
			var dataJson map[string]interface{}

			err = json.Unmarshal(jso, &dataJson)
			cobra.CheckErr(err)

			t := tpl.Must(tpl.New(template).ParseFiles(paths...))
			if with != "" {
				if each {
					dataset, ok := dataJson[with].([]interface{})
					if !ok {
						err = errors.New("dataset is not a slice")
						cobra.CheckErr(err)
					}
					for _, thing := range dataset {
						err = t.Execute(os.Stdout, thing)
					}
				} else {
					err = t.Execute(os.Stdout, dataJson[with])
				}
			} else {
				err = t.Execute(os.Stdout, dataJson)
			}
			cobra.CheckErr(err)
		},
	}
	cmd.Flags().StringVarP(&template, "template", "t", "", "template to render")
	cmd.Flags().StringVarP(&with, "with", "w", "", "dataset to use")
	cmd.Flags().BoolVarP(&each, "each", "e", false, "render template once per item")
	root.cmd = cmd
	return root
}

var (
	template string
	with     string
	each     bool
)
