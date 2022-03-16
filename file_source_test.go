package file

import (
	"path"
	"testing"

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
)

func TestJsonObjectField(t *testing.T) {
	config := goconut.NewBuilder().
		Add(NewFileProvider(json1, false, true)).
		Build()

	res := config.Get("SimpleField", nil)
	assert.Equal(t, "<SimpleField-1>", res)
}

func TestJsonObjectFieldLayered(t *testing.T) {
	config := goconut.NewBuilder().
		Add(NewFileProvider(json1, false, true)).
		Add(NewFileProvider(json2, false, true)).
		Build()

	res := config.Get("SimpleField", nil)
	assert.Equal(t, "<SimpleField-2>", res)
}

func TestTomlObject(t *testing.T) {
	config := goconut.NewBuilder().
		Add(NewFileProvider(toml1, false, true)).
		Build()

	res := config.Get("SimpleField", nil)
	assert.Equal(t, "<SimpleField-1>", res)
}

func TestYamlObject(t *testing.T) {
	config := goconut.NewBuilder().
		Add(NewFileProvider(yaml1, false, true)).
		Build()

	res := config.Get("SimpleField", nil)
	assert.Equal(t, "<SimpleField-1>", res)
}

func TestUnkownFile(t *testing.T) {
	config := goconut.NewBuilder().
		Add(NewFileProvider(notAFile, false, true)).
		Build()

	res := config.Get("SimpleField", nil)
	assert.Nil(t, res)
}

func TestUnkownFileAsNotOptional(t *testing.T) {
	defer func() {
		assert.NotNil(t, recover())
	}()
	goconut.NewBuilder().
		Add(NewFileProvider(notAFile, true, true)).
		Build()
}
