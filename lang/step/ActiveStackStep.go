package step

import (
	"github.com/whatap/golib/io"
)

type ActiveStackStep struct {
	AbstractStep
	Seq          int64
	HasCallstack bool
}

func NewActiveStackStep() *ActiveStackStep {
	p := new(ActiveStackStep)
	return p
}
func (this *ActiveStackStep) GetStepType() byte {
	return STEP_ACTIVE_STACK
}

func (this *ActiveStackStep) Write(out *io.DataOutputX) {
	this.AbstractStep.Write(out)
	out.WriteLong(this.Seq)
	out.WriteBool(this.HasCallstack)
}

func (this *ActiveStackStep) Read(in *io.DataInputX) {
	this.AbstractStep.Read(in)
	this.Seq = in.ReadLong()
	this.HasCallstack = in.ReadBool()
}

func (this *ActiveStackStep) GetElapsed() int32 {
	return 0 //this.Elapsed
}
