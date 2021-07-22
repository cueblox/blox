package cmd

import (
	"fmt"
	"io/ioutil"

	"github.com/cueblox/blox/repository"
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

			/*
				remotes, err := cfg.GetList("remotes")
				if err == nil {
					cobra.CheckErr(parseRemotes(remotes))
				}
				if images {
					err = processImages(cfg)
					if err != nil {
						cobra.CheckErr(err)
					}
				}
			*/

			repo, err := repository.NewService(string(userConfig))

			cobra.CheckErr(err)
			err = repo.Build(referentialIntegrity)
			cobra.CheckErr(err)
			bb, err := repo.RenderJSON()
			cobra.CheckErr(err)
			fmt.Println(string(bb))
			/*buildDir, err := cfg.GetString("build_dir")
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

			funcMap := tpl.FuncMap{
				// The name "title" is what the function will be called in the template text.
				"rfcdate": func(t string) string {
					tm, err := time.Parse("2006-01-02 15:04", t)
					if err != nil {
						return err.Error()
					}
					val := tm.Format(time.RFC1123)
					return val
				},
				"now": func() string {
					tm := time.Now()
					val := tm.Format(time.RFC1123)
					return val
				},
			}

			err = json.Unmarshal(jso, &dataJson)
			cobra.CheckErr(err)

			t := tpl.Must(tpl.New(template).Funcs(funcMap).ParseFiles(paths...))
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
			*/
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
