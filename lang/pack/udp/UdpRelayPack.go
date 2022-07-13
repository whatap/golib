package udp

import (
	"fmt"

	"github.com/whatap/golib/io"
)

type UdpRelayPack struct {
	AbstractPack
	RelayType int16
	Len       int32
	Data      []byte
}

func NewUdpRelayPack() *UdpRelayPack {
	p := new(UdpRelayPack)
	p.Ver = UDP_PACK_VERSION
	p.AbstractPack.Flush = false
	return p
}
func NewUdpRelayPackVer(ver int32) *UdpRelayPack {
	p := new(UdpRelayPack)
	p.Ver = ver
	p.AbstractPack.Flush = false
	return p
}

func (this *UdpRelayPack) GetPackType() uint8 {
	return RELAY_PACK
}

func (this *UdpRelayPack) ToString() string {
	return fmt.Sprint(this.AbstractPack.ToString(), ",relay Type=", this.RelayType, ",len=", len(this.Data))
}

func (this *UdpRelayPack) Clear() {
	this.AbstractPack.Clear()
	this.AbstractPack.Flush = false
	this.RelayType = 0
	this.Len = 0
}

func (this *UdpRelayPack) Write(dout *io.DataOutputX) {
	dout.WriteBytes(this.Data)
}

func (this *UdpRelayPack) Read(din *io.DataInputX) {
	this.Data = din.ReadBytes(this.Len)
}

func (this *UdpRelayPack) Process() {

}
