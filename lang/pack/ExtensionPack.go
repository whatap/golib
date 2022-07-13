package pack

import (
	"github.com/whatap/golib/io"
	"github.com/whatap/golib/lang/value"
	"github.com/whatap/golib/util/hmap"
)

type ExtensionPack struct {
	AbstractPack

	IsProjectWide bool
	Header        *hmap.StringIntLinkedMap
	Value         *value.IntMapValue
}

func NewExtensionPack() *ExtensionPack {
	p := new(ExtensionPack)
	p.Header = hmap.NewStringIntLinkedMap()
	p.Value = value.NewIntMapValue()

	return p
}

func (this *ExtensionPack) GetPackType() int16 {
	return PACK_EXTENSION
}

func (this *ExtensionPack) Write(dout *io.DataOutputX) {
	this.AbstractPack.Write(dout)
	dout.WriteByte(0) //version
	dout.WriteBool(this.IsProjectWide)
	toHeaderBytes(dout, this.Header)
	value.WriteValue(dout, this.Value)
}
func (this *ExtensionPack) Read(din *io.DataInputX) {
	this.AbstractPack.Read(din)
	din.ReadByte() //version
	this.IsProjectWide = din.ReadBool()
	this.Header = toHeaderObject(din)
	this.Value = value.ReadValue(din).(*value.IntMapValue)
}

func toHeaderBytes(dout *io.DataOutputX, m *hmap.StringIntLinkedMap) {
	dout.WriteDecimal(int64(m.Size()))
	en := m.Entries()
	for en.HasMoreElements() {
		e := en.NextElement().(*hmap.StringIntLinkedEntry)
		dout.WriteText(e.GetKey())
		dout.WriteInt(e.GetValue())
	}
}

func toHeaderObject(in *io.DataInputX) *hmap.StringIntLinkedMap {
	m := hmap.NewStringIntLinkedMap()
	cnt := int(in.ReadDecimal())
	for i := 0; i < cnt; i++ {
		key := in.ReadText()
		value := in.ReadInt()
		m.Put(key, value)
	}
	return m
}
