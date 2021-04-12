package blox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const base = `{
    data_dir: string
    schema_dir: string | *"schemas"
    build_dir: string | *"_build"
}
`

func TestGetString(t *testing.T) {
	runtime, err := New(base)

	if err != nil {
		t.Fatal("Failed to create a NewRuntime, which should never happen")
	}

	err = runtime.loadConfigString(`{
		data_dir: "my-data-dir"
}`)
	assert.Equal(t, nil, err)

	_, err = runtime.GetString("data_dir_non_exist")
	assert.NotEqual(t, nil, err)

	configString, err := runtime.GetString("data_dir")
	assert.Equal(t, nil, err)
	assert.Equal(t, "my-data-dir", configString)
}
