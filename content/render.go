package content

import (
	"os"
	"path"

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
	filename := "data.json"
	filePath := path.Join(buildDir, filename)
	err = os.WriteFile(filePath, bb, 0o755)
	if err != nil {
		return err
	}
	pterm.Success.Printf("Data blox written to '%s'\n", filePath)
	return nil
}
