package conffile

import (
	"context"
	_ "time"

	"github.com/whatap/golib/config"
	"github.com/whatap/golib/logger"
)

type fileConfigConfig struct {
	ctx            context.Context
	cancel         context.CancelFunc
	configObserver *config.ConfigObserver
	Log            logger.Logger
	Parser         FileParser

	// for file
	prefix      string
	suffix      string
	excludeKeys []string

	homePath string
}

func defaultFileConfigConfig() *fileConfigConfig {
	return &fileConfigConfig{
		Log:         &logger.EmptyLogger{},
		Parser:      NewDefaultFileParser(),
		excludeKeys: make([]string, 0),
	}
}

type FileConfigOption interface {
	apply(*fileConfigConfig)
}
type funcFileConfigOption struct {
	f func(*fileConfigConfig)
}

func (wfo *funcFileConfigOption) apply(c *fileConfigConfig) {
	wfo.f(c)
}

func newFuncFileConfigOption(f func(*fileConfigConfig)) *funcFileConfigOption {
	return &funcFileConfigOption{
		f: f,
	}
}

func WithContext(ctx context.Context, cancel context.CancelFunc) FileConfigOption {
	return newFuncFileConfigOption(func(c *fileConfigConfig) {
		c.ctx = ctx
		c.cancel = cancel
	})
}
func WithLogger(logger logger.Logger) FileConfigOption {
	return newFuncFileConfigOption(func(c *fileConfigConfig) {
		c.Log = logger
	})
}
func WithConfigObserver(obj *config.ConfigObserver) FileConfigOption {
	return newFuncFileConfigOption(func(c *fileConfigConfig) {
		c.configObserver = obj
	})
}

func WithHomePath(home string) FileConfigOption {
	return newFuncFileConfigOption(func(c *fileConfigConfig) {
		c.homePath = home
	})
}

func WithParser(ps FileParser) FileConfigOption {
	return newFuncFileConfigOption(func(c *fileConfigConfig) {
		if ps == nil {
			panic("FileParser is nil ")
		}
		c.Parser = ps
	})
}
func WithPrefix(prefix string) FileConfigOption {
	return newFuncFileConfigOption(func(c *fileConfigConfig) {
		c.prefix = prefix
	})
}

func WithSuffix(suffix string) FileConfigOption {
	return newFuncFileConfigOption(func(c *fileConfigConfig) {
		c.suffix = suffix
	})
}

func WithExcludeKeys(keys []string) FileConfigOption {
	return newFuncFileConfigOption(func(c *fileConfigConfig) {
		c.excludeKeys = keys
	})
}
