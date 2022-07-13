package pack

import (
	"github.com/whatap/golib/io"
)

type CompositePack struct {
	AbstractPack
	pack []Pack
}

func NewCompositePack() *CompositePack {
	p := new(CompositePack)
	return p
}

func (this *CompositePack) GetPackType() int16 {
	return PACK_COMPOSITE
}

func (this *CompositePack) Write(dout *io.DataOutputX) {
	this.AbstractPack.Write(dout)

	sz := len(this.pack)
	dout.WriteShort(int16(sz))

	for i := 0; i < sz; i++ {
		WritePack(dout, this.pack[i])
	}

}
func (this *CompositePack) Read(din *io.DataInputX) {
	this.AbstractPack.Read(din)

	sz := int(din.ReadShort())
	this.pack = make([]Pack, sz)
	for i := 0; i < sz; i++ {
		this.pack[i] = ReadPack(din)
	}
}
