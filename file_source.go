package file

import (
	"encoding/json"
	"log"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/andreasisnes/goconut"
	"github.com/andreasisnes/goflat"
	"github.com/fsnotify/fsnotify"
	"github.com/pelletier/go-toml"
	"gopkg.in/yaml.v3"
)

const (
	DefaultFile = "settings.json"
)

type Options struct {
	goconut.SourceOptions
	Path string
}

type FileSource struct {
	goconut.SourceBase
	FileOptions   Options
	Configuration map[string]interface{}
	Content       []byte
	WaitGroup     sync.WaitGroup
	QuitC         chan interface{}
}

func NewFileSource(options *Options) goconut.ISource {
	if options == nil {
		options = &Options{
			Path: DefaultFile,
		}
	}

	source := &FileSource{
		SourceBase:    *goconut.NewSourceBase(options.SourceOptions),
		FileOptions:   *options,
		Configuration: make(map[string]interface{}),
		QuitC:         make(chan interface{}),
		WaitGroup:     sync.WaitGroup{},
	}

	if source.SourceOptions.SentinelOptions != nil || source.SourceOptions.ReloadOnChange {
		go source.watcher()
	}

	return source
}

func (source *FileSource) GetRefreshedValue(key string) interface{} {
	return nil
}

func (source *FileSource) Deconstruct() {
	source.QuitC <- struct{}{}
}

func (source *FileSource) Load() {
	if fileExists(source.FileOptions.Path) {
		if content, err := os.ReadFile(source.FileOptions.Path); err != nil {
			panic(err)
		} else {
			source.Content = content
			extension := strings.ToLower(path.Ext(source.FileOptions.Path))
			switch extension {
			case ".json":
				source.unmarshal(json.Unmarshal)
			case ".yml", ".yaml":
				source.unmarshal(yaml.Unmarshal)
			case ".toml":
				source.unmarshal(toml.Unmarshal)
			default:
				if !source.FileOptions.Optional {
					log.Fatalf("'%s' is not a <.json, yml, yaml, toml> file", path.Base(source.FileOptions.Path))
				}
			}
		}
	} else {
		if !source.FileOptions.Optional {
			panic(os.ErrNotExist)
		}
	}
}

func (source *FileSource) watcher() {
	source.WaitGroup.Add(1)
	defer source.WaitGroup.Done()

	watcher, shouldReturn := source.createWatcher()
	if shouldReturn {
		return
	}
	defer watcher.Close()

	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write {
				source.NotifyDirtyness(source)
			}
		case err := <-watcher.Errors:
			if source.FileOptions.Optional {
				panic(err)
			}
			return
		case <-source.QuitC:
			return
		}
	}
}

func (source *FileSource) createWatcher() (*fsnotify.Watcher, bool) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		if !source.FileOptions.Optional {
			panic(err)
		}

		return nil, true
	}

	if !fileExists(source.FileOptions.Path) {
		if !source.FileOptions.Optional {
			panic(os.ErrNotExist)
		}

		return nil, true
	}

	if err = watcher.Add(source.FileOptions.Path); err != nil {
		if !source.FileOptions.Optional {
			panic(os.ErrNotExist)
		}

		return nil, true
	}

	return watcher, false
}

func (source *FileSource) unmarshal(fn func([]byte, interface{}) error) error {
	if err := fn(source.Content, &source.Configuration); err != nil {
		return err
	}

	source.Flatmap = goflat.Map(source.Configuration, &goflat.Options{
		Delimiter: goflat.DefaultDelimiter,
		Fold:      goflat.UpperCaseFold,
	})

	return nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}
