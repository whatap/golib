package pack

import (
	"github.com/whatap/golib/io"
)

type DownCheckRec struct {
	Name string
	Host string
	Port int32
	Ok   bool
}

type SMDownCheckPack struct {
	AbstractPack
	ver         byte
	Records     []byte
	RecordCount int32
}

func NewSMDownCheckPack() *SMDownCheckPack {
	p := new(SMDownCheckPack)
	p.ver = 1
	return p
}

func (this *SMDownCheckPack) GetPackType() int16 {
	return PACK_SM_DOWN_CHECK
}

func (this *SMDownCheckPack) Write(dout *io.DataOutputX) {
	this.AbstractPack.Write(dout)
	dout.WriteByte(this.ver)
	dout.WriteBlob(this.Records)
	dout.WriteDecimal(int64(this.RecordCount))

}
func (this *SMDownCheckPack) Read(din *io.DataInputX) {
	this.AbstractPack.Read(din)
	this.ver = din.ReadByte()
	this.Records = din.ReadBlob()
	this.RecordCount = int32(din.ReadDecimal())
}
func (this *SMDownCheckPack) WriteRec(o *io.DataOutputX, m *DownCheckRec) {
	o.WriteText(m.Name)
	o.WriteText(m.Host)
	o.WriteInt(m.Port)
	o.WriteBool(m.Ok)
}

func (this *SMDownCheckPack) ReadRec(in *io.DataInputX) *DownCheckRec {
	m := new(DownCheckRec)
	m.Name = in.ReadText()
	m.Host = in.ReadText()
	m.Port = in.ReadInt()
	m.Ok = in.ReadBool()
	return m
}

func (this *SMDownCheckPack) GetRecords() []*DownCheckRec {
	in := io.NewDataInputX(this.Records)
	sz := int(in.ReadShort())
	items := make([]*DownCheckRec, sz)
	for i := 0; i < sz; i++ {
		items[i] = this.ReadRec(in)
	}
	return items
}

func (this *SMDownCheckPack) SetRecords(items []*DownCheckRec) {
	out := io.NewDataOutputX()
	sz := len(items)
	out.WriteShort(int16(sz))
	for i := 0; i < sz; i++ {
		this.WriteRec(out, items[i])
	}
	this.Records = out.ToByteArray()
	this.RecordCount = int32(sz)
}
