package gonfigenvironmentvariables

import (
	"os"
	"strings"
	"sync"

	"github.com/andreasisnes/goconut"
)

type EnvironmentVariablesOptions struct {
	ReloadOnChange bool
	
}

type EnvironmentVariablesProvider struct {
	WaitGroup     sync.WaitGroup
	configuration map[string]interface{}
}

func NewEnvironmentVariablesProvider(options *EnvironmentVariablesOptions) goconut.ISource {
	return &EnvironmentVariablesProvider{
		WaitGroup:     sync.WaitGroup{},
		configuration: make(map[string]interface{}),
	}
}

func (e *EnvironmentVariablesProvider) Exists(key string) bool {
	return false
}

func (e *EnvironmentVariablesProvider) Get(key string) interface{} {
	return nil
}

func (e *EnvironmentVariablesProvider) GetKeys() []string {
	return nil
}

func (e *EnvironmentVariablesProvider) IsDirty() bool {
	return false
}

func (e *EnvironmentVariablesProvider) Load() {
	for _, variable := range os.Environ() {
		keyIdx := strings.Index(variable, "=")
		e.configuration[variable[:keyIdx]] = variable[keyIdx+1:]
	}
}

func (e *EnvironmentVariablesProvider) Deconstruct() {

}

func (e *EnvironmentVariablesProvider) watcher() {
	e.WaitGroup.Add(1)
	defer e.WaitGroup.Done()
}
