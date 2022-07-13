package step

import (
	"github.com/whatap/golib/io"
)

type MessageStep struct {
	AbstractStep
	Hash  int32
	Time  int32
	Value int32
	Desc  string
}

func NewMessageStep() *MessageStep {
	p := new(MessageStep)
	return p
}
func (this *MessageStep) GetStepType() byte {
	return STEP_MESSAGE
}

func (this *MessageStep) Write(out *io.DataOutputX) {
	this.AbstractStep.Write(out)
	out.WriteDecimal(int64(this.Hash))
	out.WriteDecimal(int64(this.Time))
	out.WriteDecimal(int64(this.Value))
	out.WriteText(this.Desc)
}

func (this *MessageStep) Read(in *io.DataInputX) {
	this.AbstractStep.Read(in)
	this.Hash = int32(in.ReadDecimal())
	this.Time = int32(in.ReadDecimal())
	this.Value = int32(in.ReadDecimal())
	this.Desc = in.ReadText()
}

func (this *MessageStep) GetElapsed() int32 {
	return 0 //this.Elapsed
}
