package step

import (
	"github.com/whatap/golib/io"
)

type SocketStep struct {
	AbstractStep
	IpAddr  []byte
	Port    int32
	Elapsed int32
	Error   int64
}

func NewSocketStep() *SocketStep {
	p := new(SocketStep)
	return p
}
func (this *SocketStep) GetStepType() byte {
	return STEP_SOCKET
}

func (this *SocketStep) Write(out *io.DataOutputX) {
	this.AbstractStep.Write(out)
	out.WriteBlob(this.IpAddr)
	out.WriteDecimal(int64(this.Port))
	out.WriteDecimal(int64(this.Elapsed))
	out.WriteDecimal(this.Error)
}

func (this *SocketStep) Read(in *io.DataInputX) {
	this.AbstractStep.Read(in)
	this.IpAddr = in.ReadBlob()
	this.Port = int32(in.ReadDecimal())
	this.Elapsed = int32(in.ReadDecimal())
	this.Error = in.ReadDecimal()
}

func (this *SocketStep) GetElapsed() int32 {
	return this.Elapsed
}
