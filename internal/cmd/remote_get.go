package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/cueblox/blox"
	"github.com/cueblox/blox/internal/repository"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type remoteGetCmd struct {
	cmd *cobra.Command
}

func newRemoteGetCmd() *remoteGetCmd {
	root := &remoteGetCmd{}
	cmd := &cobra.Command{
		Use:   "get <repository> <schema name> <version>",
		Short: "Add a remote schema to your repository",
		Args:  cobra.ExactArgs(3),

		Run: func(cmd *cobra.Command, args []string) {
			repo := args[0]
			schemaName := args[1]
			version := args[2]

			err := ensureRemote(schemaName, version, repo)
			cobra.CheckErr(err)
		},
	}

	root.cmd = cmd
	return root
}

// ensureRemote downloads a remote schema at a specific
// version if it doesn't exist locally
// TODO: Duplicated in content package. Consolidate
func ensureRemote(name, version, repo string) error {
	// load config
	userConfig, err := ioutil.ReadFile("blox.cue")
	cobra.CheckErr(err)

	cfg, err := blox.NewConfig(BaseConfig)
	cobra.CheckErr(err)

	err = cfg.LoadConfigString(string(userConfig))
	cobra.CheckErr(err)

	// Load Schemas!
	schemataDir, err := cfg.GetString("schemata_dir")
	cobra.CheckErr(err)

	schemaFilePath := path.Join(schemataDir, fmt.Sprintf("%s_%s.cue", name, version))
	_, err = os.Stat(schemaFilePath)
	if os.IsNotExist(err) {
		pterm.Info.Printf("Schema does not exist locally: %s_%s.cue\n", name, version)
		manifest := fmt.Sprintf("https://%s/manifest.json", repo)
		res, err := http.Get(manifest)
		cobra.CheckErr(err)

		var repos repository.Repository
		err = json.NewDecoder(res.Body).Decode(&repos)
		cobra.CheckErr(err)

		var selectedVersion *repository.Version
		for _, s := range repos.Schemas {
			if s.Name == name {
				for _, v := range s.Versions {
					if v.Name == version {
						selectedVersion = v
						pterm.Debug.Println(v.Name, v.Schema, v.Definition)
					}
				}
			}
		}

		// make schemata directory
		cobra.CheckErr(os.MkdirAll(schemataDir, 0o755))

		// TODO: don't overwrite each time
		filename := fmt.Sprintf("%s_%s.cue", name, version)
		filePath := path.Join(schemataDir, filename)
		err = os.WriteFile(filePath, []byte(selectedVersion.Definition), 0o755)
		cobra.CheckErr(err)
		pterm.Info.Printf("Schema downloaded: %s_%s.cue\n", name, version)
		return nil
	}
	pterm.Info.Println("Schema already exists locally, skipping download.")
	return nil
}
