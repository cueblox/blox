package config

import (
	_ "embed"
)

//go:embed config.cue
var BaseConfig string
