package zip

import (
	"bytes"
	//"fmt"
	"context"
	"sync"

	"github.com/whatap/golib/config"
	"github.com/whatap/golib/io"
	"github.com/whatap/golib/lang/pack"
	"github.com/whatap/golib/logger"
	wnet "github.com/whatap/golib/net"
	"github.com/whatap/golib/util/dateutil"
	"github.com/whatap/golib/util/queue"
)

const (
	defaultLogsinkMaxWaitTime   = 5000
	defaultLogsinkQueueSize     = 1000
	defaultLogsinkMaxBufferSize = 1024 * 64
	defaultLogsinkZipMinSize    = 100
)

var zipSendProxyThread *ZipSendProxyThread
var zipSendProxyThreadMutex = sync.Mutex{}

type ZipSendProxyThread struct {
	ctx                  context.Context
	cancel               context.CancelFunc
	Log                  logger.Logger
	configObserver       *config.ConfigObserver
	Queue                *queue.RequestQueue
	buffer               bytes.Buffer
	packCount            int
	firstTime            int64
	client               wnet.TcpClient
	useQueue             bool
	logsinkMaxWaitTime   int64
	logsinkQueueSize     int
	logsinkMaxBufferSize int
	logsinkZipMinSize    int
}

var (
	defaultZipSendProxyThreadConfig = &zipSendProxyThreadConfig{}
)

func GetInstance(opts ...ZipSendProxyThreadOption) *ZipSendProxyThread {
	zipSendProxyThreadMutex.Lock()
	defer zipSendProxyThreadMutex.Unlock()
	if zipSendProxyThread != nil {
		return zipSendProxyThread
	}
	p := new(ZipSendProxyThread)
	p.logsinkMaxWaitTime = defaultLogsinkMaxWaitTime
	p.logsinkQueueSize = defaultLogsinkQueueSize
	p.logsinkMaxBufferSize = defaultLogsinkMaxBufferSize
	p.logsinkZipMinSize = defaultLogsinkZipMinSize

	o := &zipSendProxyThreadConfig{}
	for _, opt := range opts {
		opt.apply(o)
	}
	if o.ctx == nil {
		p.ctx, p.cancel = context.WithCancel(context.Background())
	} else {
		if o.cancel == nil {
			p.ctx, p.cancel = context.WithCancel(o.ctx)
		} else {
			p.ctx = o.ctx
			p.cancel = o.cancel
		}
	}
	p.useQueue = o.useQueue
	p.logsinkMaxWaitTime = o.logsinkMaxWaitTime
	p.logsinkQueueSize = o.logsinkQueueSize
	p.logsinkMaxBufferSize = o.logsinkMaxBufferSize
	p.logsinkZipMinSize = o.logsinkZipMinSize

	p.client = o.client
	if p.client == nil {
		p.client = &wnet.EmptyTcpClient{}
	}

	p.Log = o.Log
	if p.Log == nil {
		p.Log = &logger.EmptyLogger{}
	}

	p.configObserver = o.configObserver
	if p.configObserver != nil {
		p.configObserver.Add("zipSendProxyThread", p)
	}
	zipSendProxyThread = p
	if p.useQueue {
		p.Queue = queue.NewRequestQueue(int(p.logsinkQueueSize))
		go zipSendProxyThread.run()
	}

	return zipSendProxyThread
}

// set tcp channel to write
func (this *ZipSendProxyThread) SetTcpClient(c wnet.TcpClient) {
	this.client = c
}

func (this *ZipSendProxyThread) Add(p *pack.LogSinkPack) {
	this.Queue.Put(p)
}

func (this *ZipSendProxyThread) run() {
	for true {
		select {
		case <-this.ctx.Done():
			// this.Log.Info("ZipSend stoped")
			this.sendAndClear()
			return
		default:
			// this.Log.Info("ZipSend Run  ", this.Queue.Size())
			if tmp := this.Queue.GetTimeout(int(this.logsinkMaxWaitTime)); tmp != nil {
				data := tmp.(*pack.LogSinkPack)
				if data != nil {
					this.Append(data)
				} else {
					this.sendAndClear()
				}
			}
		}
	}
}

func (this *ZipSendProxyThread) Append(p *pack.LogSinkPack) {
	defer func() {
		if r := recover(); r != nil {
			this.Log.Info("ZipSend Append Recover  ", r)
		}
	}()

	dout := io.NewDataOutputX()
	pack.WritePack(dout, p)
	this.buffer.Write(dout.ToByteArray())
	this.packCount += 1

	if this.firstTime == 0 {
		this.firstTime = p.Time
		if this.buffer.Len() >= this.logsinkMaxBufferSize {
			this.sendAndClear()
		}
	} else {
		if this.buffer.Len() >= this.logsinkMaxBufferSize || p.Time-this.firstTime >= this.logsinkMaxWaitTime {
			this.sendAndClear()
		}
	}
}

func (this *ZipSendProxyThread) sendAndClear() {
	// this.Log.Info("ZipSend sendAndClear ", this.buffer.Len())
	if this.buffer.Len() == 0 {
		return
	}

	p := pack.NewZipPack()
	p.Time = dateutil.SystemNow()
	p.RecordCount = this.packCount
	p.Records = this.buffer.Bytes()

	this.doZip(p)

	if err := this.client.SendFlush(p, true); err != nil {
		this.Log.Errorf("TcpClient SendFlush error %v", err)
	}

	this.buffer.Reset()
	this.firstTime = 0
	this.packCount = 0
}

func (this *ZipSendProxyThread) doZip(p *pack.ZipPack) {
	if p.Status != 0 {
		return
	}
	if len(p.Records) < this.logsinkZipMinSize {
		return
	}

	z := NewDefaultZipMod()
	p.Status = z.ID()
	var err error
	if p.Records, err = z.Compress(p.Records); err != nil {
		this.Log.Error("WA-LOGS-103", "Compress Error ", err)
	}
}

func (this *ZipSendProxyThread) SendDirect(arr []*pack.LogSinkPack) {
	// this.Log.Info("ZipSend SendDirect ", len(arr))
	var buffer bytes.Buffer
	p := pack.NewZipPack()
	p.RecordCount = 0

	for _, it := range arr {
		dout := io.NewDataOutputX()
		pack.WritePack(dout, it)
		buffer.Write(dout.ToByteArray())
		p.RecordCount++

		if buffer.Len() >= this.logsinkMaxBufferSize {
			p.Records = buffer.Bytes()
			p.Time = dateutil.SystemNow()

			this.doZip(p)

			if err := this.client.SendFlush(p, true); err != nil {
				this.Log.Errorf("TcpClient SendFlush error %v", err)
			}

			buffer.Reset()

			//  init after send
			p = pack.NewZipPack()
			p.RecordCount = 0
		}
	}

	if buffer.Len() > 0 {
		p.Records = buffer.Bytes()
		p.Time = dateutil.SystemNow()
		this.doZip(p)

		if err := this.client.SendFlush(p, true); err != nil {
			this.Log.Errorf("TcpClient SendFlush error %v", err)
		}
		buffer.Reset()
	}
}

// implements config.ConfigObserver
func (this *ZipSendProxyThread) ApplyConfig(conf config.Config) {

	queueSize := int(conf.GetInt("logsink_queue_size", 1000))
	if this.logsinkQueueSize != queueSize {
		this.logsinkQueueSize = queueSize
		if this.Queue != nil {
			this.Queue.SetCapacity(queueSize)
		}
	}

	this.logsinkMaxWaitTime = int64(conf.GetInt("max_wait_time", 2000))
	this.logsinkMaxBufferSize = int(conf.GetInt("max_buffer_size", 1024*64))
	this.logsinkZipMinSize = int(conf.GetInt("logsink_zip_min_size", 100))
}
