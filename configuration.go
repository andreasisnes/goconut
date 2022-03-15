package goconut

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
)

var (
	ErrKeyNotFound      = errors.New("key not found")
	ErrPointerNotPassed = errors.New("pointer not passed")
)

type IConfiguration interface {
	Get(key string, value interface{}) interface{}
	Deconstruct() IConfiguration
	Refresh() bool
	Unmarshal(value interface{}) error
}

type Configuration struct {
	sources   []ISource
	delimiter string
}

func newConfiguration(sources []ISource) IConfiguration {
	config := &Configuration{
		sources:   sources,
		delimiter: ".",
	}

	for _, source := range config.sources {
		source.Load()
	}

	return config
}

func (c *Configuration) Get(key string, result interface{}) interface{} {
	key = strings.ToUpper(key)
	for idx := range c.sources {
		source := c.sources[len(c.sources)-1-idx]
		if source.Exists(key) {
			value := source.Get(key)
			if result == nil {
				return value
			}

			return CastAndTryAssignValue(value, result)
		}
	}

	return result
}

func (c *Configuration) Refresh() bool {
	successfullyRefreshed := true
	defer func() {
		if r := recover(); r != nil {
			successfullyRefreshed = false
			fmt.Println("Recovered. Error:\n", r)
		}
	}()

	wg := sync.WaitGroup{}
	for _, source := range c.sources {
		wg.Add(1)
		go func(sourceArg ISource) {
			defer wg.Done()
			sourceArg.Load()
		}(source)
	}

	wg.Wait()

	return successfullyRefreshed
}

func (c *Configuration) Unmarshal(value interface{}) error {
	rv := reflect.ValueOf(value)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return ErrPointerNotPassed
	}

	keys := make(map[string]ISource)
	for _, source := range c.sources {
		for _, key := range source.GetKeys() {
			keys[key] = source
		}
	}

	flat := make(map[string]interface{})
	for key, source := range keys {
		flat[key] = source.Get(key)
	}

	return nil
}

func (c *Configuration) Deconstruct() IConfiguration {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from error:\n", r)
		}
	}()

	wg := sync.WaitGroup{}
	for _, source := range c.sources {
		wg.Add(1)
		go func(sourceArg ISource) {
			defer wg.Done()
			sourceArg.Deconstruct()
		}(source)
	}
	wg.Wait()

	return c
}
