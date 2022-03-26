package gonfigenvironmentvariables

import (
	"os"
	"testing"
	"time"

	"github.com/andreasisnes/goconut"
	"github.com/stretchr/testify/assert"
)

func newConfig(option *EnvironmentVariablesOptions) goconut.IConfiguration {
	return goconut.NewBuilder().
		Add(NewEnvironmentVariablesSource(option)).
		Build()
}

func TestGet(t *testing.T) {
	key := "TestGet"
	expected := "TEST_VALUE"
	os.Setenv(key, expected)
	config := newConfig(nil)

	assert.Equal(t, expected, config.Get(key, nil))
}

func TestGetAfterAddedValue(t *testing.T) {
	config := newConfig(&EnvironmentVariablesOptions{
		RefreshInterval: time.Second,
		SourceOptions: goconut.SourceOptions{
			ReloadOnChange: true,
		},
	})

	key := "TestGetAfterAddedValue"
	expected := "TEST_VALUE"
	os.Setenv(key, expected)
	time.Sleep(time.Second * 3)
	result := config.Get(key, nil)

	assert.Equal(t, expected, result)
	config.Deconstruct()
}
