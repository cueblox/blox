package cuedb

import (
	"testing"

	"cuelang.org/go/cue"
	"github.com/stretchr/testify/assert"
)

func TestNewRuntime(t *testing.T) {
	runtime, err := NewRuntime()

	if err != nil {
		t.Fatal("Failed to create a NewRuntime, which should never happen")
	}

	// Should be created with no datasets
	assert.Equal(t, runtime.CountDataSets(), 0)
}

func TestExtractSchemaMetadata(t *testing.T) {
	schemaCue := `{
	_schema: {
		namespace: "TestNS"
		name: "TestSchema"
	}
}`
	runtime, err := NewRuntime()

	if err != nil {
		t.Fatal("Failed to create a NewRuntime, which should never happen")
	}

	cueInstance, err := runtime.cueRuntime.Compile("", schemaCue)

	schemaMetdata, err := runtime.extractSchemaMetadata(cueInstance.Value())
	if err != nil {
		t.Fatal("Failed to extract SchemaMetadata")
	}

	assert.Equal(t, "TestNS", schemaMetdata.Namespace)
	assert.Equal(t, "TestSchema", schemaMetdata.Name)
}

func TestExtractDataSetMetadata(t *testing.T) {
	dataSetCue := `{
	_dataset: {
		plural: "DataSetPlural"
		supportedExtensions: ["txt", "tar.gz"]
	}
}`
	runtime, err := NewRuntime()

	if err != nil {
		t.Fatal("Failed to create a NewRuntime, which should never happen")
	}

	cueInstance, err := runtime.cueRuntime.Compile("", dataSetCue)

	dataSetMetdata, err := runtime.extractDataSetMetadata(cueInstance.Value())
	if err != nil {
		t.Fatal("Failed to extract DataSetMetadata")
	}

	assert.Equal(t, "DataSetPlural", dataSetMetdata.Plural)
	assert.Equal(t, []string{"txt", "tar.gz"}, dataSetMetdata.SupportedExtensions)
}

func TestRegisterSchema(t *testing.T) {
	schemaCue := `{
		_schema: {
			namespace: "TestNS1"
			name: "TestSchema1"
		}

		#One: {
			_dataset: {
				plural: "OnePlural"
				supportedExtensions: ["txt", "tar.gz"]
			}
			name: string
		}

		#Two: {
			_dataset: {
				plural: "TwoPlural"
				supportedExtensions: ["txt", "tar.gz"]
			}

			sport: string
		}
}`
	runtime, err := NewRuntime()

	if err != nil {
		t.Fatal("Failed to create a NewRuntime, which should never happen")
	}

	err = runtime.RegisterSchema(schemaCue)
	if err != nil {
		t.Fatal(err)
	}

	assert.Len(t, runtime.dataSets, 2)
}

func TestGetDataMapCue(t *testing.T) {
	somePath := cue.ParsePath("some.data.path")

	dataSet := DataSet{
		name: "MyDataSet",
		metadata: DataSetMetadata{
			Plural: "MyDataSets",
		},
		cuePath: somePath,
	}

	assert.Equal(t, dataSet.GetDataMapCue(), `{
		some: data: path: MyDataSet: _
		data: MyDataSets: [ID=string]: some.data.path.MyDataSet
	}`)
}

func TestGetInlinePath(t *testing.T) {
	somePath := cue.ParsePath("some.data.path")

	dataSet := DataSet{
		cuePath: somePath,
	}

	assert.Equal(t, dataSet.GetInlinePath(), "some: data: path")

	somePath = cue.ParsePath("another.random.path.of.random.length")

	dataSet = DataSet{
		cuePath: somePath,
	}

	assert.Equal(t, dataSet.GetInlinePath(), "another: random: path: of: random: length")
}
