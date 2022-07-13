package pack

import (
	"github.com/whatap/golib/io"
)

//##################################################
// SYSBASE PACK
//##################################################
type SMPingPack struct {
	AbstractPack
	IP   int32
	OS   int16
	Core int16
}

func NewSMPingPack() *SMPingPack {
	p := new(SMPingPack)
	return p
}

func (this *SMPingPack) GetPackType() int16 {
	return PACK_SM_PING
}
func (this *SMPingPack) Write(doutx *io.DataOutputX) {
	dout := io.NewDataOutputX()
	this.AbstractPack.Write(dout)
	dout.WriteInt(this.IP)
	dout.WriteShort(this.OS)
	dout.WriteShort(this.Core)

	doutx.WriteBlob(dout.ToByteArray())
}
func (this *SMPingPack) Read(dinx *io.DataInputX) {
	din := io.NewDataInputX(dinx.ReadBlob())
	this.AbstractPack.Read(din)
	this.IP = din.ReadInt()
	this.OS = din.ReadShort()
	this.Core = din.ReadShort()
}
