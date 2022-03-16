package gonfigenvironmentvariables

import (
	"os"
	"strings"
	"sync"

	"github.com/andreasisnes/goconut"
)

type EnvironmentVariablesOptions struct {
	goconut.SourceOptions
	Delimiter string
}

type EnvironmentVariablesSource struct {
	EnvOptions    EnvironmentVariablesOptions
	WaitGroup     sync.WaitGroup
	Configuration map[string]interface{}
}

func NewEnvironmentVariablesSource(options *EnvironmentVariablesOptions) goconut.ISource {
	return &EnvironmentVariablesSource{
		EnvOptions:    *options,
		WaitGroup:     sync.WaitGroup{},
		Configuration: make(map[string]interface{}),
	}
}

func (e *EnvironmentVariablesSource) Options() goconut.SourceOptions {
	return e.EnvOptions.SourceOptions
}

func (e *EnvironmentVariablesSource) Exists(key string) bool {
	return false
}

func (e *EnvironmentVariablesSource) Get(key string) interface{} {
	return nil
}

func (e *EnvironmentVariablesSource) GetKeys() []string {
	return nil
}

func (e *EnvironmentVariablesSource) IsDirty() bool {
	return false
}

func (e *EnvironmentVariablesSource) Load() {
	for _, variable := range os.Environ() {
		keyIdx := strings.Index(variable, "=")
		e.Configuration[variable[:keyIdx]] = variable[keyIdx+1:]
	}
}

func (e *EnvironmentVariablesSource) Deconstruct() {

}

func (e *EnvironmentVariablesSource) watcher() {
	e.WaitGroup.Add(1)
	defer e.WaitGroup.Done()
}
