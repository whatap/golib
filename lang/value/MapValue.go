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

func (this *MapValue) GetInt(key string) int32 {
	o := this.table.Get(key)
	if o == nil {
		return 0
	}
	switch o.(Value).GetValueType() {
	case VALUE_DECIMAL:
		return int32(o.(*DecimalValue).Val)
	case VALUE_DECIMAL_INT:
		return o.(*IntValue).Val
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

func (this *MapValue) GetDouble(key string) float64 {
	o := this.table.Get(key)
	if o != nil && o.(Value).GetValueType() == VALUE_DOUBLE {
		t := o.(*DoubleValue)
		return t.Val
	}
	return 0
}

func (this *MapValue) GetStringDefault(key string, def string) string {
	o := this.table.Get(key)
	if o != nil && o.(Value).GetValueType() == VALUE_TEXT {
		return o.(*TextValue).Val
	}
	return def
}

func (this *MapValue) GetLongDefault(key string, def int64) int64 {
	o := this.table.Get(key)
	if o != nil && o.(Value).GetValueType() == VALUE_DECIMAL {
		return o.(*DecimalValue).Val
	}
	return def
}

func (this *MapValue) GetIntDefault(key string, def int32) int32 {
	o := this.table.Get(key)
	if o == nil {
		return def
	}
	switch o.(Value).GetValueType() {
	case VALUE_DECIMAL:
		return int32(o.(*DecimalValue).Val)
	case VALUE_DECIMAL_INT:
		return o.(*IntValue).Val
	}
	return def
}

func (this *MapValue) GetFloatDefault(key string, def float32) float32 {
	o := this.table.Get(key)
	if o != nil && o.(Value).GetValueType() == VALUE_FLOAT {
		return o.(*FloatValue).Val
	}
	return def
}

func (this *MapValue) GetBoolDefault(key string, def bool) bool {
	o := this.table.Get(key)
	if o != nil && o.(Value).GetValueType() == VALUE_BOOLEAN {
		return o.(*BoolValue).Val
	}
	return def
}

func (this *MapValue) GetList(key string) *ListValue {
	o := this.table.Get(key)
	if o != nil && o.(Value).GetValueType() == VALUE_LIST {
		return o.(*ListValue)
	}
	return nil
}

func (this *MapValue) GetListNotNull(key string) *ListValue {
	list := this.GetList(key)
	if list == nil {
		return NewListValue(nil)
	}
	return list
}

func (this *MapValue) KeyArray() []string {
	keys := this.Keys()
	result := make([]string, 0, this.Size())
	for keys.HasMoreElements() {
		result = append(result, keys.NextString())
	}
	return result
}

func (this *MapValue) PutString(key string, v string) {
	this.Put(key, NewTextValue(v))
}
func (this *MapValue) PutLong(key string, v int64) {
	this.Put(key, NewDecimalValue(v))
}
func (this *MapValue) PutInt(key string, v int32) {
	this.Put(key, NewIntValue(v))
}
func (this *MapValue) PutFloat(key string, v float32) {
	this.Put(key, NewFloatValue(v))
}
func (this *MapValue) PutDouble(key string, v float64) {
	this.Put(key, NewDoubleValue(v))
}
func (this *MapValue) PutBool(key string, v bool) {
	this.Put(key, NewBoolValue(v))
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

func (this *MapValue) Remove(key string) Value {
	o := this.table.Remove(key)
	if o == nil {
		return nil
	}
	return o.(Value)
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

func (this *MapValue) ToJsonString() string {
	var buf bytes.Buffer

	keys := this.Keys()

	buf.WriteString("{")
	for keys.HasMoreElements() {
		if buf.Len() > 1 {
			buf.WriteString(",")
		}
		key := keys.NextString()
		value := this.table.Get(key).(Value)

		buf.WriteString(quoteString(key))
		buf.WriteString(":")
		buf.WriteString(valueToJsonString(value))
	}
	buf.WriteString("}")
	return buf.String()
}
