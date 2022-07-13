package pack

import (
	"github.com/whatap/golib/io"
	"github.com/whatap/golib/lang/service"
	"github.com/whatap/golib/lang/step"
)

type ProfilePack struct {
	AbstractPack
	// DEBUG TxRecord
	//Transaction service.Service
	Transaction *service.TxRecord
	Steps       []byte
}

func NewProfilePack() *ProfilePack {
	p := new(ProfilePack)
	return p
}

func (this *ProfilePack) GetPackType() int16 {
	return PACK_PROFILE
}

func (this *ProfilePack) Write(dout *io.DataOutputX) {
	this.AbstractPack.Write(dout)

	// DEBUG TxRecord
	//service.ToBytes(this.Transaction, dout)
	this.Transaction.Write(dout)
	dout.WriteBlob(this.Steps)

}
func (this *ProfilePack) Read(din *io.DataInputX) {
	this.AbstractPack.Read(din)

	service.ToObject(din)
	this.Steps = din.ReadBlob()
}
func (this *ProfilePack) SetProfile(steps []step.Step) {
	this.Steps = step.ToBytesStep(steps)
}
