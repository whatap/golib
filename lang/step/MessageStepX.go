package step

import (
	// "encoding/json"

	"github.com/whatap/golib/io"
	"github.com/whatap/golib/lang/value"
)

const (
	SINGLE_LINE_DISPLAY = 0x00000001
)

type MessageStepX struct {
	AbstractStep

	Title string
	Desc  string
	Ctr   int32
	Attr  *value.MapValue
}

func NewMessageStepX() *MessageStepX {
	p := new(MessageStepX)
	return p
}

func NewMessageStepXWithStartTime(startTime int32) *MessageStepX {
	p := NewMessageStepX()
	p.StartTime = startTime
	return p
}

func (this *MessageStepX) GetStepType() byte {
	return STEP_MESSAGE_X
}

func (this *MessageStepX) Write(out *io.DataOutputX) {
	this.AbstractStep.Write(out)
	// version
	out.WriteByte(0)
	out.WriteBlob(this.WriteVer0())
}

func (this *MessageStepX) WriteVer0() []byte {
	o := io.NewDataOutputX()
	o.WriteText(this.Title)
	o.WriteText(this.Desc)
	o.WriteInt(this.Ctr)
	if this.Attr != nil {
		value.WriteMapValue(o, this.Attr)
	}
	return o.ToByteArray()
}

func (this *MessageStepX) Read(in *io.DataInputX) {
	this.AbstractStep.Read(in)
	ver := in.ReadByte()
	bytes := in.ReadBlob()
	if ver == 0 {
		this.ReadVer0(bytes)
	}
}

func (this *MessageStepX) ReadVer0(bytes []byte) {
	i := io.NewDataInputX(bytes)
	this.Title = i.ReadText()
	this.Desc = i.ReadText()
	this.Ctr = i.ReadInt()
	val := value.ReadValue(i)
	if mv, ok := val.(*value.MapValue); ok {
		this.Attr = mv
	}

}

func (this *MessageStepX) GetElapsed() int32 {
	return 0 //this.Elapsed
}

func (this *MessageStepX) SetCtr(key int) {
	this.Ctr = this.Ctr | int32(key)
}

func (this *MessageStepX) CtrToJson() map[string]interface{} {
	obj := make(map[string]interface{})
	if this.isA(SINGLE_LINE_DISPLAY) {
		obj["SINGLE_LINE_DISPLAY"] = true
	}
	return obj
}

func (this *MessageStepX) isA(k int) bool {
	return (int(this.Ctr) & k) != 0
}
