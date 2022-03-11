package gonfigenvironmentvariables

import (
	"fmt"
	"os"
	"testing"

	"github.com/andreasisnes/gonfig"
	"github.com/stretchr/testify/assert"
)

func TestField(t *testing.T) {
	config := gonfig.NewBuilder().
		Add(NewEnvironmentVariablesProvider(&EnvironmentVariablesOptions{})).
		Add(NewEnvironmentVariablesProvider(&EnvironmentVariablesOptions{})).
		Build()

	expected := "TEST_ENV"
	err := os.Setenv("TMP_TEST_VALUE", expected)
	fmt.Println(err)
	var result string
	err = config.Get("TMP_TEST_VALUE", &result)

	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}
