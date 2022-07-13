//github.com/whatap/golib/net/udp
package udp

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/whatap/golib/io"
	"github.com/whatap/golib/lang/pack"
	"github.com/whatap/golib/lang/pack/udp"
	"github.com/whatap/golib/util/dateutil"
	"github.com/whatap/go-api/config"
)

const (
	UDP_READ_MAX                    = 64 * 1024
	UDP_PACKET_BUFFER               = 64 * 1024
	UDP_PACKET_BUFFER_CHUNKED_LIMIT = 48 * 1024
	UDP_PACKET_CHANNEL_MAX          = 255
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
	addr string

	udp net.Conn
	wr  *bufio.Writer

	sendCh       chan *UdpData
	lastSendTime int64

	lock sync.Mutex
	conf *config.Config
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

func GetUdpClient() *UdpClient {
	udpClientLock.Lock()
	defer udpClientLock.Unlock()

	if udpClient != nil {
		return udpClient
	}
	udpClient = new(UdpClient)
	udpClient.conf = config.GetConfig()
	udpClient.open()
	udpClient.sendCh = make(chan *UdpData, UDP_PACKET_CHANNEL_MAX)
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

func (this *UdpClient) open() (ret bool) {
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.isOpen() {
		return true
	}
	addr := fmt.Sprintf("%s:%d", this.conf.NetUdpHost, this.conf.NetUdpPort)

	if this.addr != addr {
		this.Close()
	}

	con, err := net.DialTimeout("udp", addr, time.Duration(60000)*time.Millisecond)
	if err != nil {
		if this.conf.Debug {
			log.Println("[WA-GO-01001]", "Error. UDP Connect, ", addr)
		}
		this.Close()
		return false
	}
	if this.conf.Debug {
		log.Println("[WA-GO-01002]", "UDP Connected: ", addr)
	}
	this.addr = addr
	this.wr = bufio.NewWriterSize(con, UDP_PACKET_BUFFER)
	this.udp = con

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
	this.sendQueue(p.GetPackType(), p.GetVersion(), dout.ToByteArray(), p.IsFlush())
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
		this.sendQueue(rp.GetPackType(), rp.GetVersion(), b, flush)
	}
	udp.ClosePack(relayPack)
}

func (this *UdpClient) sendQueue(t uint8, ver int32, b []byte, flush bool) bool {
	this.lock.Lock()
	defer func() {
		if r := recover(); r != nil {
			if this.conf.Debug {
				log.Println("[WA-GO-01003]", "Recover UdpClient.sendQueue ", r)
			}
		}
		this.lock.Unlock()
	}()
	if this.isOpen() {
		if len(b) >= UDP_PACKET_BUFFER {
			log.Println("[WA-GO-01003-01]", "big data ")
			return false
		}
		buff := make([]byte, len(b))
		copy(buff, b)
		this.sendCh <- &UdpData{t, ver, buff, flush}
		return true
	} else {
		if this.conf.Debug {
			log.Println("[WA-GO-01004]", "Before a UDP connection is established.")
		}
		return false
	}
}

func (this *UdpClient) process() {
	for {
		func() {
			defer func() {
				if r := recover(); r != nil {
					if this.conf.Debug {
						log.Println("[WA-GO-01005]", "Recover UdpClient.process ", r)
					}
				}
			}()
			select {
			case sendData := <-this.sendCh:
				this.send(sendData)
			default:
				if !this.isOpen() {
					return
				}
				time.Sleep(1 * time.Second)
				// 시간 비교하여 발송.
				if this.wr != nil {
					if this.wr.Buffered() > 0 && dateutil.SystemNow()-this.lastSendTime > UDP_PACKET_FLUSH_TIMEOUT {
						this.lastSendTime = dateutil.Now()
						if err := this.wr.Flush(); err != nil {
							if this.conf.Debug {
								log.Println("[WA-GO-01006]", "Error. Flush bufferd writer ", err)
							}
							this.Close()
							return
						}
					}
				}

			}
		}()
	}
}
func (this *UdpClient) send(sendData *UdpData) {
	defer func() {
		if r := recover(); r != nil {
			if this.conf.Debug {
				log.Println("[WA-GO-01007]", "Recover UdpClient.send ", r)
			}
		}
	}()
	if !this.isOpen() {
		if this.conf.Debug {
			log.Println("[WA-GO-01008]", "Before a UDP connection is established.")
		}
		return
	}
	if sendData == nil {
		if this.conf.Debug {
			log.Println("[WA-GO-01008]", "Data is nil")
		}
		return
	}

	out := io.NewDataOutputX()
	out.WriteByte(sendData.Type)
	out.WriteInt(sendData.Ver)
	out.WriteIntBytes(sendData.Data)
	sendBytes := out.ToByteArray()
	if this.wr.Buffered() > 0 && this.wr.Buffered()+len(sendBytes) > UDP_PACKET_BUFFER_CHUNKED_LIMIT {
		this.lastSendTime = dateutil.Now()
		if err := this.wr.Flush(); err != nil {
			if this.conf.Debug {
				log.Println("[WA-GO-01009]", "Error. Flush bufferd writer", err)
			}
			this.Close()
			return
		}
	}
	if _, err := this.wr.Write(sendBytes); err != nil {
		if this.conf.Debug {
			log.Println("[WA-GO-01010]", "Error. Write to bufferd writer len=", len(sendBytes), ", err=", err)
		}
		this.Close()
		return
	}
	// flush == true
	if this.wr.Buffered() > 0 && sendData.Flush {
		this.lastSendTime = dateutil.Now()
		if err := this.wr.Flush(); err != nil {
			if this.conf.Debug {
				log.Println("[WA-GO-01011]", "Error. Flush bufferd writer ", err)
			}
			this.Close()
			return
		}
	}
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
						if this.conf.Debug {
							log.Println("[WA-GO-01012]", "Recover UdpClient.receive ", r)
						}
					}
				}()
				udpNet := this.udp.(*net.UDPConn)
				if _, _, err := udpNet.ReadFromUDP(buff); err != nil {
					if this.conf.Debug {
						log.Println("[WA-GO-01013]", "Error. ReadFromUDP ", err)
					}
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
						this.conf.ApplyConfig(p.(*udp.UdpConfigPack).MapData)
					}
				}
			}()
		}
	}
}
func (this *UdpClient) Close() {
	defer func() {
		recover()
	}()
	if this.udp != nil {
		this.udp.Close()
	}
	this.udp = nil
	this.wr = nil
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
			this.send(sendData)
		}
		this.Close()
	}
}

func UdpShutdown() {
	if udpClient != nil {
		udpClient.Shutdown()
	}
}
