package data

// DataSource represents a generic key-value store for template replacements.
type DataSource interface {
	Get(key string) (string, bool)
	GetList(key string) ([]string, bool)
}

// ChainedSource implements DataSource by polling a list of sources.
// Later sources have higher priority (last-wins merge).
type chainedSource struct {
	sources []DataSource
}

// NewChainedSource creates a DataSource from multiple sources.
func NewChainedSource(sources ...DataSource) DataSource {
	return &chainedSource{sources: sources}
}

func (c *chainedSource) Get(key string) (string, bool) {
	for i := len(c.sources) - 1; i >= 0; i-- {
		if val, ok := c.sources[i].Get(key); ok {
			return val, ok
		}
	}
	return "", false
}

func (c *chainedSource) GetList(key string) ([]string, bool) {
	for i := len(c.sources) - 1; i >= 0; i-- {
		if val, ok := c.sources[i].GetList(key); ok {
			return val, ok
		}
	}
	return nil, false
}

type emptySource struct{}

func (e emptySource) Get(string) (string, bool)         { return "", false }
func (e emptySource) GetList(string) ([]string, bool) { return nil, false }

// NewEmptySource returns a DataSource that always returns missing.
func NewEmptySource() DataSource { return emptySource{} }
