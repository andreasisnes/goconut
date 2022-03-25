package goconut

type ISource interface {
	Exists(key string) bool
	Get(key string) interface{}
	GetKeys() []string
	Options() SourceOptions
	Connect(configuration *Configuration)
	Load()
	Deconstruct()
}

func NewSourceBase(c *Configuration) ISource {
	return &SourceBase{}
}

type SourceBase struct {
	Flatmap        map[string]interface{}
	Configurations []*Configuration
	SourceOptions  SourceOptions
	RefreshC       chan ISource
}

func InitSourceBase(base *SourceBase, options *SourceOptions) {
	base.Flatmap = make(map[string]interface{})
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
	return s.Get(key) != nil
}

func (s *SourceBase) Get(key string) (value interface{}) {
	if s.Flatmap == nil {
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

func (s *SourceBase) Deconstruct() {

}
