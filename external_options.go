package goconut

const (
	RefreshAll RefreshPolicy = iota
	RefreshCurrentAndOver
	RefreshCurrentAndUnder
	RefreshCurrent
)

type SourceOptions struct {
	Optional        bool
	ReloadOnChange  bool
	SentinelOptions *SentinelOptions
}

type SentinelOptions struct {
	Key           string
	RefreshPolicy RefreshPolicy
}

type RefreshPolicy int
