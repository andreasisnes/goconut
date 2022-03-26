package gonfigenvironmentvariables

import (
	"os"
	"strings"
	"sync"
	"time"

	"github.com/andreasisnes/goconut"
)

type EnvironmentVariablesOptions struct {
	goconut.SourceOptions
	Delimiter       string
	Prefix          string
	RefreshInterval time.Duration
}

type EnvironmentVariablesSource struct {
	goconut.SourceBase
	RWTex      sync.RWMutex
	EnvOptions EnvironmentVariablesOptions
	WaitGroup  sync.WaitGroup
	QuitC      chan interface{}
}

func NewEnvironmentVariablesSource(options *EnvironmentVariablesOptions) goconut.ISource {
	if options == nil {
		options = &EnvironmentVariablesOptions{
			Delimiter:       "__",
			RefreshInterval: time.Second,
		}
	}
	if options.Delimiter == "" {
		options.Delimiter = "__"
	}

	env := &EnvironmentVariablesSource{
		QuitC:      make(chan interface{}),
		EnvOptions: *options,
		WaitGroup:  sync.WaitGroup{},
		RWTex:      sync.RWMutex{},
	}
	goconut.InitSourceBase(&env.SourceBase, &options.SourceOptions)

	if env.SourceOptions.SentinelOptions != nil || env.SourceOptions.ReloadOnChange {
		go env.watcher()
	}

	return env
}

func (e *EnvironmentVariablesSource) Load() {
	for _, variable := range os.Environ() {
		keyIdx := strings.Index(variable, "=")
		e.Flatmap[variable[:keyIdx]] = variable[keyIdx+1:]
		e.Flatmap[e.formatKey(variable[:keyIdx])] = variable[keyIdx+1:]
	}
}

func (e *EnvironmentVariablesSource) Deconstruct(configuration *goconut.Configuration) {
	for idx, c := range e.Configurations {
		if c == configuration {
			e.Configurations = append(e.Configurations[:idx], e.Configurations[idx+1:]...)
		}
	}

	if len(e.Configurations) == 0 {
		e.QuitC <- struct{}{}
		e.WaitGroup.Wait()
	}
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
			for _, variable := range os.Environ() {
				keyIdx := strings.Index(variable, "=")
				key := variable[:keyIdx]
				formattedKey := e.formatKey(key)
				if val, ok := e.Flatmap[key]; !ok || val != variable[keyIdx+1:] {
					e.NotifyDirtyness()
					break
				}

				if val, ok := e.Flatmap[formattedKey]; !ok || val != variable[keyIdx+1:] {
					e.NotifyDirtyness()
					break
				}
			}
		}
		timer.Reset(e.EnvOptions.RefreshInterval)
	}
}

func (e *EnvironmentVariablesSource) formatKey(key string) string {
	key = strings.ToUpper(strings.ReplaceAll(key, e.EnvOptions.Delimiter, "."))
	prefix := strings.ToUpper(e.EnvOptions.Prefix)
	return strings.TrimPrefix(key, prefix)
}
