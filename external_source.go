package goconut

import (
	"sync"
)

type ISource interface {
	Connect(refreshC chan ISource)
	Exists(key string) bool
	Get(key string) interface{}
	GetKeys() []string
	Options() SourceOptions

	GetRefreshedValue(key string) interface{}
	Load()
	Deconstruct()
}

type SourceBase struct {
	RefreshC      chan ISource
	Flatmap       map[string]interface{}
	RWTex         sync.RWMutex
	SourceOptions SourceOptions
}

// Used by external Sources
func NewSourceBase(options SourceOptions) *SourceBase {
	return &SourceBase{
		Flatmap:       make(map[string]interface{}),
		RWTex:         sync.RWMutex{},
		SourceOptions: options,
	}
}

// Used by Configuration
func (s *SourceBase) Connect(refreshC chan ISource) {
	s.RWTex.Lock()
	defer s.RWTex.Unlock()

	s.RefreshC = refreshC
}

// Checks if a key exists
func (s *SourceBase) Exists(key string) bool {
	return s.Get(key) != nil
}

// Get Config Values
func (s *SourceBase) Get(key string) (value interface{}) {
	s.RWTex.RLock()
	defer s.RWTex.RUnlock()

	if value, ok := s.Flatmap[key]; ok {
		return value
	}

	return nil
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

func (source *SourceBase) Options() SourceOptions {
	return source.SourceOptions
}

func (source *SourceBase) NotifyDirtyness(externalSource ISource) {
	source.RWTex.RLock()
	defer source.RWTex.RUnlock()

	if source.RefreshC != nil {
		source.RefreshC <- externalSource
	}
}
