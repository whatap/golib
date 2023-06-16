package udp

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/whatap/golib/config"
	"github.com/whatap/golib/io"
	"github.com/whatap/golib/lang/pack"
	"github.com/whatap/golib/lang/pack/udp"
	"github.com/whatap/golib/util/dateutil"
	"github.com/whatap/golib/util/queue"
)

const (
	UDP_READ_MAX                    = 64 * 1024
	UDP_PACKET_BUFFER               = 64 * 1024
	UDP_PACKET_BUFFER_CHUNKED_LIMIT = 48 * 1024
	UDP_PACKET_CHANNEL_MAX          = 2048
	UDP_PACKET_FLUSH_TIMEOUT        = 10 * 1000

	UDP_PACKET_HEADER_SIZE = 9
	// typ pos 0
	UDP_PACKET_HEADER_TYPE_POS = 0
	// ver pos 1
	UDP_PACKET_HEADER_VER_POS = 1
	// len pos 5
	UDP_PACKET_HEADER_LEN_POS = 5

	UDP_PACKET_SQL_MAX_SIZE = 32768
)

type UdpClient struct {
	addr       string
	serverAddr *net.UDPAddr
	localAddr  *net.UDPAddr

	udp *net.UDPConn
	wr  *bufio.Writer

	sendCh       chan *UdpData
	sendQ        *queue.RequestQueue
	lastSendTime int64

	lock sync.Mutex
	conf *udpClientConfig

	lockCount     sync.Mutex
	packCount     int
	sendCount     int
	flushErrCount int
	errCount      int
	chanCount     int

	buffer bytes.Buffer
}

type UdpData struct {
	Type  byte
	Ver   int32
	Data  []byte
	Flush bool
}

//
var udpClient *UdpClient
var udpClientLock sync.Mutex

func GetUdpClient(opts ...UdpClientOption) *UdpClient {
	udpClientLock.Lock()
	defer udpClientLock.Unlock()

	if udpClient != nil {
		return udpClient
	}

	udpClient = newUdpClient(opts...)
	udpClient = new(UdpClient)
	udpClient.open()

	udpClient.addr = fmt.Sprintf("%s:%d", udpClient.conf.NetUdpHost, udpClient.conf.NetUdpPort)
	if serverAddr, err := net.ResolveUDPAddr("udp", udpClient.addr); err == nil {
		udpClient.serverAddr = serverAddr
	}
	if localAddr, err := net.ResolveUDPAddr("udp", ":0"); err == nil {
		udpClient.localAddr = localAddr
	}

	go func() {
		for {
			for udpClient.open() == false {
				time.Sleep(3000 * time.Millisecond)
			}
			for udpClient.isOpen() {
				time.Sleep(5000 * time.Millisecond)
			}
		}
	}()

	go udpClient.receive()
	go udpClient.process()

	return udpClient
}

func newUdpClient(opts ...UdpClientOption) *UdpClient {
	p := new(UdpClient)
	p.conf = newDefaultConfig()

	for _, opt := range opts {
		opt.apply(p.conf)
	}

	if p.conf.ConfigObserver != nil {
		p.conf.ConfigObserver.Add("UdpClient", p)
	}

	if p.conf.ctx == nil {
		p.conf.ctx, p.conf.cancel = context.WithCancel(context.Background())
	} else if p.conf.cancel == nil {
		p.conf.ctx, p.conf.cancel = context.WithCancel(context.Background())
	}

	if p.conf.Timeout == 0 {
		p.conf.Timeout = netTimeout
	}

	p.sendCh = make(chan *UdpData, UDP_PACKET_CHANNEL_MAX)
	p.sendQ = queue.NewRequestQueue(UDP_PACKET_CHANNEL_MAX)

	return p
}

func (this *UdpClient) open() (ret bool) {
	this.lock.Lock()
	defer this.lock.Unlock()

	addr := fmt.Sprintf("%s:%d", this.conf.NetUdpHost, this.conf.NetUdpPort)
	if this.addr != addr {
		this.Close()
		if this.udp != nil {
			this.udp.Close()
			this.udp = nil
		}
		if serverAddr, err := net.ResolveUDPAddr("udp", addr); err == nil {
			this.serverAddr = serverAddr
		}
		this.addr = addr
	}

	if this.isOpen() {
		return true
	}

	conn, err := net.DialTimeout("udp", addr, time.Duration(60000)*time.Millisecond)
	if err != nil {
		this.conf.Log.Error("[WA-GO-01001]", "Error. UDP Connect, ", addr)
		return false
	}

	if udpConn, ok := conn.(*net.UDPConn); ok {
		this.udp = udpConn
		this.wr = bufio.NewWriterSize(this.udp, UDP_PACKET_BUFFER)
		this.udp = udpConn
		return false
	}

	this.addr = addr
	this.conf.Log.Debug("[WA-GO-01000] Connected to ", conn.RemoteAddr().(*net.UDPAddr))
	return true
}

func (this *UdpClient) isOpen() bool {
	return this.udp != nil && this.wr != nil
}

func (this *UdpClient) GetLocalAddr() net.Addr {
	return this.udp.LocalAddr()
}

func (this *UdpClient) Send(p udp.UdpPack) {
	dout := udp.WritePack(io.NewDataOutputX(), p)
	this.sendByBuffer(&UdpData{p.GetPackType(), p.GetVersion(), dout.ToByteArray(), p.IsFlush()})
	udp.ClosePack(p)
}

func (this *UdpClient) SendRelay(p pack.Pack, flush bool) {
	pb := pack.ToBytesPack(p)

	relayPack := udp.CreatePack(udp.RELAY_PACK, udp.UDP_PACK_VERSION)
	if rp, ok := relayPack.(*udp.UdpRelayPack); ok {
		rp.RelayType = p.GetPackType()
		rp.Data = pb
		rp.Flush = flush
		rp.Time = dateutil.SystemNow()
		b := udp.ToBytesPack(rp)
		this.sendByBuffer(&UdpData{rp.GetPackType(), rp.GetVersion(), b, flush})
	}
	udp.ClosePack(relayPack)
}

func (this *UdpClient) process() {
	for {
		sendData := <-this.sendCh
		this.lastSendTime = dateutil.Now()
		if _, err := this.sendUDP(sendData); err != nil {
			log.Println("[WA-GO-01005]", "Error. send by udp ", err)
			// DEBUG add send count, err
			this.AddCount(0, 0, 1, 0, 1, false)
		} else {
			// DEBUG add send count
			this.AddCount(0, 0, 1, 0, 0, false)
		}
	}
}
func (this *UdpClient) processRemain() {
	for {
		if !this.isOpen() {
			continue
		}
		select {
		case <-time.After(1 * time.Second):
			// 시간 비교하여 발송.
			if this.buffer.Len() > 0 && dateutil.SystemNow()-this.lastSendTime > UDP_PACKET_FLUSH_TIMEOUT {
				//fmt.Println(">>>>", "timeout flush")
				this.sendBuffer()
			}
		}
	}
}

func (this *UdpClient) processMon() {
	for {
		select {
		case <-time.After(10 * time.Second):
			//reset
			fmt.Printf(">>>>Udp t=%d, pack=%d, chan=%d, send=%d, ferr=%d, err=%d, chanlen=%d\n", this.lastSendTime, this.packCount, this.chanCount, this.sendCount, this.flushErrCount, this.errCount, len(this.sendCh))
			this.AddCount(0, 0, 0, 0, 0, true)
		}
	}
}

func (this *UdpClient) sendByBuffer(sendData *UdpData) {
	this.lock.Lock()
	defer func() {
		if r := recover(); r != nil {
			this.conf.Log.Error("[WA-GO-01007-01]", "Recover UdpClient.sendBuffer ", r)

		}
		this.lock.Unlock()
	}()

	// DEBUG add pack count
	this.AddCount(1, 0, 0, 0, 0, false)

	if !this.isOpen() {
		this.conf.Log.Error("[WA-GO-01008]", "Before a UDP connection is established.")
		return
	}
	if sendData == nil {
		this.conf.Log.Error("[WA-GO-01008]", "Data is nil")
		return
	}

	out := io.NewDataOutputX()
	out.WriteByte(sendData.Type)
	out.WriteInt(sendData.Ver)
	out.WriteIntBytes(sendData.Data)
	sendBytes := out.ToByteArray()

	if this.buffer.Len() > 0 && this.buffer.Len()+len(sendBytes) > UDP_PACKET_BUFFER_CHUNKED_LIMIT {
		// add chanCount
		this.AddCount(0, 1, 0, 0, 0, false)

		if len(this.sendCh) == cap(this.sendCh) {
			fmt.Println("1 sendCh full ", len(this.sendCh))
		}
		data := make([]byte, this.buffer.Len())
		copy(data, this.buffer.Bytes())
		select {
		case this.sendCh <- &UdpData{Data: data}:
		case <-time.After(5 * time.Second):
			// fmt.Println("<<<<", "send to chan timeout 1")
		}
		this.buffer.Reset()
	}
	if _, err := this.buffer.Write(sendBytes); err != nil {
		this.conf.Log.Error("[WA-GO-01010-01]", "Error. Write to buffer len=", len(sendBytes), ", err=", err)
		this.Close()
		// DEBUG add write errCount,
		this.AddCount(0, 0, 0, 0, 1, false)
		return
	}
	// flush == true
	if this.buffer.Len() > 0 && sendData.Flush {
		// add chanCount
		this.AddCount(0, 1, 0, 0, 0, false)

		if len(this.sendCh) == cap(this.sendCh) {
			// fmt.Println("2 sendCh full ", len(this.sendCh))
		}
		data := make([]byte, this.buffer.Len())
		copy(data, this.buffer.Bytes())
		this.buffer.Reset()
		select {
		case this.sendCh <- &UdpData{Data: data}:
		case <-time.After(5 * time.Second):
			// fmt.Println("<<<<", "send to chan timeout 1")
		}

	}
}

func (this *UdpClient) sendBuffer() {
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.buffer.Len() > 0 {
		if len(this.sendCh) == cap(this.sendCh) {
			fmt.Println("3 sendCh full ", len(this.sendCh))
		}
		fmt.Println(">>>> Send to chan 3")
		data := make([]byte, this.buffer.Len())
		copy(data, this.buffer.Bytes())
		this.buffer.Reset()
		select {
		case this.sendCh <- &UdpData{Data: data}:
		case <-time.After(5 * time.Second):
			fmt.Println("<<<<", "send to chan timeout 1")
		}

		// add chanCount
		this.AddCount(0, 1, 0, 0, 0, false)
	}
}

func (this *UdpClient) sendUDP(udpData *UdpData) (int, error) {
	if this.isOpen() == false {
		return 0, fmt.Errorf("udp disconnected")
	}
	return this.udp.Write(udpData.Data)
}

func (this *UdpClient) receive() {
	buff := make([]byte, UDP_PACKET_BUFFER)

	for {
		for this.isOpen() == false {
			time.Sleep(1000 * time.Millisecond)
		}
		if this.udp != nil {
			func() {
				defer func() {
					if r := recover(); r != nil {
						this.conf.Log.Error("[WA-GO-01012]", "Recover UdpClient.receive ", r)
					}
				}()
				if _, _, err := this.udp.ReadFrom(buff); err != nil {
					this.conf.Log.Error("[WA-GO-01013]", "Error. ReadFromUDP ", err)
					this.Close()
					return
				}

				offset := 0
				t := uint8(buff[offset])
				v := io.ToInt(buff[offset+UDP_PACKET_HEADER_VER_POS:offset+UDP_PACKET_HEADER_VER_POS+4], 0)
				l := io.ToInt(buff[offset+UDP_PACKET_HEADER_LEN_POS:offset+UDP_PACKET_HEADER_LEN_POS+4], 0)

				offset += UDP_PACKET_HEADER_SIZE

				tmp := buff[offset : offset+int(l)]
				offset += int(l)
				switch t {
				case udp.CONFIG_INFO:
					p := udp.ToPack(t, v, tmp)
					if p != nil {
						//this.conf.ApplyConfig(p.(*udp.UdpConfigPack).MapData)
					}
				}
			}()
		}
	}
}

func (this *UdpClient) Close() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("[WA-GO-01014]", "Recover Close() ", r)
		}
	}()
	// if this.wr != nil {
	// 	this.wr.Reset(this.udp)
	// }
}
func (this *UdpClient) Shutdown() {
	if this.udp != nil {
		defer func() {
			recover()
		}()
		// block to receive
		close(this.sendCh)
		// send all remaining data
		for sendData := range this.sendCh {
			this.sendUDP(sendData)
		}
		this.Close()
		this.udp.Close()
		this.udp = nil
	}
}

func (this *UdpClient) AddCount(packCount, chanCount, sendCount, flushErrCount, errCount int, reset bool) {
	this.lockCount.Lock()
	defer this.lockCount.Unlock()
	if reset {
		this.packCount = packCount
		this.chanCount = chanCount
		this.sendCount = sendCount
		this.flushErrCount = flushErrCount
		this.errCount = errCount

	} else {
		this.packCount += packCount
		this.chanCount += chanCount
		this.sendCount += sendCount
		this.flushErrCount += flushErrCount
		this.errCount += errCount
	}
}

func UdpShutdown() {
	if udpClient != nil {
		udpClient.Shutdown()
	}
}

// implements common.ConfigObserver
func (this *UdpClient) ApplyConfig(conf config.Config) {

	//this.Pcode = conf.GetLong("pcode")
	host := conf.GetValueDef("net_udp_host", "127.0.0.1")
	port := conf.GetValueDef("net_udp_port", "6600")

	srv := fmt.Sprintf("%s:%s", host, port)
	if this.conf.Server != srv {
		UdpShutdown()
		this.open()
	}
}
