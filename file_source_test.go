package file

import (
	"path"
	"testing"

	"github.com/andreasisnes/goconut"
	"github.com/stretchr/testify/assert"
)

var (
	dataDir = "data"
	yaml1   = path.Join(dataDir, "config1.yaml")
	toml1   = path.Join(dataDir, "config1.toml")
	json1   = path.Join(dataDir, "jsonconfig1.json")
	json2   = path.Join(dataDir, "jsonconfig2.json")
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
