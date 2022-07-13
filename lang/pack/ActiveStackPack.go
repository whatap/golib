package pack

import (
	"github.com/whatap/golib/io"
)

type ActiveStackPack struct {
	AbstractPack
	Seq           int64
	ProfileSeq    int64
	Service       int32
	CallStack     []int32
	CallStackHash int32

	Elapsed int32
}

func NewActiveStackPack() *ActiveStackPack {
	p := new(ActiveStackPack)
	return p
}

func (this *ActiveStackPack) GetPackType() int16 {
	return PACK_ACTIVESTACK_1
}

func (this *ActiveStackPack) Write(dout *io.DataOutputX) {
	this.AbstractPack.Write(dout)

	dout.WriteByte(1)
	dout.WriteLong(this.Seq)
	dout.WriteLong(this.ProfileSeq)
	dout.WriteInt(this.Service)
	dout.WriteInt(this.CallStackHash)
	dout.WriteIntArray(this.CallStack)

	dout.WriteDecimal(int64(this.Elapsed))

}
func (this *ActiveStackPack) Read(din *io.DataInputX) {
	this.AbstractPack.Read(din)

	var ver = din.ReadByte()
	this.Seq = din.ReadLong()
	this.ProfileSeq = din.ReadLong()
	this.Service = din.ReadInt()
	this.CallStackHash = din.ReadInt()
	this.CallStack = din.ReadIntArray()
	if ver > 0 {
		this.Elapsed = int32(din.ReadDecimal())
	}
}
