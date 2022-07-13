package step

import (
	"github.com/whatap/golib/io"
)

type SecureMsgStep struct {
	AbstractStep
	Hash  int32
	Opt   byte
	Crc   byte
	Value []byte
}

func NewSecureMsgStep() *SecureMsgStep {
	p := new(SecureMsgStep)
	return p
}
func (this *SecureMsgStep) GetStepType() byte {
	return STEP_SECURE_MESSAGE
}

func (this *SecureMsgStep) Write(out *io.DataOutputX) {
	this.AbstractStep.Write(out)
	out.WriteDecimal(int64(this.Hash))
	out.WriteByte(this.Opt)
	out.WriteByte(this.Crc)
	out.WriteBlob(this.Value)
}

func (this *SecureMsgStep) Read(in *io.DataInputX) {
	this.AbstractStep.Read(in)
	this.Hash = int32(in.ReadDecimal())
	this.Opt = in.ReadByte()
	this.Crc = in.ReadByte()
	this.Value = in.ReadBlob()
}

func (this *SecureMsgStep) GetElapsed() int32 {
	return 0 //this.Elapsed
}
