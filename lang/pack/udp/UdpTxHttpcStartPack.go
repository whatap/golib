package udp

import (
	"fmt"

	"github.com/whatap/golib/io"
	"github.com/whatap/golib/util/hash"
)

type UdpTxHttpcStartPack struct {
	AbstractPack
	Url     string
	UrlHash int32
}

func NewUdpTxHttpcStartPack() *UdpTxHttpcStartPack {
	p := new(UdpTxHttpcStartPack)
	p.Ver = UDP_PACK_VERSION
	p.AbstractPack.Flush = false
	return p
}
func NewUdpTxHttpcStartPackVer(ver int32) *UdpTxHttpcStartPack {
	p := new(UdpTxHttpcStartPack)
	p.Ver = ver
	p.AbstractPack.Flush = false
	return p
}
func (this *UdpTxHttpcStartPack) GetPackType() uint8 {
	return TX_HTTPC_START
}

func (this *UdpTxHttpcStartPack) ToString() string {
	return fmt.Sprint(this.AbstractPack.ToString(), ",Url=", this.Url)
}

func (this *UdpTxHttpcStartPack) Clear() {
	this.AbstractPack.Clear()
	this.AbstractPack.Flush = false
	this.Url = ""
	this.UrlHash = 0
}

func (this *UdpTxHttpcStartPack) Write(dout *io.DataOutputX) {
	this.AbstractPack.Write(dout)
	dout.WriteTextShortLength(this.Url)
}

func (this *UdpTxHttpcStartPack) Read(din *io.DataInputX) {
	this.AbstractPack.Read(din)
	this.Url = din.ReadTextShortLength()
}

func (this *UdpTxHttpcStartPack) Process() {
	this.UrlHash = hash.HashStr(this.Url)
}
