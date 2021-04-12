package cuedb

import (
	"fmt"
	"path"
	"strings"
	"testing"

	"cuelang.org/go/cue"
	"github.com/cueblox/blox"
	"github.com/stretchr/testify/assert"
)

func testConfig(t *testing.T) *blox.Config {
	c, err := blox.NewConfig(baseConfig)
	if err != nil {
		t.Error(err)
	}
	err = c.LoadConfig(path.Join("..", "..", "blox.cue"))
	if err != nil {
		t.Error()
	}
	return c
}
func TestNewRuntime(t *testing.T) {
	runtime, err := NewEngineWithConfig(testConfig(t))

	if err != nil {
		fmt.Println(err)
		t.Fatal("Failed to create a NewRuntime, which should never happen")
	}

	// Should be created with no datasets
	assert.Equal(t, 0, runtime.CountDataSets())
}

func TestExtractSchemaMetadata(t *testing.T) {
	schemaCue := `{
	_schema: {
		namespace: "TestNS"
		name: "TestSchema"
	}
}`
	runtime, err := NewEngineWithConfig(testConfig(t))

	if err != nil {
		t.Fatal("Failed to create a NewRuntime, which should never happen")
	}

	cueInstance, err := runtime.CueRuntime.Compile("", schemaCue)
	assert.Equal(t, nil, err)

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
	runtime, err := NewEngineWithConfig(testConfig(t))

	if err != nil {
		t.Fatal("Failed to create a NewRuntime, which should never happen")
	}

	cueInstance, err := runtime.CueRuntime.Compile("", dataSetCue)
	assert.Equal(t, nil, err)

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
	runtime, err := NewEngineWithConfig(testConfig(t))

	if err != nil {
		t.Fatal("Failed to create a NewRuntime, which should never happen")
	}

	err = runtime.RegisterSchema(schemaCue)
	if err != nil {
		t.Fatal(err)
	}

	assert.Len(t, runtime.dataSets, 2)
}

func TestGetDataSets(t *testing.T) {
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
	runtime, err := NewEngineWithConfig(testConfig(t))

	if err != nil {
		t.Fatal("Failed to create a NewRuntime, which should never happen")
	}

	err = runtime.RegisterSchema(schemaCue)
	if err != nil {
		t.Fatal(err)
	}

	dataSets := runtime.GetDataSets()
	assert.Len(t, dataSets, 2)

	_, err = runtime.GetDataSet("Nope")
	assert.NotEqual(t, nil, err)

	// Currently we/Cue is leaking "#" prefix on DataSet names
	// I would like to remove this, but for now we'll work with
	// either
	dataSet, err := runtime.GetDataSet("#One")
	assert.Equal(t, nil, err)
	assert.Equal(t, "#One", dataSet.name)

	dataSet, err = runtime.GetDataSet("One")
	assert.Equal(t, "#One", dataSet.name)

	dataSet, err = runtime.GetDataSet("Two")
	assert.Equal(t, "#Two", dataSet.name)
}

func TestInsert(t *testing.T) {
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
	runtime, err := NewEngineWithConfig(testConfig(t))

	if err != nil {
		t.Fatal("Failed to create a NewRuntime, which should never happen")
	}

	err = runtime.RegisterSchema(schemaCue)
	if err != nil {
		t.Fatal(err)
	}

	// Get DataSet so we can insert
	dataSet, err := runtime.GetDataSet("#One")
	assert.Equal(t, "#One", dataSet.name)

	recordOne := map[string]interface{}{"david": map[string]string{"name": "David"}}
	assert.Equal(t, nil, runtime.Insert(dataSet, recordOne))

	nameValueOne := runtime.Database.LookupPath(cue.ParsePath("data.OnePlural.david.name"))
	nameValueOneStr, err := nameValueOne.String()
	assert.Equal(t, nil, err)
	assert.Equal(t, "David", nameValueOneStr)

	recordTwo := map[string]interface{}{"brian": map[string]string{"name": "Brian"}}
	assert.Equal(t, nil, runtime.Insert(dataSet, recordTwo))

	nameValueTwo := runtime.Database.LookupPath(cue.ParsePath("data.OnePlural.brian.name"))
	nameValueTwoStr, err := nameValueTwo.String()
	assert.Equal(t, nil, err)
	assert.Equal(t, "Brian", nameValueTwoStr)
}

func TestDataSetID(t *testing.T) {
	dataSet := DataSet{
		name: "MyDataSet",
	}

	assert.Equal(t, dataSet.ID(), "mydataset")

	dataSet = DataSet{
		name: "randomname",
	}
	assert.Equal(t, dataSet.ID(), "randomname")
}

func TestCueDataPath(t *testing.T) {
	dataSet := DataSet{
		metadata: DataSetMetadata{
			Plural: "testPlural",
		},
	}

	assert.Equal(t, dataSet.CueDataPath(), cue.ParsePath(fmt.Sprintf("%s.%s", dataPathRoot, dataSet.metadata.Plural)))
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

	assert.Equal(t, strings.ReplaceAll(dataSet.GetDataMapCue(), " ", ""), strings.ReplaceAll(`{
		some: data: path: MyDataSet: _
		data: MyDataSets: [ID=string]: some.data.path.MyDataSet & {id: (ID)}
	}`, " ", ""))
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
