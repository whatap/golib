package pack

import (
	"fmt"
	"strings"

	"github.com/whatap/golib/io"
	"github.com/whatap/golib/lang/value"
	"github.com/whatap/golib/util/castutil"
	"github.com/whatap/golib/util/hmap"
)

type AgentPropertyPack struct {
	AbstractPack
	ctr   *hmap.StringKeyLinkedMap
	table *hmap.StringKeyLinkedMap
}

func NewAgentPropertyPack() *AgentPropertyPack {
	p := new(AgentPropertyPack)
	p.ctr = hmap.NewStringKeyLinkedMap()
	p.table = hmap.NewStringKeyLinkedMap()
	return p
}

func (this *AgentPropertyPack) Size() int {
	return this.table.Size()
}

func (this *AgentPropertyPack) IsEmpty() bool {
	return this.table.Size() == 0
}

func (this *AgentPropertyPack) ContainsKey(key string) bool {
	return this.table.ContainsKey(key)
}

func (this *AgentPropertyPack) Keys() hmap.StringEnumer {
	return this.table.Keys()
}

func (this *AgentPropertyPack) GetPropertyAsMapValue() *value.MapValue {
	m := value.NewMapValue()
	en := this.table.Keys()
	for en.HasMoreElements() {
		key := en.NextString()
		val := castutil.CString(m.Get(key))
		m.PutString(key, val)
	}
	return m
}

func (this *AgentPropertyPack) Get(key string) string {
	return castutil.CString(this.table.Get(key))
}

func (this *AgentPropertyPack) GetInt(key string) int32 {
	v := this.Get(key)
	return castutil.CInt(v)
}

func (this *AgentPropertyPack) GetLong(key string) int64 {
	v := this.Get(key)
	return castutil.CLong(v)
}

func (this *AgentPropertyPack) GetFloat(key string) float32 {
	v := this.Get(key)
	return castutil.CFloat(v)
}

func (this *AgentPropertyPack) Put(key string, value string) {
	this.table.Put(key, value)
}

func (this *AgentPropertyPack) PutAll(m *value.MapValue) {
	if m == nil {
		return
	}
	keys := m.Keys()
	for keys.HasMoreElements() {
		key := keys.NextString()
		val := m.Get(key)
		this.table.Put(key, castutil.CString(val))
	}
}

func (this *AgentPropertyPack) PutAllMap(m *hmap.StringKeyLinkedMap) {
	if m != nil {
		en := m.Keys()
		for en.HasMoreElements() {
			key := en.NextString()
			val := castutil.CString(m.Get(key))
			this.table.Put(key, val)
		}
	}
}
func (this *AgentPropertyPack) PutMapValue(m *value.MapValue) {
	if m != nil {
		en := m.Keys()
		for en.HasMoreElements() {
			key := en.NextString()
			val := castutil.CString(m.Get(key))
			this.table.Put(key, val)
		}
	}
}

func (this *AgentPropertyPack) Remove(key string) {
	this.table.Remove(key)
}

func (this *AgentPropertyPack) Clear() {
	this.table.Clear()
}

func (this *AgentPropertyPack) PutCtr(key string, value string) {
	this.ctr.Put(key, value)
}

func (this *AgentPropertyPack) GetCtr(key string) string {
	return castutil.CString(this.ctr.Get(key))
}

func (this *AgentPropertyPack) String() string {
	var buf strings.Builder
	buf.WriteString("AgentProperty ")
	buf.WriteString(fmt.Sprintf("pcode=%d,oid=%d,okind=%d,onode=%d",
		this.Pcode, this.Oid, this.Okind, this.Onode))

	en := this.table.Keys()
	for en.HasMoreElements() {
		key := en.NextString()
		val := castutil.CString(this.table.Get(key))
		buf.WriteString(fmt.Sprintf(" %s=%s", key, val))
	}
	return buf.String()
}

func (this *AgentPropertyPack) ToPropertyString() string {
	var buf strings.Builder
	en := this.table.Keys()
	for en.HasMoreElements() {
		key := en.NextString()
		val := castutil.CString(this.table.Get(key))
		buf.WriteString(fmt.Sprintf(" %s=%s", key, val))
	}
	return buf.String()
}

func (this *AgentPropertyPack) GetPackType() int16 {
	return PACK_AGENT_PROPERTY
}

func (this *AgentPropertyPack) Write(dout *io.DataOutputX) {
	this.AbstractPack.Write(dout)
	dout.WriteByte(0) // version
	this.writeTable(dout, this.ctr)
	this.writeTable(dout, this.table)
}

func (this *AgentPropertyPack) writeTable(dout *io.DataOutputX, t *hmap.StringKeyLinkedMap) {
	dout.WriteDecimal(int64(t.Size()))
	en := t.Keys()
	for en.HasMoreElements() {
		key := en.NextString()
		val := castutil.CString(t.Get(key))
		dout.WriteText(key)
		dout.WriteText(val)
	}
}

func (this *AgentPropertyPack) Read(din *io.DataInputX) {
	this.AbstractPack.Read(din)
	_ = din.ReadByte() // version
	this.ctr = this.readTable(din)
	this.table = this.readTable(din)
}

func (this *AgentPropertyPack) readTable(din *io.DataInputX) *hmap.StringKeyLinkedMap {
	o := hmap.NewStringKeyLinkedMap()
	sz := int(din.ReadDecimal())
	for i := 0; i < sz; i++ {
		k := din.ReadText()
		v := din.ReadText()
		o.Put(k, v)
	}
	return o
}

func (this *AgentPropertyPack) SetMapValue(mapValue *value.MapValue) *AgentPropertyPack {
	if mapValue == nil {
		return this
	}
	this.PutAll(mapValue)
	return this
}

func (this *AgentPropertyPack) ToFormatString() string {
	var buf strings.Builder
	buf.WriteString(this.AbstractPack.ToString())
	en := this.table.Keys()
	for en.HasMoreElements() {
		key := en.NextString()
		val := castutil.CString(this.table.Get(key))
		buf.WriteString(fmt.Sprintf("\t%s=%s\n", key, val))
	}
	return buf.String()
}
