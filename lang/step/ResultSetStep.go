package step

import (
	"github.com/whatap/golib/io"
)

type ResultSetStep struct {
	AbstractStep
	Dbc     int32
	SqlHash int32
	Elapsed int32
	Fetch   int32
}

func NewResultSetStep() *ResultSetStep {
	p := new(ResultSetStep)
	return p
}
func (this *ResultSetStep) GetStepType() byte {
	return STEP_RESULTSET
}

func (this *ResultSetStep) Write(out *io.DataOutputX) {
	this.AbstractStep.Write(out)
	out.WriteDecimal(int64(this.Dbc))
	out.WriteDecimal(int64(this.SqlHash))
	out.WriteDecimal(int64(this.Elapsed))
	out.WriteDecimal(int64(this.Fetch))
}

func (this *ResultSetStep) Read(in *io.DataInputX) {
	this.AbstractStep.Read(in)
	this.Dbc = int32(in.ReadDecimal())
	this.SqlHash = int32(in.ReadDecimal())
	this.Elapsed = int32(in.ReadDecimal())
	this.Fetch = int32(in.ReadDecimal())
}

func (this *ResultSetStep) GetElapsed() int32 {
	return this.Elapsed
}
