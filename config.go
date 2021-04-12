package blox

import (
	"fmt"
	"io/ioutil"

	"cuelang.org/go/cue"
)

type Config struct {
	*Runtime
}

// New setup a new config type with
// base as the defaults
func NewConfig(base string) (*Config, error) {
	r, err := NewRuntimeWithBase(base)
	if err != nil {
		return nil, err
	}
	config := &Config{
		Runtime: r,
	}

	return config, nil
}

// LoadConfig opens the configuration file
// specified in `path` and validates it against
// the configuration provided in when the `Engine`
// was initialized with `New()`
func (r *Config) LoadConfig(path string) error {
	cueConfig, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	return r.LoadConfigString(string(cueConfig))
}

func (r *Config) LoadConfigString(cueConfig string) error {
	cueInstance, err := r.CueRuntime.Compile("", cueConfig)
	if err != nil {
		return err
	}

	cueValue := cueInstance.Value()

	r.Database = r.Database.Unify(cueValue)
	if err = r.Database.Validate(); err != nil {
		return err
	}

	return nil
}
func (r *Config) GetString(key string) (string, error) {
	keyValue := r.Database.LookupPath(cue.ParsePath(key))

	if keyValue.Exists() {
		return keyValue.String()
	}

	return "", fmt.Errorf("couldn't find key '%s'", key)
}

func (r *Config) GetStringOr(key string, def string) string {
	cueValue, err := r.GetString(key)

	if err != nil {
		return def
	}

	return cueValue
}
