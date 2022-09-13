package logfile

import (
	"context"
	_ "time"

	"github.com/whatap/golib/config"
)

type fileLoggerConfig struct {
	ctx            context.Context
	cancel         context.CancelFunc
	configObserver *config.ConfigObserver

	level           int
	rotationEnabled bool
	keepDays        int
	cacheInterval   int

	oname    string
	logID    string
	homePath string

	IsStdout bool
}

func defaultFileLoggerConfig() fileLoggerConfig {
	return fileLoggerConfig{
		level:           defaultLevel,
		rotationEnabled: defaultRotationEnabled,
		keepDays:        defaultKeepDays,
		cacheInterval:   defaultCacheInterval,
		oname:           defaultOname,
		logID:           defaultLogID,
	}
}

type FileLoggerOption interface {
	apply(*fileLoggerConfig)
}
type funcFileLoggerOption struct {
	f func(*fileLoggerConfig)
}

func (wfo *funcFileLoggerOption) apply(c *fileLoggerConfig) {
	wfo.f(c)
}

func newFuncFileLoggerOption(f func(*fileLoggerConfig)) *funcFileLoggerOption {
	return &funcFileLoggerOption{
		f: f,
	}
}

func WithContext(ctx context.Context, cancel context.CancelFunc) FileLoggerOption {
	return newFuncFileLoggerOption(func(c *fileLoggerConfig) {
		c.ctx = ctx
		c.cancel = cancel
	})
}

func WithConfigObserver(obj *config.ConfigObserver) FileLoggerOption {
	return newFuncFileLoggerOption(func(c *fileLoggerConfig) {
		c.configObserver = obj
	})
}

func WithHomePath(home string) FileLoggerOption {
	return newFuncFileLoggerOption(func(c *fileLoggerConfig) {
		c.homePath = home
	})
}

func WithOnameLogID(oname, logID string) FileLoggerOption {
	return newFuncFileLoggerOption(func(c *fileLoggerConfig) {
		c.oname = oname
		c.logID = logID
	})
}

func WithLevel(lv int) FileLoggerOption {
	return newFuncFileLoggerOption(func(c *fileLoggerConfig) {
		c.level = lv
	})
}
func WithStdout(b bool) FileLoggerOption {
	return newFuncFileLoggerOption(func(c *fileLoggerConfig) {
		c.IsStdout = b
	})
}
