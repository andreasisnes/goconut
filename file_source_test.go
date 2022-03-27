package file

import (
	"encoding/json"
	"os"
	"path"
	"testing"
	"time"

	"github.com/andreasisnes/goconut"
	"github.com/stretchr/testify/assert"
)

var (
	dataDir  = "data"
	notAFile = path.Join(dataDir, "notafile.yaml")
	json1    = path.Join(dataDir, "config1.json")
	toml1    = path.Join(dataDir, "config1.toml")
	yaml1    = path.Join(dataDir, "config1.yaml")
	json2    = path.Join(dataDir, "config2.json")
	toml2    = path.Join(dataDir, "config2.toml")
	yaml2    = path.Join(dataDir, "config2.yaml")
	json3    = path.Join(dataDir, "config3.json")
)

func TestJsonObjectField(t *testing.T) {
	config := goconut.NewBuilder().
		Add(NewFileSource(&Options{Path: json1})).
		Build()

	res := config.Get("SimpleField", nil)
	assert.Equal(t, "<SimpleField-1>", res)
}

func TestJsonObjectFieldLayered(t *testing.T) {
	config := goconut.NewBuilder().
		Add(NewFileSource(&Options{Path: json1})).
		Add(NewFileSource(&Options{Path: json2})).
		Build()

	res := config.Get("SimpleField", nil)
	assert.Equal(t, "<SimpleField-2>", res)
}

func TestTomlObject(t *testing.T) {
	config := goconut.NewBuilder().
		Add(NewFileSource(&Options{Path: toml1})).
		Build()

	res := config.Get("SimpleField", nil)
	assert.Equal(t, "<SimpleField-1>", res)
}

func TestYamlObject(t *testing.T) {
	config := goconut.NewBuilder().
		Add(NewFileSource(&Options{Path: yaml1})).
		Build()

	res := config.Get("SimpleField", nil)
	assert.Equal(t, "<SimpleField-1>", res)
}

func TestUnkownFile(t *testing.T) {
	config := goconut.NewBuilder().
		Add(NewFileSource(&Options{Path: notAFile, SourceOptions: goconut.SourceOptions{Optional: true, ReloadOnChange: true}})).
		Build()

	res := config.Get("SimpleField", nil)
	assert.Nil(t, res)
}

func TestUnkownFileAsNotOptional(t *testing.T) {
	defer func() {
		assert.NotNil(t, recover())
	}()
	goconut.NewBuilder().
		Add(NewFileSource(&Options{Path: notAFile, SourceOptions: goconut.SourceOptions{Optional: false}})).
		Build()
}

func TestJsonReloadOnChange(t *testing.T) {
	t.Parallel()
	key := "ChangedField"
	originalValue := "UnchangedValue"
	alteredValue := "AlteredValue"
	config := goconut.NewBuilder().
		Add(NewFileSource(&Options{Path: json3, SourceOptions: goconut.SourceOptions{ReloadOnChange: true}})).
		Build()
	assert.Equal(t, config.Get(key, nil), originalValue)

	data := readFile(json3)
	data[key] = alteredValue
	writeFile(json3, data)
	defer func() {
		data[key] = originalValue
		writeFile(json3, data)
	}()

	time.Sleep(time.Second * 5)

	assert.Equal(t, config.Get(key, nil), alteredValue)
}

func writeFile(file string, data map[string]interface{}) {
	content, _ := json.Marshal(data)
	os.WriteFile(json3, content, os.ModePerm)
}

func readFile(file string) map[string]interface{} {
	content, _ := os.ReadFile(json3)
	result := make(map[string]interface{})
	json.Unmarshal(content, &result)
	return result
}
