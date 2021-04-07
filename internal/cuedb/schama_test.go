package cuedb

import (
	"testing"

	"cuelang.org/go/cue"
)

func TestGetSchemaV1Metadata(t *testing.T) {
	// Can we get the version from the schema's metadata?
	schemaV1Metadata := `{
		_schema: {
			namespace: "cueblox.com"
			name: "profile"
		}
}`

	var cueRuntime cue.Runtime
	cueInstance, err := cueRuntime.Compile("schemaWithVersion", schemaV1Metadata)

	metadata, err := GetSchemaV1Metadata(cueInstance.Value())
	if nil != err {
		t.FailNow()
	}

	if "cueblox.com" != metadata.Namespace {
		t.FailNow()
	}

	if "profile" != metadata.Name {
		t.FailNow()
	}

	invalidSchemaV1Metadata := `{
	_schema: {
		namespace: 123
	}
}`

	cueInstance, err = cueRuntime.Compile("invalidSchema", invalidSchemaV1Metadata)

	_, err = GetSchemaV1Metadata(cueInstance.Value())
	if nil == err {
		t.FailNow()
	}
}

func TestGetV1Model(t *testing.T) {
	// Can we get the version from the schema's metadata?
	cueWithModel := `{
		_model: {
			plural: "iAmPlurals"
		}
}`

	var cueRuntime cue.Runtime
	cueInstance, err := cueRuntime.Compile("", cueWithModel)

	plural, err := GetV1Model(cueInstance.Value())
	if nil != err {
		t.FailNow()
	}

	if plural != "iAmPlurals" {
		t.FailNow()
	}

	// Do we error if there is no eversion?
	cueWithoutModel := `{
		model: {
			plural: "nope"
		}
}`

	cueInstance, err = cueRuntime.Compile("", cueWithoutModel)

	_, err = GetV1Model(cueInstance.Value())
	if nil == err {
		t.FailNow()
	}
}
