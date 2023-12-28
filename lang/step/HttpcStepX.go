package step

import (
	"github.com/whatap/golib/io"
)

const (
	HTTPC_STEP_DEFAULT_VERSION = 2
)

type HttpcStepX struct {
	AbstractStep

	Version byte

	Url  int32
	Host int32
	Port int32

	Elapsed int32
	Error   int64

	Status int32

	StartCpu int32
	StartMem int64
	Stack    []int32

	StepId    int64
	Driver    string
	OriginUrl string
	Param     string
}

func NewHttpcStepX() *HttpcStepX {
	p := new(HttpcStepX)
	p.Version = HTTPC_STEP_DEFAULT_VERSION
	return p
}
func NewHttpcStepXVersion(ver byte) *HttpcStepX {
	p := new(HttpcStepX)
	p.Version = ver
	return p
}

func (this *HttpcStepX) GetStepType() byte {
	return STEP_HTTPCALL_X
}

func (this *HttpcStepX) Write(out *io.DataOutputX) {
	this.AbstractStep.Write(out)

	out.WriteByte(this.Version)
	out.WriteDecimal(int64(this.Url))
	out.WriteDecimal(int64(this.Elapsed))
	out.WriteDecimal(this.Error)
	out.WriteDecimal(int64(this.Host))
	out.WriteDecimal(int64(this.Port))

	out.WriteDecimal(int64(this.Status))
	out.WriteDecimal(int64(this.StartCpu))
	out.WriteDecimal(int64(this.StartMem))
	out.WriteIntArray(this.Stack)
	switch this.Version {
	case 1:
		out.WriteDecimal(0)
	case 2:
		out.WriteDecimal(this.StepId)
		out.WriteText(this.Driver)
		out.WriteText(this.OriginUrl)
		out.WriteText(this.Param)
	}

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

	switch ver {
	case 1:
		in.ReadDecimal()
	case 2:
		this.StepId = in.ReadDecimal()
		this.Driver = in.ReadText()
		this.OriginUrl = in.ReadText()
		this.Param = in.ReadText()
	}
}

func (this *HttpcStepX) GetElapsed() int32 {
	return this.Elapsed
}
