package goconfig

type ISource interface {
	Exists(key string) bool
	Get(key string) interface{}
	GetKeys() []string
	IsDirty() bool
	Load()
	Deconstruct()
}
