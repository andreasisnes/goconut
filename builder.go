package goconfig

type IBuilder interface {
	Clear()
	Add(source ISource) *Builder
	Sources() []ISource
}

type Builder struct {
	sources []ISource
}

func NewBuilder() *Builder {
	return &Builder{
		sources: make([]ISource, 0),
	}
}

func (c *Builder) Clear() {
	c.sources = make([]ISource, 0)
}

func (c *Builder) Add(source ISource) *Builder {
	c.sources = append(c.sources, source)
	return c
}

func (c *Builder) Sources() []ISource {
	return c.sources
}

func (c *Builder) Build() IConfiguration {
	return newConfiguration(c.sources)
}
