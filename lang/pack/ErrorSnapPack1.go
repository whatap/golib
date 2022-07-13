package pack

import (
	"github.com/whatap/golib/io"
	"github.com/whatap/golib/lang/step"
)

type ErrorSnapPack1 struct {
	AbstractPack
	Seq     int64
	Profile []byte
	Stack   []byte

	AppendType byte
	AppendHash int32
}

func NewErrorSnapPack1() *ErrorSnapPack1 {
	p := new(ErrorSnapPack1)
	return p
}

func (this *ErrorSnapPack1) GetPackType() int16 {
	return PACK_ERROR_SNAP_1
}

func (this *ErrorSnapPack1) Write(dout *io.DataOutputX) {
	this.AbstractPack.Write(dout)

	dout.WriteLong(this.Seq)
	dout.WriteBlob(this.Profile)
	dout.WriteBlob(this.Stack)
	dout.WriteByte(this.AppendType)
	dout.WriteDecimal(int64(this.AppendHash))

}
func (this *ErrorSnapPack1) Read(din *io.DataInputX) {
	this.AbstractPack.Read(din)

	this.Seq = din.ReadLong()
	this.Profile = din.ReadBlob()
	this.Stack = din.ReadBlob()
	this.AppendType = din.ReadByte()
	this.AppendHash = int32(din.ReadDecimal())
}

func (this *ErrorSnapPack1) SetProfile(steps []step.Step) {
	this.Profile = step.ToBytesStep(steps)
}
func (this *ErrorSnapPack1) SetStack(callstack []int32) {
	dout := io.NewDataOutputX()
	dout.WriteIntArray(callstack)
	this.Stack = dout.ToByteArray()
}
