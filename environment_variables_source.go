package gonfigenvironmentvariables

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/andreasisnes/goconut"
)

type EnvironmentVariablesOptions struct {
	goconut.SourceOptions
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
			Delimiter:       "__",
			RefreshInterval: time.Second,
		}
	}
	if options.Delimiter == "" {
		options.Delimiter = "__"
	}

	env := &EnvironmentVariablesSource{
		EnvOptions: *options,
		WaitGroup:  sync.WaitGroup{},
	}
	goconut.InitSourceBase(&env.SourceBase, &options.SourceOptions)

	if env.SourceOptions.SentinelOptions == nil || env.SourceOptions.ReloadOnChange {
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

func (e *EnvironmentVariablesSource) Deconstruct() {
	e.QuitC <- struct{}{}
	e.WaitGroup.Wait()
}

func (e *EnvironmentVariablesSource) watcher() {
	e.WaitGroup.Add(1)
	defer e.WaitGroup.Done()
	for {
		timer := time.NewTimer(e.EnvOptions.RefreshInterval)
		select {
		case <-timer.C:
			fmt.Println("kicked")
			for _, variable := range os.Environ() {
				keyIdx := strings.Index(variable, "=")
				key := variable[:keyIdx]
				formattedKey := e.formatKey(key)
				fmt.Println(formattedKey)
				if val, ok := e.Flatmap[key]; !ok || val == variable[keyIdx+1:] {
					e.NotifyDirtyness()
					break
				}

				if val, ok := e.Flatmap[formattedKey]; !ok || val == variable[keyIdx+1:] {
					e.NotifyDirtyness()
					break
				}
			}
			continue
		case <-e.QuitC:
			return
		}
	}
}

func (e *EnvironmentVariablesSource) formatKey(key string) string {
	return strings.ToUpper(strings.ReplaceAll(key, e.EnvOptions.Delimiter, "."))
}
