package blox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const base = `{
    build_dir:    string | *"_build"
    data_dir:     string | *"data"
    schemata_dir: string | *"schemata"
	static_dir:   string | *"static"
}`

func TestGetString(t *testing.T) {
	runtime, err := NewConfig(base)

	if err != nil {
		t.Fatal("Failed to create a NewRuntime, which should never happen")
	}

	err = runtime.LoadConfigString(`{
		data_dir: "my-data-dir"
}`)
	assert.Equal(t, nil, err)

	// test non-existent key
	_, err = runtime.GetString("data_dir_non_exist")
	assert.NotEqual(t, nil, err)

	// test defaulted key
	schDir, err := runtime.GetString("schemata_dir")
	assert.Equal(t, nil, err)
	assert.Equal(t, "schemata", schDir)

	// test parsed key
	configString, err := runtime.GetString("data_dir")
	assert.Equal(t, nil, err)
	assert.Equal(t, "my-data-dir", configString)
}
