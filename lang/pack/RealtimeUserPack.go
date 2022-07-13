package pack

import (
	"github.com/whatap/golib/io"
)

type RealtimeUserPack struct {
	AbstractPack
	Logbits []byte
}

func NewRealtimeUserPack() *RealtimeUserPack {
	p := new(RealtimeUserPack)
	return p
}

func (this *RealtimeUserPack) GetPackType() int16 {
	return PACK_REALTIME_USER
}

func (this *RealtimeUserPack) Write(dout *io.DataOutputX) {
	this.AbstractPack.Write(dout)

	dout.WriteBlob(this.Logbits)

}
func (this *RealtimeUserPack) Read(din *io.DataInputX) {
	this.AbstractPack.Read(din)

	this.Logbits = din.ReadBlob()
}
