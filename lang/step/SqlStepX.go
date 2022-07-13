package step

import (
	"github.com/whatap/golib/io"
)

type SqlStepX struct {
	AbstractStep
	Xtype   byte
	Hash    int32
	Elapsed int32
	Error   int64
	Dbc     int32

	P1   []byte
	P2   []byte
	Pcrc byte

	StartCpu int32
	StartMem int64

	Stack []int32
}

func NewSqlStepX() *SqlStepX {
	p := new(SqlStepX)
	return p
}
func (this *SqlStepX) GetStepType() byte {
	return STEP_SQL_X
}

func (this *SqlStepX) GetElapsed() int32 {
	return this.Elapsed
}

func (this *SqlStepX) Write(out *io.DataOutputX) {
	this.AbstractStep.Write(out)
	out.WriteByte(0)
	out.WriteDecimal(int64(this.Hash))
	out.WriteDecimal(int64(this.Elapsed))
	out.WriteDecimal(this.Error)

	out.WriteByte(this.Xtype)
	out.WriteDecimal(int64(this.Dbc))

	out.WriteBlob(this.P1)
	out.WriteBlob(this.P2)
	out.WriteByte(this.Pcrc)
	out.WriteDecimal(int64(this.StartCpu))
	out.WriteDecimal(int64(this.StartMem))
	out.WriteIntArray(this.Stack)
}

func (this *SqlStepX) Read(in *io.DataInputX) {
	this.AbstractStep.Read(in)
	in.ReadByte()
	this.Hash = int32(in.ReadDecimal())
	this.Elapsed = int32(in.ReadDecimal())
	this.Error = in.ReadDecimal()

	this.Xtype = in.ReadByte()
	this.Dbc = int32(in.ReadDecimal())

	this.P1 = in.ReadBlob()
	this.P2 = in.ReadBlob()
	this.Pcrc = in.ReadByte()

	this.StartCpu = int32(in.ReadDecimal())
	this.StartMem = int64(in.ReadDecimal())

	this.Stack = in.ReadIntArray()
}
