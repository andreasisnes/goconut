package gonfigenvironmentvariables

import (
	"os"
	"strings"
	"sync"
	"time"

	"github.com/andreasisnes/goconut"
)

const (
	DefaultDelimiter = "__"
)

type EnvironmentVariablesOptions struct {
	goconut.SourceOptions
	Prefix          string
	Delimiter       string
	RefreshInterval time.Duration
}

type EnvironmentVariablesSource struct {
	goconut.SourceBase
	EnvOptions EnvironmentVariablesOptions
	WaitGroup  sync.WaitGroup
	QuitC      chan interface{}
}

func NewEnvironmentVariablesSource(options *EnvironmentVariablesOptions) goconut.ISource {
	if options == nil {
		options = &EnvironmentVariablesOptions{
			Delimiter:       DefaultDelimiter,
			RefreshInterval: time.Second,
		}
	}

	if options.Delimiter == "" {
		options.Delimiter = DefaultDelimiter
	}

	env := &EnvironmentVariablesSource{
		SourceBase: *goconut.NewSourceBase(options.SourceOptions),
		QuitC:      make(chan interface{}),
		EnvOptions: *options,
		WaitGroup:  sync.WaitGroup{},
	}

	if env.SourceOptions.SentinelOptions != nil || env.SourceOptions.ReloadOnChange {
		go env.watcher()
	}

	return env
}

func (e *EnvironmentVariablesSource) Load() {
	e.RWTex.Lock()
	defer e.RWTex.Unlock()
	for _, variable := range os.Environ() {
		keyIdx := strings.Index(variable, "=")
		e.Flatmap[variable[:keyIdx]] = variable[keyIdx+1:]
		e.Flatmap[e.formatKey(variable[:keyIdx])] = variable[keyIdx+1:]
	}
}

func (e *EnvironmentVariablesSource) GetRefreshedValue(key string) interface{} {
	return e.Flatmap[key]
}

func (e *EnvironmentVariablesSource) Deconstruct() {
	e.QuitC <- struct{}{}
}

func (e *EnvironmentVariablesSource) watcher() {
	e.WaitGroup.Add(1)
	defer e.WaitGroup.Done()
	timer := time.NewTimer(e.EnvOptions.RefreshInterval)
	for {
		select {
		case <-e.QuitC:
			return
		case <-timer.C:
			e.RWTex.RLock()
			for _, variable := range os.Environ() {
				keyIdx := strings.Index(variable, "=")

				key := variable[:keyIdx]
				if val, ok := e.Flatmap[key]; !ok || val != variable[keyIdx+1:] {
					e.NotifyDirtyness()
					break
				}

				formattedKey := e.formatKey(key)
				if val, ok := e.Flatmap[formattedKey]; !ok || val != variable[keyIdx+1:] {
					e.NotifyDirtyness()
					break
				}
			}
			e.RWTex.RUnlock()
		}
		timer.Reset(e.EnvOptions.RefreshInterval)
	}
}

func (e *EnvironmentVariablesSource) formatKey(key string) string {
	key = strings.ToUpper(strings.ReplaceAll(key, e.EnvOptions.Delimiter, "."))
	prefix := strings.ToUpper(e.EnvOptions.Prefix)
	return strings.TrimPrefix(key, prefix)
}
