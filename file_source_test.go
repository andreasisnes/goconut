package file

import (
	"path"
	"testing"

	"github.com/andreasisnes/goconut"
	"github.com/stretchr/testify/assert"
)

var (
	dataDir = "data"
	json1   = path.Join(dataDir, "jsonconfig1.json")
	json2   = path.Join(dataDir, "jsonconfig2.json")
)

func TestJsonFile(t *testing.T) {
	builder := goconut.NewBuilder().
		Add(NewFileProvider(json1, false, true)).
		Add(NewFileProvider(json1, false, true)).
		Build()

	var g string
	builder.Get("SimpleField", &g)
	assert.Equal(t, "<SimpleField>", g)
}
