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

func TestGetWithReloadOnChange(t *testing.T) {
	t.Parallel()
	config := newConfig(&EnvironmentVariablesOptions{
		RefreshInterval: time.Second,
		SourceOptions: goconut.SourceOptions{
			ReloadOnChange: true,
		},
	})

	key := "TestGetWithReloadOnChange"
	expected := "TEST_VALUE"
	os.Setenv(key, expected)
	time.Sleep(time.Second * 2)
	result := config.Get(key, nil)

	assert.Equal(t, expected, result)
	config.Deconstruct()
}

func TestGetWithSentinel(t *testing.T) {
	t.Parallel()
	config := newConfig(&EnvironmentVariablesOptions{
		RefreshInterval: time.Second,
		SourceOptions: goconut.SourceOptions{
			ReloadOnChange: false,
			SentinelOptions: &goconut.SentinelOptions{
				Key:           "TestGetWithSentinel",
				RefreshPolicy: goconut.RefreshCurrent,
			},
		},
	})

	key := "TestGetWithSentinel"
	expected := "TEST_VALUE"
	os.Setenv(key, expected)
	time.Sleep(time.Second * 2)
	result := config.Get(key, nil)

	assert.Equal(t, expected, result)
	config.Deconstruct()
}

func TestGetNilWithSentinel(t *testing.T) {
	t.Parallel()
	config := newConfig(&EnvironmentVariablesOptions{
		RefreshInterval: time.Second,
		SourceOptions: goconut.SourceOptions{
			ReloadOnChange: false,
			SentinelOptions: &goconut.SentinelOptions{
				Key:           "TestGetNilWithSentinelUnkownKey",
				RefreshPolicy: goconut.RefreshCurrent,
			},
		},
	})

	key := "TestGetNilWithSentinel"
	notExpected := "TEST_VALUE"
	os.Setenv(key, notExpected)
	time.Sleep(time.Second * 2)
	result := config.Get(key, nil)

	assert.Nil(t, result)
	config.Deconstruct()
}
