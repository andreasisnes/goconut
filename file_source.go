package file

import (
	"encoding/json"
	"errors"
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

var (
	ErrFileProviderFileNotExist = errors.New("provided argument 'filename' does not exists")
)

type FileSource struct {
	Shutdown       chan struct{}
	Content        []byte
	Filename       string
	Optional       bool
	ReloadOnChange bool
	DirtyFile      bool
	WaitGroup      *sync.WaitGroup
	Config         map[string]interface{}
	ConfigFlat     map[string]interface{}
}

func NewFileProvider(filename string, optional, reloadOnChange bool) goconut.ISource {
	fs := &FileSource{
		WaitGroup:      &sync.WaitGroup{},
		Filename:       filename,
		Optional:       optional,
		ReloadOnChange: reloadOnChange,
		DirtyFile:      false,
		Content:        make([]byte, 0),
		Config:         make(map[string]interface{}),
		ConfigFlat:     make(map[string]interface{}),
	}

	if reloadOnChange {
		go fs.watcher()
	}

	return fs
}

func (f *FileSource) Exists(key string) bool {
	return f.Get(key) != nil
}

func (f *FileSource) Get(key string) interface{} {
	if val, ok := f.ConfigFlat[key]; ok {
		return val
	}

	return nil
}

func (f *FileSource) IsDirty() bool {
	return f.DirtyFile
}

func (f *FileSource) GetKeys() []string {
	res := make([]string, 0)
	for k := range f.ConfigFlat {
		res = append(res, k)
	}

	return res
}

func (f *FileSource) Deconstruct() {
	f.Shutdown <- struct{}{}
	f.WaitGroup.Wait()
}

func (f *FileSource) Load() {
	if fileExists(f.Filename) {
		if content, err := os.ReadFile(f.Filename); err != nil {
			panic(err)
		} else {
			f.Content = content
			extension := strings.ToLower(path.Ext(f.Filename))
			switch extension {
			case ".json":
				f.unmarshal(json.Unmarshal)
			case ".yml", ".yaml":
				f.unmarshal(yaml.Unmarshal)
			case ".toml":
				f.unmarshal(toml.Unmarshal)
			default:
				if !f.Optional {
					log.Fatalf("'%s' is not a <.json, yml, yaml, toml> file", path.Base(f.Filename))
				}
			}
			f.DirtyFile = false
		}
	} else {
		if f.Optional {
			panic(ErrFileProviderFileNotExist)
		}
	}
}

func (f *FileSource) watcher() {
	f.WaitGroup.Add(1)
	defer f.WaitGroup.Done()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	if !fileExists(f.Filename) {
		return
	}
	err = watcher.Add(f.Filename)
	if err != nil {
		return
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				f.DirtyFile = true
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println(err)
		case <-f.Shutdown:
			return
		}
	}
}

func (f *FileSource) unmarshal(fn func([]byte, interface{}) error) error {
	if err := fn(f.Content, &f.Config); err != nil {
		return err
	}

	f.ConfigFlat = goflat.Map(f.Config, &goflat.Options{
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
