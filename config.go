package blox

import (
	"errors"
	"fmt"
	"io/ioutil"

	"cuelang.org/go/cue"
	"github.com/pterm/pterm"
)

type Config struct {
	runtime *Runtime
}

// New setup a new config type with
// base as the defaults
func NewConfig(base string) (*Config, error) {
	r, err := NewRuntimeWithBase(base)
	if err != nil {
		return nil, err
	}
	config := &Config{
		runtime: r,
	}

	return config, nil
}

// LoadConfig opens the configuration file
// specified in `path` and validates it against
// the configuration provided in when the `Engine`
// was initialized with `New()`
func (r *Config) LoadConfig(path string) error {
	pterm.Debug.Printf("\t\tLoading config: %s\n", path)
	cueConfig, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	return r.LoadConfigString(string(cueConfig))
}

func (r *Config) LoadConfigString(cueConfig string) error {
	cueInstance, err := r.runtime.CueRuntime.Compile("", cueConfig)
	if err != nil {
		return err
	}

	cueValue := cueInstance.Value()

	r.runtime.Database = r.runtime.Database.Unify(cueValue)
	if err = r.runtime.Database.Validate(); err != nil {
		return err
	}

	return nil
}
func (r *Config) GetString(key string) (string, error) {
	keyValue := r.runtime.Database.LookupPath(cue.ParsePath(key))

	if keyValue.Exists() {
		return keyValue.String()
	}

	return "", fmt.Errorf("couldn't find key '%s'", key)
}

func (r *Config) GetList(key string) (cue.Value, error) {

	keyValue := r.runtime.Database.LookupPath(cue.ParsePath(key))

	if keyValue.Exists() {
		_, err := keyValue.List()
		if err != nil {
			return cue.Value{}, err
		}

		return keyValue, nil
	}
	return cue.Value{}, errors.New("not found")
}

func (r *Config) GetStringOr(key string, def string) string {
	cueValue, err := r.GetString(key)

	if err != nil {
		return def
	}

	return cueValue
}
