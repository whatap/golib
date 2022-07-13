package zip

import (
	"context"

	"github.com/whatap/golib/config"
	"github.com/whatap/golib/logger"
	wnet "github.com/whatap/golib/net"
)

type zipSendProxyThreadConfig struct {
	ctx                  context.Context
	cancel               context.CancelFunc
	packCount            int
	firstTime            int64
	client               wnet.TcpClient
	Log                  logger.Logger
	useQueue             bool
	logsinkMaxWaitTime   int64
	logsinkQueueSize     int
	logsinkMaxBufferSize int
	logsinkZipMinSize    int
	configObserver       *config.ConfigObserver
}

type ZipSendProxyThreadOption interface {
	apply(*zipSendProxyThreadConfig)
}
type funcZipSendProxyThreadOption struct {
	f func(*zipSendProxyThreadConfig)
}

func (wfo *funcZipSendProxyThreadOption) apply(c *zipSendProxyThreadConfig) {
	wfo.f(c)
}

func newFuncZipSendProxyThreadOption(f func(*zipSendProxyThreadConfig)) *funcZipSendProxyThreadOption {
	return &funcZipSendProxyThreadOption{
		f: f,
	}
}

func WithContext(ctx context.Context, cancel context.CancelFunc) ZipSendProxyThreadOption {
	return newFuncZipSendProxyThreadOption(func(c *zipSendProxyThreadConfig) {
		c.ctx = ctx
		c.cancel = cancel
	})
}
func WithLogger(logger logger.Logger) ZipSendProxyThreadOption {
	return newFuncZipSendProxyThreadOption(func(c *zipSendProxyThreadConfig) {
		c.Log = logger
	})
}

func WithTcpClient(tcpClient wnet.TcpClient) ZipSendProxyThreadOption {
	return newFuncZipSendProxyThreadOption(func(c *zipSendProxyThreadConfig) {
		c.client = tcpClient
	})
}
func WithUseQueue() ZipSendProxyThreadOption {
	return newFuncZipSendProxyThreadOption(func(c *zipSendProxyThreadConfig) {
		c.useQueue = true
	})
}

func WithConfigObserver(obj *config.ConfigObserver) ZipSendProxyThreadOption {
	return newFuncZipSendProxyThreadOption(func(c *zipSendProxyThreadConfig) {
		c.configObserver = obj
	})
}
