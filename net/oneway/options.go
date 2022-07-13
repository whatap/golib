package oneway

import (
	"context"
	"time"

	"github.com/whatap/golib/config"
	"github.com/whatap/golib/logger"
	wnet "github.com/whatap/golib/net"
)

type oneWayTcpClientConfig struct {
	Log            logger.Logger
	ctx            context.Context
	cancel         context.CancelFunc
	Timeout        time.Duration
	License        string
	Servers        []string
	Pcode          int64
	Oid            int32
	UseQueue       bool
	ConfigObserver *config.ConfigObserver
}

type OneWayTcpClientOption interface {
	apply(*oneWayTcpClientConfig)
}
type funcOneWayTcpClientOption struct {
	f func(*oneWayTcpClientConfig)
}

func (wfo *funcOneWayTcpClientOption) apply(c *oneWayTcpClientConfig) {
	wfo.f(c)
}

var (
	defaultOneWayTcpClientConfig = &oneWayTcpClientConfig{
		Timeout:  netTimeout,
		UseQueue: false,
	}
)

func newFuncOneWayTcpClientOption(f func(*oneWayTcpClientConfig)) *funcOneWayTcpClientOption {
	return &funcOneWayTcpClientOption{
		f: f,
	}
}

func WithContext(ctx context.Context, cancel context.CancelFunc) OneWayTcpClientOption {
	return newFuncOneWayTcpClientOption(func(c *oneWayTcpClientConfig) {
		c.ctx = ctx
		c.cancel = cancel
	})
}
func WithLogger(logger logger.Logger) OneWayTcpClientOption {
	return newFuncOneWayTcpClientOption(func(c *oneWayTcpClientConfig) {
		c.Log = logger
	})
}

func WithWhatapTcpServer(info *wnet.WhatapTcpServerInfo) OneWayTcpClientOption {
	return newFuncOneWayTcpClientOption(func(c *oneWayTcpClientConfig) {
		c.License = info.License
		c.Servers = info.Hosts
		c.Pcode = info.Pcode
		c.Oid = info.Oid
	})
}

func WithServers(servers []string) OneWayTcpClientOption {
	return newFuncOneWayTcpClientOption(func(c *oneWayTcpClientConfig) {
		c.Servers = servers
	})
}

func WithPcode(pcode int64) OneWayTcpClientOption {
	return newFuncOneWayTcpClientOption(func(c *oneWayTcpClientConfig) {
		c.Pcode = pcode
	})
}

func WithOid(oid int32) OneWayTcpClientOption {
	return newFuncOneWayTcpClientOption(func(c *oneWayTcpClientConfig) {
		c.Oid = oid
	})
}

func WithLicense(license string) OneWayTcpClientOption {
	return newFuncOneWayTcpClientOption(func(c *oneWayTcpClientConfig) {
		c.License = license
	})
}

func WithUseQueue() OneWayTcpClientOption {
	return newFuncOneWayTcpClientOption(func(c *oneWayTcpClientConfig) {
		c.UseQueue = true
	})
}

func WithConfigObserver(obj *config.ConfigObserver) OneWayTcpClientOption {
	return newFuncOneWayTcpClientOption(func(c *oneWayTcpClientConfig) {
		c.ConfigObserver = obj
	})
}
