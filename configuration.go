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
	waitgroup sync.WaitGroup
	sources   []ISource
	delimiter string
	RefreshC  chan ISource
	QuitC     chan interface{}
}

func newConfiguration(sources []ISource) IConfiguration {
	config := &Configuration{
		waitgroup: sync.WaitGroup{},
		sources:   sources,
		delimiter: ".",
		RefreshC:  make(chan ISource),
		QuitC:     make(chan interface{}),
	}

	config.Refresh()
	config.waitgroup.Add(1)
	go config.autoRefresh()

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
	c.QuitC <- struct{}{}
	wg := sync.WaitGroup{}
	for _, source := range c.sources {
		wg.Add(1)
		go func(sourceArg ISource) {
			defer wg.Done()
			sourceArg.Deconstruct()
		}(source)
	}
	wg.Wait()
	c.waitgroup.Done()

	return c
}

func (c *Configuration) autoRefresh() {
	defer c.waitgroup.Done()
	for {
		select {
		case source := <-c.RefreshC:
			if source.Options().ReloadOnChange {
				source.Load()
			}
			if source.Options().SentinelOptions != nil {
				c.newMethod(source)
			}
		case <-c.QuitC:
			return
		}
	}
}

func (c *Configuration) newMethod(source ISource) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from error:\n", r)
		}
	}()

	switch source.Options().SentinelOptions.RefreshPolicy {
	case RefreshAll:
		c.Refresh()
	case RefreshCurrent:
		source.Load()
	case RefreshCurrentAndOver:
		c.refreshCurrendAndAbove(source)
	case RefreshCurrentAndUnder:
		c.refreshCurrendAndUnder(source)
	}
}

func (c *Configuration) refreshCurrendAndAbove(source ISource) {
	wg := sync.WaitGroup{}
	isAbove := false
	for _, s := range c.sources {
		if s == source {
			isAbove = true
		}
		if isAbove {
			wg.Add(1)
			go func(sourceArg ISource) {
				sourceArg.Load()
				defer wg.Done()
			}(source)
		}
	}
	wg.Wait()
}

func (c *Configuration) refreshCurrendAndUnder(source ISource) {
	wg := sync.WaitGroup{}
	isUnder := true
	for _, s := range c.sources {
		if isUnder {
			wg.Add(1)
			go func(sourceArg ISource) {
				sourceArg.Load()
				defer wg.Done()
			}(source)
		}
		if s == source {
			isUnder = false
		}
	}
	wg.Wait()
}
