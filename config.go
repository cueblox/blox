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
	cueValue := r.runtime.CueContext.CompileString(cueConfig)

	if cueValue.Err() != nil {
		return cueValue.Err()
	}

	r.runtime.Database = r.runtime.Database.Unify(cueValue)
	if err := r.runtime.Database.Validate(); err != nil {
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

func (r *Config) GetBool(key string) (bool, error) {
	keyValue := r.runtime.Database.LookupPath(cue.ParsePath(key))

	if keyValue.Exists() {
		return keyValue.Bool()
	}

	return false, fmt.Errorf("couldn't find key '%s'", key)
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

const BaseConfig = `{
	#Remote: {
	    name: string
	    version: string
	    repository: string
	}
	#Plugin: {
		name: string
		executable: string
	}
	build_dir:    string | *"_build"
	data_dir:     string | *"data"
	schemata_dir: string | *"schemata"
	static_dir: string | *"static"
	template_dir: string | *"templates"
	output_cue: bool | *false
	output_recordsets: bool | *false
	remotes: [ ...#Remote ]
	prebuild: [...#Plugin]
	postbuild: [...#Plugin]
}`

const DefaultConfigName = "blox.cue"
