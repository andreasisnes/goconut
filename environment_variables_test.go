package gonfigenvironmentvariables

import (
	"os"
	"testing"
	"time"

	"github.com/andreasisnes/goconut"
	"github.com/stretchr/testify/assert"
)

func newBuilder(option *EnvironmentVariablesOptions) goconut.IConfiguration {
	return goconut.NewBuilder().
		Add(NewEnvironmentVariablesSource(option)).
		Build()
}

func TestField(t *testing.T) {
	expected := "TEST_VALUE"
	os.Setenv("TMP_TEST_VALUE", expected)
	b := newBuilder(nil)

	assert.Equal(t, expected, b.Get("TMP_TEST_VALUE", nil))
}

func TestFieldWithRefresh(t *testing.T) {
	b := newBuilder(&EnvironmentVariablesOptions{
		RefreshInterval: time.Second,
		SourceOptions: goconut.SourceOptions{
			ReloadOnChange: true,
		},
	})

	key := "TestFieldWithRefresh"
	expected := "TEST_VALUE"
	os.Setenv(key, expected)
	time.Sleep(time.Second * 5)
	result := b.Get(key, nil)
	b.Deconstruct()

	assert.Equal(t, expected, result)
}
