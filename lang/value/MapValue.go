package value

import (
	"bytes"
	"fmt"

	"github.com/whatap/golib/io"
	"github.com/whatap/golib/util/hmap"
)

type MapValue struct {
	table *hmap.StringKeyLinkedMap
}

func NewMapValue() *MapValue {
	v := new(MapValue)
	v.table = hmap.NewStringKeyLinkedMap()
	return v
}

func (this *MapValue) ContainsKey(key string) bool {
	return this.table.ContainsKey(key)
}

func (this *MapValue) CompareTo(o Value) int {
	if o == nil {
		return 0
	}
	if o.GetValueType() != this.GetValueType() {
		return int(this.GetValueType() - o.GetValueType())
	}
	that := o.(*MapValue)
	if this.table.Size() != that.table.Size() {
		return this.table.Size() - that.table.Size()
	}
	keys := this.Keys()
	for keys.HasMoreElements() {
		key := keys.NextString()
		v1 := this.table.Get(key).(Value)
		v2 := that.table.Get(key).(Value)
		if v2 == nil {
			return 1
		}
		c := v1.CompareTo(v2)
		if c != 0 {
			return c
		}
	}
	return 0

}

func (this *MapValue) Equals(o Value) bool {
	if o == nil || o.GetValueType() != this.GetValueType() {
		return false
	}
	that := o.(*MapValue)
	if this.table.Size() != that.table.Size() {
		return false
	}
	keys := this.Keys()
	for keys.HasMoreElements() {
		key := keys.NextString()
		v1 := this.table.Get(key).(Value)
		v2 := that.table.Get(key).(Value)
		if v2 == nil {
			return false
		}
		if v1.Equals(v2) == false {
			return false
		}
	}
	return true
}

func (this *MapValue) IsEmpty() bool {
	return this.table.IsEmpty()
}

func (this *MapValue) Get(key string) Value {
	o := this.table.Get(key)
	if o == nil {
		return nil
	}
	return o.(Value)
}
func (this *MapValue) GetString(key string) string {
	o := this.table.Get(key)
	if o != nil && o.(Value).GetValueType() == VALUE_TEXT {
		t := o.(*TextValue)
		return t.Val
	}
	return ""
}
func (this *MapValue) GetBool(key string) bool {
	o := this.table.Get(key)
	if o != nil && o.(Value).GetValueType() == VALUE_BOOLEAN {
		t := o.(*BoolValue)
		return t.Val
	}
	return false
}

func (this *MapValue) GetLong(key string) int64 {
	o := this.table.Get(key)
	if o != nil && o.(Value).GetValueType() == VALUE_DECIMAL {
		t := o.(*DecimalValue)
		return t.Val
	}
	return 0
}

func (this *MapValue) GetFloat(key string) float32 {
	o := this.table.Get(key)
	if o != nil && o.(Value).GetValueType() == VALUE_FLOAT {
		t := o.(*FloatValue)
		return t.Val
	}
	return 0
}

func (this *MapValue) PutString(key string, v string) {
	this.Put(key, NewTextValue(v))
}
func (this *MapValue) PutLong(key string, v int64) {
	this.Put(key, NewDecimalValue(v))
}
func (this *MapValue) Put(key string, v Value) {
	this.table.Put(key, v)
}
func (this *MapValue) PutAll(m *MapValue) {
	strEnumer := m.Keys()
	for strEnumer.HasMoreElements() {
		key := strEnumer.NextString()
		this.table.Put(key, m.Get(key))
	}
}

func (this *MapValue) Clear() {
	this.table.Clear()
}
func (this *MapValue) Size() int {
	return this.table.Size()
}
func (this *MapValue) Keys() hmap.StringEnumer {
	return this.table.Keys()
}

func (this *MapValue) GetValueType() byte {
	return VALUE_MAP
}
func (this *MapValue) Write(dout *io.DataOutputX) {
	dout.WriteDecimal(int64(this.table.Size()))
	keys := this.Keys()
	for keys.HasMoreElements() {
		key := keys.NextString()
		value := this.table.Get(key).(Value)

		dout.WriteText(key)
		WriteValue(dout, value)
	}
}
func (this *MapValue) Read(din *io.DataInputX) {
	count := int(din.ReadDecimal())
	for t := 0; t < count; t++ {
		key := din.ReadText()
		value := ReadValue(din)
		this.table.Put(key, value)
	}
}

func (this *MapValue) NewList(name string) *ListValue {
	list := NewListValue(nil)
	this.Put(name, list)
	return list
}

func (this *MapValue) ToString() string {
	return this.String()
}
func (this *MapValue) String() string {
	var buf bytes.Buffer

	keys := this.Keys()

	buf.WriteString("{")
	for keys.HasMoreElements() {
		if buf.Len() > 1 {
			buf.WriteString(",")
		}
		key := keys.NextString()
		value := this.table.Get(key).(Value)

		buf.WriteString(key)
		buf.WriteString("=")
		buf.WriteString(fmt.Sprintf("%s", value))
	}
	buf.WriteString("}")
	return buf.String()
}
