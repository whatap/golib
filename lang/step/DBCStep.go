package step

import (
	"github.com/whatap/golib/io"
)

type DBCStep struct {
	AbstractStep
	Hash    int32
	Elapsed int32
	Error   int32
}

func NewDBCStep() *DBCStep {
	p := new(DBCStep)
	return p
}
func (this *DBCStep) GetStepType() byte {
	return STEP_DBC
}

func (this *DBCStep) Write(out *io.DataOutputX) {
	this.AbstractStep.Write(out)
	out.WriteDecimal(int64(this.Hash))
	out.WriteDecimal(int64(this.Elapsed))
	out.WriteDecimal(int64(this.Error))
}

func (this *DBCStep) Read(in *io.DataInputX) {
	this.AbstractStep.Read(in)
	this.Hash = int32(in.ReadDecimal())
	this.Elapsed = int32(in.ReadDecimal())
	this.Error = int32(in.ReadDecimal())
}

func (this *DBCStep) GetElapsed() int32 {
	return this.Elapsed
}
