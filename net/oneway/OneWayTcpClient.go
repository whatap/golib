package oneway

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"net"
	"reflect"
	"sync"
	"time"

	wio "github.com/whatap/golib/io"
	"github.com/whatap/golib/lang/pack"
	"github.com/whatap/golib/util/dateutil"
	whash "github.com/whatap/golib/util/hash"
	"github.com/whatap/golib/util/queue"

	"github.com/whatap/golib/config"
	"github.com/whatap/golib/logger"
	wnet "github.com/whatap/golib/net"
)

const (
	netSrcAgentOneway   = 10
	netSrcAgentVersion  = 0
	netScheme           = "tcp"
	netTimeout          = 60 * time.Second
	netQueueSize        = 1000
	queueMaxWaitTime    = 5000
	netWriteBufferSize  = 2 * 1024 * 1024
	netFlushMaxWaitTime = 5000
)

var oneWayClient *OneWayTcpClient
var oneWayClientLock sync.Mutex
var oneWayClientSendLock sync.Mutex

type OneWayTcpClient struct {
	ctx              context.Context
	cancel           context.CancelFunc
	License          string        `toml:"license"`
	Servers          []string      `toml:"servers"`
	Pcode            int64         `toml:"project_code"`
	Timeout          time.Duration `toml:"timeout"`
	Log              logger.Logger
	Oname            string
	Oid              int32
	conn             net.Conn
	wr               *bufio.Writer
	hosts            []string
	Queue            *queue.RequestQueue
	lastTime         int64
	flushMaxWaitTime int64
	UseQueue         bool
	ConfigObserver   *config.ConfigObserver
}

func GetOneWayTcpClient(opts ...OneWayTcpClientOption) *OneWayTcpClient {
	if oneWayClient != nil {
		return oneWayClient
	}
	oneWayClientLock.Lock()
	defer oneWayClientLock.Unlock()

	if oneWayClient != nil {
		return oneWayClient
	}
	oneWayClient = newOneWayTcpClient(opts...)
	if err := oneWayClient.Connect(); err != nil {
		fmt.Println("GetOneWayTcpClient  Connect err", err)

	}
	go oneWayClient.process()

	return oneWayClient
}
func newOneWayTcpClient(opts ...OneWayTcpClientOption) *OneWayTcpClient {
	p := new(OneWayTcpClient)
	p.Queue = queue.NewRequestQueue(netQueueSize)

	o := &oneWayTcpClientConfig{}
	for _, opt := range opts {
		opt.apply(o)
	}
	p.Log = o.Log
	if p.Log == nil {
		p.Log = &logger.EmptyLogger{}
	}
	p.License = o.License
	p.Servers = o.Servers
	p.Pcode = o.Pcode
	p.Oid = o.Oid
	p.ctx = o.ctx
	p.cancel = o.cancel
	p.Timeout = o.Timeout
	p.lastTime = dateutil.SystemNow()
	p.UseQueue = o.UseQueue
	// p.Log.Info("newOneWayTcpClient license=", p.License)
	if p.ConfigObserver != nil {
		p.ConfigObserver.Add("oneWayTcpClient", p)
	}

	if p.ctx == nil {
		p.ctx, p.cancel = context.WithCancel(context.Background())
	} else if p.cancel == nil {
		p.ctx, p.cancel = context.WithCancel(context.Background())
	}

	if p.Timeout == 0 {
		p.Timeout = defaultOneWayTcpClientConfig.Timeout
	}

	if p.flushMaxWaitTime == 0 {
		p.flushMaxWaitTime = netFlushMaxWaitTime
	}

	return p
}
func (this *OneWayTcpClient) Connect() error {
	//fmt.Println("Connect host len =", len(this.Servers))
	if this.conn != nil {
		return nil
	}
	// this.Log.Info("Connect host len =", len(this.Servers))
	// Change and connect multiple servers sequentially.
	for _, host := range this.Servers {
		// this.Log.Info("Connect host =", host)
		client, err := net.DialTimeout(netScheme, host, time.Duration(this.Timeout))
		if err != nil {
			this.Log.Errorf("connecting to %q failed: %v", host, err)
			continue
		}
		this.conn = client.(*net.TCPConn)
		this.wr = bufio.NewWriterSize(client, int(netWriteBufferSize))
		this.Log.Infof("Connected %s", host)
		return nil
	}
	return fmt.Errorf("could not connect to any server")
}

func (this *OneWayTcpClient) Close() error {
	if this.conn == nil {
		return nil
	}
	err := this.conn.Close()
	this.conn = nil

	return err
}
func (this *OneWayTcpClient) Destroy() error {
	if this.cancel != nil {
		this.cancel()
	}
	oneWayClient = nil
	return nil
}
func (this *OneWayTcpClient) Send(p pack.Pack, opts ...wnet.TcpClientOption) error {
	return this.SendFlush(p, false, opts...)
}

func (this *OneWayTcpClient) SendFlush(p pack.Pack, flush bool, opts ...wnet.TcpClientOption) error {
	if this.UseQueue {
		ret := this.Queue.Put(&wnet.TcpSend{Flag: 0, Pack: p, Flush: flush, Opts: opts})
		if ret == true {
			return nil
		} else {
			return errors.New("Enqueue Failed")
		}

	} else {
		return this.sendDirect(p, opts...)
	}
}

func (this *OneWayTcpClient) makeData(tcpSend *wnet.TcpSend) *wio.DataOutputX {

	o := &wnet.TcpClientConfig{}

	for _, opt := range tcpSend.Opts {
		opt.Apply(o)
	}

	p := tcpSend.Pack
	dout := wio.NewDataOutputX()
	dout.WriteShort(p.GetPackType())
	p.Write(dout)

	if o.License != "" {
		dout.WriteHeader(netSrcAgentOneway, netSrcAgentVersion, p.GetPCODE(),
			whash.Hash64Str(o.License))
	} else {
		dout.WriteHeader(netSrcAgentOneway, netSrcAgentVersion, p.GetPCODE(),
			whash.Hash64Str(this.License))
	}

	return dout
}

func (this *OneWayTcpClient) sendDirect(p pack.Pack, opts ...wnet.TcpClientOption) error {
	oneWayClientSendLock.Lock()
	defer oneWayClientSendLock.Unlock()
	// this.Log.Infof("sendDirect lic=%s, pcode=%d, oid=%d", this.License, this.Pcode, this.Oid)

	dout := this.makeData(&wnet.TcpSend{Pack: p, Opts: opts})

	if err := this.send(dout.ToByteArray()); err != nil {
		_ = this.Close()
		return err
	}
	if _, err := this.Flush(); err != nil {
		return fmt.Errorf("cannot flush : %v", err)
	}
	return nil
}

func (this *OneWayTcpClient) send(sendbuf []byte) (err error) {
	defer func() {
		if r := recover(); r != nil {
			this.Log.Errorf("recover panic %v", r)
		}
	}()

	if this.conn == nil {
		if err := this.Connect(); err != nil {
			return fmt.Errorf("cannot connect %v", err)
		}
	}

	for pos := 0; pos < len(sendbuf); {
		deadline := time.Now().Add(time.Duration(this.Timeout))
		if err := this.conn.SetWriteDeadline(deadline); err != nil {
			return fmt.Errorf("cannot set write deadline: %v", err)
		}
		nbytethistime, err := this.wr.Write(sendbuf[pos:])
		if err != nil {
			return fmt.Errorf("buffered writer cannot write : %v", err)
		}
		pos += nbytethistime
	}
	return nil
}

func (this *OneWayTcpClient) SendAndClear() error {
	// this.Log.Info("OneWayTcpClient SendAndClear\n")
	if this.Queue.Size() > 0 {
		tmp := this.Queue.GetNoWait()
		for tmp != nil {
			if tcpSend, ok := tmp.(*wnet.TcpSend); ok {

				dout := this.makeData(tcpSend)

				if err := this.send(dout.ToByteArray()); err != nil {
					_ = this.Close()
					return err
				} else {
					// this.Log.Info("send pack t=", p.GetPackType())
				}
			}
			tmp = this.Queue.GetNoWait()
		}
	}
	if _, err := this.Flush(); err != nil {
		return fmt.Errorf("cannot flush : %v", err)
	}
	return nil
}

func (this *OneWayTcpClient) Flush() (n int, err error) {
	n = this.wr.Buffered()
	if err = this.wr.Flush(); err != nil {
		return 0, err
	}
	return n, nil
}

func (this *OneWayTcpClient) process() {
	for {
		select {
		case <-this.ctx.Done():
			return
		default:
			if this.conn == nil {
				if err := this.Connect(); err != nil {
					time.Sleep(5 * time.Second)
					continue
				}
			}
			if tmp := this.Queue.GetTimeout(int(queueMaxWaitTime)); tmp != nil {
				if tcpSend, ok := tmp.(*wnet.TcpSend); ok {
					dout := this.makeData(tcpSend)

					if err := this.send(dout.ToByteArray()); err != nil {
						_ = this.Close()
						//return err
					}
					if tcpSend.Flush || this.lastTime-this.flushMaxWaitTime > 5000 {
						if _, err := this.Flush(); err != nil {
							this.Close()
						}
					}
				}
			}
		}
	}
}

// implements common.ConfigObserver
func (this *OneWayTcpClient) ApplyConfig(conf config.Config) {
	//this.Pcode = conf.GetLong("pcode")
	license := conf.GetValue("license")
	host := conf.GetValue("whatap.server.host")
	port := conf.GetValueDef("whatap.server.port", "6600")
	servers := wnet.GetWhatapHosts(host, port)

	this.Pcode = conf.GetLong("pcode", 0)
	this.Oid = conf.GetInt("oid", 0)
	this.Timeout = time.Duration(conf.GetInt("", 60000)) * time.Millisecond
	this.lastTime = dateutil.SystemNow()

	if this.License != license || !reflect.DeepEqual(this.Servers, servers) {
		this.License = license
		this.Servers = servers
		this.Close()
		this.Connect()
	}
}
