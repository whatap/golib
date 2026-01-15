package udp

import (
	"fmt"
	"github.com/whatap/golib/io"
	"github.com/whatap/golib/util/stringutil"
)

type UdpTxWebSocketPack struct {
	AbstractPack
	IpAddr  string
	Port    int32
	Elapsed int32
	Error   int64
}

func NewUdpTxWebSocketPack() *UdpTxWebSocketPack {
	p := new(UdpTxWebSocketPack)
	p.Ver = UDP_PACK_VERSION
	p.AbstractPack.Flush = false
	return p
}

func NewUdpTxWebSocketPackVer(ver int32) *UdpTxWebSocketPack {
	p := new(UdpTxWebSocketPack)
	p.Ver = ver
	p.AbstractPack.Flush = false
	return p
}

func (this *UdpTxWebSocketPack) GetPackType() uint8 {
	return TX_WEB_SOCKET
}

func (this *UdpTxWebSocketPack) ToString() string {
	return fmt.Sprint(this.AbstractPack.ToString(), ",ip=", this.IpAddr, ",port=", this.Port, ",elapsed=", this.Elapsed,",error=", this.Error)
}

func (this *UdpTxWebSocketPack) Clear() {
	this.AbstractPack.Clear()
	this.AbstractPack.Flush = false

	this.IpAddr = ""
	this.Port = 0
	this.Elapsed = 0
	this.Error = 0
}

func (this *UdpTxWebSocketPack) Write(dout *io.DataOutputX) {
	this.AbstractPack.Write(dout)
	dout.WriteTextShortLength(this.IpAddr)
	dout.WriteTextShortLength(stringutil.ParseStringZeroToEmpty(int64(this.Port)))
	dout.WriteTextShortLength(stringutil.ParseStringZeroToEmpty(int64(this.Elapsed)))
	dout.WriteTextShortLength(stringutil.ParseStringZeroToEmpty(int64(this.Error)))


	if this.Ver > 50000 {
		// Golang
	} else if this.Ver > 40000 {
		// Batch
	} else if this.Ver > 30000 {
		// Dotnet
	} else if this.Ver > 20000 {
		// Python
	} else {
		// PHP
	}
}

func (this *UdpTxWebSocketPack) Read(din *io.DataInputX) {
	this.AbstractPack.Read(din)
	this.IpAddr = din.ReadTextShortLength()
	this.Port = stringutil.ParseInt32(din.ReadTextShortLength())
	this.Elapsed = stringutil.ParseInt32(din.ReadTextShortLength())
	this.Error = stringutil.ParseInt64(din.ReadTextShortLength())

	if this.Ver > 50000 {
		// Golang
	} else if this.Ver > 40000 {
		// Batch
	} else if this.Ver > 30000 {
		// Dotnet
	} else if this.Ver > 20000 {
		// Python
	} else {
		// PHP
	}
}

func (this *UdpTxWebSocketPack) Process() {
}
