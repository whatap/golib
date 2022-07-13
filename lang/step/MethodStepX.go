package step

import (
	"github.com/whatap/golib/io"
)

type MethodStepX struct {
	AbstractStep
	Hash    int32
	Elapsed int32

	StartCpu int32
	StartMem int32

	Stack []int32
}

func NewMethodStepX() *MethodStepX {
	p := new(MethodStepX)
	return p
}
func (this *MethodStepX) GetStepType() byte {
	return STEP_METHOD_X
}

func (this *MethodStepX) Write(out *io.DataOutputX) {
	this.AbstractStep.Write(out)
	out.WriteByte(0)

	out.WriteDecimal(int64(this.Hash))
	out.WriteDecimal(int64(this.Elapsed))
	out.WriteDecimal(int64(this.StartCpu))
	out.WriteDecimal(int64(this.StartMem))
	out.WriteIntArray(this.Stack)
}

func (this *MethodStepX) Read(in *io.DataInputX) {
	this.AbstractStep.Read(in)
	in.ReadByte()

	this.Hash = int32(in.ReadDecimal())
	this.Elapsed = int32(in.ReadDecimal())
	this.StartCpu = int32(in.ReadDecimal())
	this.StartMem = int32(in.ReadDecimal())
	this.Stack = in.ReadIntArray()
}

func (this *MethodStepX) GetElapsed() int32 {
	return this.Elapsed
}
