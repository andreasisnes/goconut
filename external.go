package goconut

const (
	RefreshAll RefreshPolicy = iota
	RefreshAllUnder
	RefreshAllOver
	RefreshCurrent
)

type RefreshPolicy int

type SourceOptions struct {
	ReloadOnChange  bool
	Optional        bool
	SentinelOptions *SentinelOptions
}

type SentinelOptions struct {
	Key           string
	RefreshPolicy RefreshPolicy
}

type ISource interface {
	Exists(key string) bool
	Get(key string) interface{}
	GetKeys() []string
	IsDirty() bool
	Options() SourceOptions
	Load()
	Deconstruct()
}
