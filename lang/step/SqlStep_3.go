package step

import (
	"github.com/whatap/golib/io"
)

type SqlStep_3 struct {
	AbstractStep
	Xtype   byte
	Hash    int32
	Elapsed int32
	Error   int64
	Dbc     int32

	Updated int32
	Crud    byte

	////////////////////
	Opt byte

	P1   []byte
	P2   []byte
	Pcrc byte

	StartCpu int32
	StartMem int32

	Cpu int32
	Mem int32

	Stack []int32
}

func NewSqlStep_3() *SqlStep_3 {
	p := new(SqlStep_3)
	return p
}
func (this *SqlStep_3) GetStepType() byte {
	return STEP_SQL_X
}

func (this *SqlStep_3) Write(out *io.DataOutputX) {
	this.AbstractStep.Write(out)
	out.WriteDecimal(int64(this.Hash))
	out.WriteDecimal(int64(this.Elapsed))
	out.WriteDecimal(this.Error)

	out.WriteByte(this.Xtype)
	out.WriteDecimal(int64(this.Updated))
	out.WriteByte(byte(this.Crud))
	out.WriteDecimal(int64(this.Dbc))

	out.WriteByte(this.Opt)
	if this.IsTrue(1) {
		out.WriteBlob(this.P1)
		out.WriteBlob(this.P2)
		out.WriteByte(this.Pcrc)
	}
	if this.IsTrue(2) {
		out.WriteDecimal(int64(this.StartCpu))
		out.WriteDecimal(int64(this.Cpu))
		out.WriteDecimal(int64(this.StartMem))
		out.WriteDecimal(int64(this.Mem))
	}
	if this.IsTrue(4) {
		out.WriteIntArray(this.Stack)
	}
}

func (this *SqlStep_3) Read(in *io.DataInputX) {
	this.AbstractStep.Read(in)
	this.Hash = int32(in.ReadDecimal())
	this.Elapsed = int32(in.ReadDecimal())
	this.Error = in.ReadDecimal()

	this.Xtype = in.ReadByte()
	this.Updated = int32(in.ReadDecimal())
	this.Crud = in.ReadByte()
	this.Dbc = int32(in.ReadDecimal())

	this.Opt = in.ReadByte()
	if this.IsTrue(1) {
		this.P1 = in.ReadBlob()
		this.P2 = in.ReadBlob()
		this.Pcrc = in.ReadByte()
	}
	if this.IsTrue(2) {
		this.StartCpu = int32(in.ReadDecimal())
		this.Cpu = int32(in.ReadDecimal())
		this.StartMem = int32(in.ReadDecimal())
		this.Mem = int32(in.ReadDecimal())
	}
	if this.IsTrue(4) {
		this.Stack = in.ReadIntArray()
	}
}

func (this *SqlStep_3) IsTrue(flag byte) bool {
	return (this.Opt & flag) != 0
}
func (this *SqlStep_3) SetTrue(flag byte) {
	this.Opt |= flag
}
func (this *SqlStep_3) GetElapsed() int32 {
	return this.Elapsed
}
