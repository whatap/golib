package udp

import (
	"context"
	"fmt"
	"time"

	"github.com/whatap/golib/config"
	"github.com/whatap/golib/logger"
	// wnet "github.com/whatap/golib/net"
)

const (
	netScheme           = "udp"
	netTimeout          = 60 * time.Second
	netWriteBufferSize  = 64 * 1024
	netFlushMaxWaitTime = 5000
)

type udpClientConfig struct {
	Log            logger.Logger
	ctx            context.Context
	cancel         context.CancelFunc
	Timeout        time.Duration
	NetUdpHost     string
	NetUdpPort     int
	Server         string
	ConfigObserver *config.ConfigObserver
}

func (udpConfig *udpClientConfig) ApplyConfig(conf config.Config) {

}

type UdpClientOption interface {
	apply(*udpClientConfig)
}
type funcUdpClientOption struct {
	f func(*udpClientConfig)
}

func (wfo *funcUdpClientOption) apply(c *udpClientConfig) {
	wfo.f(c)
}

var (
	defaultUdpClientConfig = newDefaultConfig()
)

func newDefaultConfig() *udpClientConfig {
	return &udpClientConfig{
		Log:        &logger.EmptyLogger{},
		Timeout:    netTimeout,
		NetUdpHost: "127.0.0.1",
		NetUdpPort: 6600,
		Server:     fmt.Sprintf("%s:%d", "127.0.0.1", 6600),
	}
}

func newFuncUdpClientOption(f func(*udpClientConfig)) *funcUdpClientOption {
	return &funcUdpClientOption{
		f: f,
	}
}

func WithContext(ctx context.Context, cancel context.CancelFunc) UdpClientOption {
	return newFuncUdpClientOption(func(c *udpClientConfig) {
		c.ctx = ctx
		c.cancel = cancel
	})
}
func WithLogger(logger logger.Logger) UdpClientOption {
	return newFuncUdpClientOption(func(c *udpClientConfig) {
		c.Log = logger
	})
}

func WithConfigObserver(obj *config.ConfigObserver) UdpClientOption {
	return newFuncUdpClientOption(func(c *udpClientConfig) {
		c.ConfigObserver = obj
	})
}

func WithUdpServer(host string, port int) UdpClientOption {
	return newFuncUdpClientOption(func(c *udpClientConfig) {
		c.NetUdpHost = host
		c.NetUdpPort = port
		c.Server = fmt.Sprintf("%s:%d", host, port)
	})
}
