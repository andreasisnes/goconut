package gonfigenvironmentvariables

import (
	"os"
	"reflect"
	"strings"
)

type EnvironmentVariablesOptions struct {
}

type EnvironmentVariablesProvider struct {
	flatmap map[string]interface{}
}

func NewEnvironmentVariablesProvider(options *EnvironmentVariablesOptions) gonfig.IProvider {
	return &EnvironmentVariablesProvider{
		flatmap: make(map[string]interface{}),
	}
}

func (e *EnvironmentVariablesProvider) TryGet(key string, value interface{}) bool {
	if val, ok := e.flatmap[key]; ok {
		va := reflect.ValueOf(val)
		reflect.ValueOf(value).Elem().Set(va)
		return true
	}

	return false
}

func (e *EnvironmentVariablesProvider) Get(key string) interface{} {
	return e.flatmap[key]
}

func (e *EnvironmentVariablesProvider) Exists(key string) bool {
	if _, ok := e.flatmap[key]; ok {
		return true
	}

	return false
}

func (e *EnvironmentVariablesProvider) GetKeys() []string {
	keys := make([]string, len(e.flatmap))
	i := 0
	for k := range e.flatmap {
		keys[i] = k
		i++
	}

	return keys
}

func (e *EnvironmentVariablesProvider) Load() {
	for _, variable := range os.Environ() {
		keyIdx := strings.Index(variable, "=")
		e.flatmap[variable[:keyIdx]] = variable[keyIdx+1:]
	}
}
