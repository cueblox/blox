package content

import (
	"encoding/json"
	"os"
	"path"
	"path/filepath"

	"github.com/pterm/pterm"
)

func (s *Service) RenderJSON() ([]byte, error) {
	if !s.built {
		err := s.build()
		if err != nil {
			return nil, err
		}
	}
	pterm.Debug.Println("Building output data blox")
	output, err := s.engine.GetOutput()
	if err != nil {
		return nil, err
	}

	pterm.Debug.Println("Rendering data blox to JSON")
	return output.MarshalJSON()
}

func (s *Service) RenderAndSave() error {
	if !s.built {
		err := s.build()
		if err != nil {
			return err
		}
	}

	bb, err := s.RenderJSON()
	if err != nil {
		return err
	}
	buildDir, err := s.Cfg.GetString("build_dir")
	if err != nil {
		return err
	}
	err = os.MkdirAll(buildDir, 0o755)
	if err != nil {
		return err
	}

	build_cue, err := s.Cfg.GetBool("output_cue")
	if err != nil {
		return err
	}
	// only output the cue if it's specified in the config file
	if build_cue {
		cval, err := s.engine.MarshalString()
		if err != nil {
			pterm.Error.Println("error getting cue output")
			return err
		}
		cval = "package blox\n" + cval

		err = os.WriteFile(filepath.Join(buildDir, "data.cue"), []byte(cval), 0o755)
		if err != nil {
			return err
		}
		pterm.Success.Printf("Cue output written to '%s'\n", filepath.Join(buildDir, "data.cue"))

	}
	// always output the full json dataset
	filename := "data.json"
	filePath := path.Join(buildDir, filename)
	err = os.WriteFile(filePath, bb, 0o755)
	if err != nil {
		return err
	}
	pterm.Success.Printf("Data blox written to '%s'\n", filePath)

	build_recordsets, err := s.Cfg.GetBool("output_recordsets")
	if err != nil {
		return err
	}
	if build_recordsets {

		var dataList map[string][]map[string]interface{}

		err = json.Unmarshal(bb, &dataList)
		if err != nil {
			return err
		}

		for k := range dataList {
			set := dataList[k]
			ss, err := json.Marshal(set)
			if err != nil {
				return err
			}
			filename := k + ".json"
			filePath := path.Join(buildDir, filename)

			// write the array
			err = os.WriteFile(filePath, ss, 0o755)
			if err != nil {
				return err
			}
			dirpath := path.Join(buildDir, k)
			err = os.MkdirAll(dirpath, 0o755)
			if err != nil {
				if err != os.ErrExist {
					return err
				}
			}
			for j := range set {
				slug := set[j]["id"].(string)
				// write each item
				filename := slug + ".json"
				filePath := path.Join(dirpath, filename)
				derp := path.Dir(filePath)
				err = os.MkdirAll(derp, 0o755)
				if err != nil {
					if err != os.ErrExist {
						return err
					}
				}
				ss, err := json.Marshal(set[j])
				if err != nil {
					return err
				}
				err = os.WriteFile(filePath, ss, 0o755)
				if err != nil {
					return err
				}
			}

		}
		pterm.Success.Printf("Recordset output written to '%s'\n", buildDir)

	}

	pterm.Info.Println("Running Postbuild Plugins")
	err = s.runPostPlugins()
	if err != nil {
		return err
	}
	return nil
}
