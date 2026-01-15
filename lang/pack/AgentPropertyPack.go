package pack

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/whatap/golib/io"
	"github.com/whatap/golib/lang/value"
)

type AgentPropertyPack struct {
	AbstractPack
	ctr   map[string]string
	table map[string]string
}

func NewAgentPropertyPack() *AgentPropertyPack {
	p := &AgentPropertyPack{
		ctr:   make(map[string]string),
		table: make(map[string]string),
	}
	return p
}

func (this *AgentPropertyPack) Size() int {
	return len(this.table)
}

func (this *AgentPropertyPack) IsEmpty() bool {
	return len(this.table) == 0
}

func (this *AgentPropertyPack) ContainsKey(key string) bool {
	_, exists := this.table[key]
	return exists
}

func (this *AgentPropertyPack) Keys() []string {
	keys := make([]string, 0, len(this.table))
	for k := range this.table {
		keys = append(keys, k)
	}
	return keys
}

func (this *AgentPropertyPack) Get(key string) string {
	return this.table[key]
}

func (this *AgentPropertyPack) GetPropertyAsMapValue() *value.MapValue {
	m := value.NewMapValue()
	for key, val := range this.table {
		m.PutString(key, val)
	}
	return m
}

func (this *AgentPropertyPack) GetInt(key string) int32 {
	v := this.Get(key)
	if i, err := strconv.ParseInt(v, 10, 32); err == nil {
		return int32(i)
	}
	return 0
}

func (this *AgentPropertyPack) GetLong(key string) int64 {
	v := this.Get(key)
	if i, err := strconv.ParseInt(v, 10, 64); err == nil {
		return i
	}
	return 0
}

func (this *AgentPropertyPack) GetFloat(key string) float32 {
	v := this.Get(key)
	if f, err := strconv.ParseFloat(v, 32); err == nil {
		return float32(f)
	}
	return 0.0
}

func (this *AgentPropertyPack) Put(key string, value string) {
	this.table[key] = value
}

func (this *AgentPropertyPack) PutAll(m *value.MapValue) {
	if m == nil {
		return
	}
	keys := m.Keys()
	for keys.HasMoreElements() {
		key := keys.NextString()
		val := m.Get(key)
		this.table[key] = castutil.cString(val)
	}
}

func (this *AgentPropertyPack) PutAllMap(m map[string]string) {
	if m != nil {
		for k, v := range m {
			this.table[k] = v
		}
	}
}

func (this *AgentPropertyPack) Remove(key string) {
	delete(this.table, key)
}

func (this *AgentPropertyPack) Clear() {
	this.table = make(map[string]string)
}

func (this *AgentPropertyPack) PutCtr(key string, value string) {
	this.ctr[key] = value
}

func (this *AgentPropertyPack) GetCtr(key string) string {
	return this.ctr[key]
}

func (this *AgentPropertyPack) String() string {
	var buf strings.Builder
	buf.WriteString("AgentProperty ")
	buf.WriteString(fmt.Sprintf("pcode=%d,oid=%d,okind=%d,onode=%d",
		this.Pcode, this.Oid, this.Okind, this.Onode))

	for key, val := range this.table {
		buf.WriteString(fmt.Sprintf(" %s=%s", key, val))
	}
	return buf.String()
}

func (this *AgentPropertyPack) ToPropertyString() string {
	var buf strings.Builder
	for key, val := range this.table {
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

func (this *AgentPropertyPack) writeTable(dout *io.DataOutputX, t map[string]string) {
	dout.WriteDecimal(int64(len(t)))
	for key, val := range t {
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

func (this *AgentPropertyPack) readTable(din *io.DataInputX) map[string]string {
	o := make(map[string]string)
	sz := int(din.ReadDecimal())
	for i := 0; i < sz; i++ {
		k := din.ReadText()
		v := din.ReadText()
		o[k] = v
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
	for key, val := range this.table {
		buf.WriteString(fmt.Sprintf("\t%s=%s\n", key, val))
	}
	return buf.String()
}
