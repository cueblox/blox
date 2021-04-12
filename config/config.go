package config

import (
	"fmt"
	"io/ioutil"

	"cuelang.org/go/cue"
)

type Engine struct {
	cueRuntime *cue.Runtime
	database   cue.Value
}

// New setup a new config engine with
// base as the default
func New(base string) (Engine, error) {
	var cueRuntime cue.Runtime
	cueInstance, err := cueRuntime.Compile("", base)

	if nil != err {
		return Engine{}, err
	}

	engine := Engine{
		cueRuntime: &cueRuntime,
		database:   cueInstance.Value(),
	}

	return engine, nil
}

// LoadConfig opens the configuration file
// specified in `path` and validates it against
// the configuration provided in when the `Engine`
// was initialized with `New()`
func (r *Engine) LoadConfig(path string) error {
	cueConfig, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	return r.loadConfigString(string(cueConfig))
}

func (r *Engine) loadConfigString(cueConfig string) error {
	cueInstance, err := r.cueRuntime.Compile("", cueConfig)
	if err != nil {
		return err
	}

	cueValue := cueInstance.Value()

	r.database = r.database.Unify(cueValue)
	if err = r.database.Validate(); err != nil {
		return err
	}

	return nil
}

func (r *Engine) GetString(key string) (string, error) {
	keyValue := r.database.LookupPath(cue.ParsePath(key))

	if keyValue.Exists() {
		return keyValue.String()
	}

	return "", fmt.Errorf("couldn't find key '%s'", key)
}

func (r *Engine) GetStringOr(key string, def string) string {
	cueValue, err := r.GetString(key)

	if err != nil {
		return def
	}

	return cueValue
}
