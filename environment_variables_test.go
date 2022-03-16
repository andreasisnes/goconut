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

func TestCreatingNewField(t *testing.T) {
	b := newBuilder(&EnvironmentVariablesOptions{
		Delimiter:       "__",
		RefreshInterval: time.Second,
		SourceOptions: goconut.SourceOptions{
			ReloadOnChange: false,
		},
	})
	expected := "TEST_VALUE"
	os.Setenv("TMP_TEST_VALUE", expected)
	//time.Sleep(time.Second)

	assert.Nil(t, b.Get("TMP_TEST_VALUE", nil))
}
