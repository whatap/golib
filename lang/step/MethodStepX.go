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

	Component string
	Exception string
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
	if this.Component == "" {
		out.WriteByte(0)

		out.WriteDecimal(int64(this.Hash))
		out.WriteDecimal(int64(this.Elapsed))
		out.WriteDecimal(int64(this.StartCpu))
		out.WriteDecimal(int64(this.StartMem))
		out.WriteIntArray(this.Stack)
	} else {
		out.WriteByte(1)

		out.WriteDecimal(int64(this.Hash))
		out.WriteDecimal(int64(this.Elapsed))
		out.WriteDecimal(int64(this.StartCpu))
		out.WriteDecimal(int64(this.StartMem))
		out.WriteIntArray(this.Stack)
		out.WriteText(this.Component)
		out.WriteText(this.Exception)
	}
}

func (this *MethodStepX) Read(in *io.DataInputX) {
	this.AbstractStep.Read(in)
	ver := in.ReadByte()

	this.Hash = int32(in.ReadDecimal())
	this.Elapsed = int32(in.ReadDecimal())
	this.StartCpu = int32(in.ReadDecimal())
	this.StartMem = int32(in.ReadDecimal())
	this.Stack = in.ReadIntArray()
	if ver == 0 {
		return
	}
	this.Component = in.ReadText()
	this.Exception = in.ReadText()
}

func (this *MethodStepX) GetElapsed() int32 {
	return this.Elapsed
}
