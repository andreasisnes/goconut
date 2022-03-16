package goconut

const (
	RefreshAll RefreshPolicy = iota
	RefreshCurrentAndOver
	RefreshCurrentAndUnder
	RefreshCurrent
)

type SourceOptions struct {
	ReloadOnChange  bool
	Optional        bool
	SentinelOptions *SentinelOptions
}

type SentinelOptions struct {
	Key           string
	RefreshPolicy RefreshPolicy
}

type RefreshPolicy int
