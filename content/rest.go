package content

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/pterm/pterm"
)

// RestHandlerFunc returns a handler function that will render
// the dataset specified as the last path parameter.
func (s *Service) RestHandlerFunc() (http.HandlerFunc, error) {
	if !s.built {
		err := s.build()
		if err != nil {
			return nil, err
		}
	}

	hf := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET")

		path := strings.TrimPrefix(r.URL.Path, "/")
		parts := strings.Split(path, "/")
		pterm.Debug.Println(parts, len(parts))
		if len(parts) == 0 {
			pterm.Warning.Println("No dataset specified")
			w.WriteHeader(http.StatusNotFound)
			return
		}
		dataset := parts[len(parts)-1]

		ds, err := s.engine.GetDataSetByPlural(dataset)
		if err != nil {
			pterm.Warning.Println("Requested dataset not found", parts, len(parts))
			w.WriteHeader(http.StatusNotFound)
			return
		}

		data := s.engine.GetAllData(ds.GetExternalName())
		err = json.NewEncoder(w).Encode(data)
		if err != nil {
			pterm.Warning.Printf("failed to encode: %v", err)
		}
	}
	return hf, nil
}
