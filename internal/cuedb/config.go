package cuedb

import (
	"fmt"
	"io/ioutil"

	"cuelang.org/go/cue"
)

const defaultConfigName = "blox.cue"

func (r *Runtime) LoadConfig() error {
	cueConfig, err := ioutil.ReadFile(defaultConfigName)
	if err != nil {
		return err
	}

	return r.loadConfigString(string(cueConfig))
}

func (r *Runtime) loadConfigString(cueConfig string) error {
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

func (r *Runtime) GetString(key string) (string, error) {
	keyValue := r.database.LookupPath(cue.ParsePath(key))

	if keyValue.Exists() {
		return keyValue.String()
	}

	return "", fmt.Errorf("Couldn't find key '%s'", key)
}

func (r *Runtime) GetStringOr(key string, def string) string {
	cueValue, err := r.GetString(key)

	if err != nil {
		return def
	}

	return cueValue
}
