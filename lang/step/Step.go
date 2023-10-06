package step

import (
	"github.com/whatap/golib/io"
)

const (
	STEP_METHOD       = 1
	STEP_SQL          = 2
	STEP_RESULTSET    = 3
	STEP_SOCKET       = 5
	STEP_ACTIVE_STACK = 6
	STEP_MESSAGE      = 7
	STEP_DBC          = 8
	STEP_SQL_2        = 9
	STEP_POSITION     = 10
	STEP_METHOD_2     = 11

	STEP_METHOD_3       = 12
	STEP_SQL_3          = 13
	STEP_HTTPCALL_3     = 14
	STEP_SECURE_MESSAGE = 15
	STEP_CHILD_THREAD   = 16

	STEP_METHOD_X   = 17
	STEP_SQL_X      = 18
	STEP_HTTPCALL_X = 19 // 14 -> 19
	STEP_REMOTE     = 20
	STEP_DBC_X      = 21
	STEP_MESSAGE_X  = 22

	STEP_RESOURCE    = 50
	STEP_RUM_ERR_MSG = 51
)

type Step interface {
	GetStepType() byte
	Write(out *io.DataOutputX)
	Read(in *io.DataInputX)

	// interface 함수 추가
	GetStartTime() int32
	SetStartTime(t int32)

	SetParent(v int32)
	GetParent() int32

	SetIndex(v int32)
	GetIndex() int32

	SetDrop(v bool)
	GetDrop() bool

	GetElapsed() int32
}

func CreateStep(t byte) Step {
	switch t {
	case STEP_METHOD_X:
		return NewMethodStepX()
	case STEP_SQL_X:
		return NewSqlStepX()
	case STEP_RESULTSET:
		return NewResultSetStep()
	case STEP_SOCKET:
		return NewSocketStep()
	case STEP_HTTPCALL_X:
		return NewHttpcStepX()
	case STEP_ACTIVE_STACK:
		return NewActiveStackStep()
	case STEP_MESSAGE:
		return NewMessageStep()
	case STEP_SECURE_MESSAGE:
		return NewSecureMsgStep()
	case STEP_DBC:
		return NewDBCStep()
	}
	return nil
}

type AbstractStep struct {
	Parent    int32
	Index     int32
	StartTime int32
	Drop      bool
	Opt       byte
}

func (this *AbstractStep) GetStartTime() int32 {
	return this.StartTime
}
func (this *AbstractStep) SetStartTime(t int32) {
	this.StartTime = t
}

func (this *AbstractStep) SetParent(v int32) {
	this.Parent = v
}
func (this *AbstractStep) GetParent() int32 {
	return this.Parent
}
func (this *AbstractStep) SetIndex(v int32) {
	this.Index = v
}
func (this *AbstractStep) GetIndex() int32 {
	return this.Index
}

func (this *AbstractStep) SetDrop(v bool) {
	this.Drop = v
}
func (this *AbstractStep) GetDrop() bool {
	return this.Drop
}

func (this *AbstractStep) IsTrue(flag int) bool {
	return (this.Opt & byte(flag)) != 0
}
func (this *AbstractStep) SetTrue(flag int) {
	this.Opt |= byte(flag)
}

func (this *AbstractStep) Write(dout *io.DataOutputX) {
	dout.WriteDecimal(int64(this.Parent))
	dout.WriteDecimal(int64(this.Index))
	dout.WriteDecimal(int64(this.StartTime))
}

func (this *AbstractStep) Read(din *io.DataInputX) {
	this.Parent = int32(din.ReadDecimal())
	this.Index = int32(din.ReadDecimal())
	this.StartTime = int32(din.ReadDecimal())
}

func WriteStep(out *io.DataOutputX, p Step) *io.DataOutputX {
	out.WriteByte(p.GetStepType())
	p.Write(out)
	return out
}
func ReadStep(in *io.DataInputX) Step {
	t := in.ReadByte()
	v := CreateStep(t)
	v.Read(in)
	return v
}

func ToBytesStep(p []Step) []byte {
	if p == nil {
		return nil
	}
	dout := io.NewDataOutputX()
	sz := len(p)
	for i := 0; i < sz; i++ {
		WriteStep(dout, p[i])
	}
	return dout.ToByteArray()
}
