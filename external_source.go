package goconut

import (
	"sync"
)

type ISource interface {
	Exists(key string) bool
	Get(key string) interface{}
	GetKeys() []string
	Options() SourceOptions
	Connect(configuration *Configuration)
	Load()
	Deconstruct(configuration *Configuration)
}

func NewSourceBase(c *Configuration) ISource {
	return &SourceBase{}
}

type SourceBase struct {
	Flatmap        map[string]interface{}
	RWTex          *sync.RWMutex
	Configurations []*Configuration
	SourceOptions  SourceOptions
}

func InitSourceBase(base *SourceBase, options *SourceOptions) {
	base.Flatmap = make(map[string]interface{})
	base.RWTex = &sync.RWMutex{}
	base.Configurations = make([]*Configuration, 0)
	base.SourceOptions = *options
}

func (s *SourceBase) Connect(configuration *Configuration) {
	for _, c := range s.Configurations {
		if c == configuration {
			return
		}
	}

	if s.Configurations == nil {
		InitSourceBase(s, &s.SourceOptions)
	}

	s.Configurations = append(s.Configurations, configuration)
}

func (s *SourceBase) Exists(key string) bool {
	s.RWTex.RLock()
	defer s.RWTex.RUnlock()

	return s.Get(key) != nil
}

func (s *SourceBase) Get(key string) (value interface{}) {
	s.RWTex.Lock()
	defer s.RWTex.Unlock()

	if s.Flatmap == nil || s.RWTex == nil {
		InitSourceBase(s, &s.SourceOptions)
		return nil
	}

	if value, ok := s.Flatmap[key]; ok {
		return value
	}

	return nil
}

func (s *SourceBase) NotifyDirtyness() {
	for _, v := range s.Configurations {
		v.RefreshC <- s
	}
}

func (s *SourceBase) GetKeys() (result []string) {
	s.RWTex.RLock()
	defer s.RWTex.RUnlock()

	result = make([]string, 0)
	for key := range s.Flatmap {
		result = append(result, key)
	}

	return result
}

func (s *SourceBase) Options() SourceOptions {
	return s.SourceOptions
}

func (s *SourceBase) Load() {
}

func (s *SourceBase) Deconstruct(configuration *Configuration) {
}
