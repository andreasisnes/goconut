package usersecrets

import (
	"sync"

	"github.com/andreasisnes/goconut"
)

type Options struct {
	goconut.SourceOptions
}

type Source struct {
	goconut.SourceBase
	options   Options
	waitGroup sync.WaitGroup
	quitC     chan interface{}
}

func New(options *Options) goconut.ISource {
	if options == nil {
		options = &Options{}
	}

	return &Source{
		SourceBase: *goconut.NewSourceBase(options.SourceOptions),
		quitC:      make(chan interface{}),
		options:    *options,
		waitGroup:  sync.WaitGroup{},
	}
}

func (source *Source) GetRefreshedValue(key string) interface{} {
	return nil
}

func (source *Source) Load() {

}

func (source *Source) Deconstruct() {

}
