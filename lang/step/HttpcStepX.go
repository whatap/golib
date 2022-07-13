package step

import (
	"github.com/whatap/golib/io"
)

type HttpcStepX struct {
	AbstractStep

	Url  int32
	Host int32
	Port int32

	Elapsed int32
	Error   int64

	Status int32

	StartCpu int32
	StartMem int64
	Stack    []int32

	Callee int64
}

func NewHttpcStepX() *HttpcStepX {
	p := new(HttpcStepX)
	return p
}
func (this *HttpcStepX) GetStepType() byte {
	return STEP_HTTPCALL_X
}

func (this *HttpcStepX) Write(out *io.DataOutputX) {
	this.AbstractStep.Write(out)

	out.WriteByte(1)
	out.WriteDecimal(int64(this.Url))
	out.WriteDecimal(int64(this.Elapsed))
	out.WriteDecimal(this.Error)
	out.WriteDecimal(int64(this.Host))
	out.WriteDecimal(int64(this.Port))

	out.WriteDecimal(int64(this.Status))
	out.WriteDecimal(int64(this.StartCpu))
	out.WriteDecimal(int64(this.StartMem))
	out.WriteIntArray(this.Stack)
	out.WriteDecimal(this.Callee)
}

func (this *HttpcStepX) Read(in *io.DataInputX) {
	this.AbstractStep.Read(in)
	ver := in.ReadByte()

	this.Url = int32(in.ReadDecimal())
	this.Elapsed = int32(in.ReadDecimal())
	this.Error = in.ReadDecimal()
	this.Host = int32(in.ReadDecimal())
	this.Port = int32(in.ReadDecimal())
	this.Status = int32(in.ReadDecimal())

	this.StartCpu = int32(in.ReadDecimal())
	this.StartMem = int64(in.ReadDecimal())
	this.Stack = in.ReadIntArray()

	if ver > 0 {
		this.Callee = in.ReadDecimal()
	}
}

func (this *HttpcStepX) GetElapsed() int32 {
	return this.Elapsed
}
