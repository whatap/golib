package net

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/whatap/golib/lang/pack"

	whash "github.com/whatap/golib/util/hash"
	"github.com/whatap/golib/util/stringutil"
)

type TcpClient interface {
	Connect() error
	Send(p pack.Pack, opts ...TcpClientOption) error
	SendFlush(p pack.Pack, flush bool, opts ...TcpClientOption) error
	Close() error
}

// Emtyp TcpClient , use without nil check
type EmptyTcpClient struct{}

func (tc *EmptyTcpClient) Connect() error                                  { return nil }
func (tc *EmptyTcpClient) Send(p pack.Pack, opts ...TcpClientOption) error { return nil }
func (tc *EmptyTcpClient) SendFlush(p pack.Pack, flush bool, opts ...TcpClientOption) error {
	return nil
}
func (tc *EmptyTcpClient) Close() error { return nil }

type Initializer interface {
	Init() error
}

type TcpClientConfig struct {
	License string
}

type TcpClientOption interface {
	Apply(conf *TcpClientConfig)
}
type funcTcpClientOption struct {
	f func(*TcpClientConfig)
}

func (this *funcTcpClientOption) Apply(c *TcpClientConfig) {
	this.f(c)
}
func newFuncTcpClientOption(f func(*TcpClientConfig)) TcpClientOption {
	return &funcTcpClientOption{
		f: f,
	}
}

func WithLicense(license string) TcpClientOption {
	return newFuncTcpClientOption(func(c *TcpClientConfig) {
		c.License = license
	})
}

type TcpSend struct {
	Flag  byte
	Pack  pack.Pack
	Flush bool
	Opts  []TcpClientOption
}

type WhatapTcpServerInfo struct {
	License string
	Hosts   []string
	Pcode   int64
	Oid     int32
}

func NewWhatapTcpServerInfo(license, host, portStr, pcode, oname string) *WhatapTcpServerInfo {
	p := new(WhatapTcpServerInfo)
	p.License = license
	p.Hosts = GetWhatapHosts(host, portStr)
	if v, err := strconv.ParseInt(pcode, 10, 64); err == nil {
		p.Pcode = v
	}
	p.Oid = whash.HashStr(oname)
	return p
}

func GetWhatapHosts(host string, portStr string) []string {
	port := 6600
	if v, err := strconv.Atoi(portStr); err != nil {
		if v != 0 {
			port = v
		}
	}
	arr := stringutil.Tokenizer(host, "/:,")
	servers := make([]string, 0)
	for _, it := range arr {
		server := fmt.Sprintf("tcp://%s:%d", it, port)
		u, err := url.Parse(server)
		if err != nil {
			//this.Log.Errorf("invalid address: %s", server)
			continue
		}
		if u.Scheme != "tcp" {
			//this.Log.Errorf("only tcp is supported: %s", server)
			continue
		}
		//fmt.Println("whatap host ", fmt.Sprintf("%s:%d", it, port))
		servers = append(servers, fmt.Sprintf("%s:%d", it, port))
	}
	return servers
}
