package cuedb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetString(t *testing.T) {
	runtime, err := NewRuntime()

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
