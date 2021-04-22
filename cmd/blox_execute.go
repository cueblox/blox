/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
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

var template string
var with string
var each bool

// executeCmd represents the execute command
var executeCmd = &cobra.Command{
	Use:   "execute",
	Short: "Execute templates with compiled data",
	Long: `Execute templates with compiled data. 
Use the 'with' parameter to restrict the data set to a single content type.
Use the 'each' parameter to execute the template once for each item.`,
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println("execute called")
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
			parseRemotes(remotes)
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
		cobra.CheckErr(os.MkdirAll(buildDir, 0755))

		filename := "data.json"
		filePath := path.Join(buildDir, filename)
		cobra.CheckErr(os.WriteFile(filePath, jso, 0755))
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

func init() {
	rootCmd.AddCommand(executeCmd)
	executeCmd.Flags().StringVarP(&template, "template", "t", "", "template to execute")
	executeCmd.Flags().StringVarP(&with, "with", "w", "", "dataset to use")
	executeCmd.Flags().BoolVarP(&each, "each", "e", false, "execute template once per item")
}
