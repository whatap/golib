package value

import (
	"github.com/whatap/golib/io"
	"github.com/whatap/golib/util/hmap"
)

type IntMapValue struct {
	table *hmap.IntKeyLinkedMap
}

func NewIntMapValue() *IntMapValue {
	v := new(IntMapValue)
	v.table = hmap.NewIntKeyLinkedMapDefault()
	return v
}

func (this *IntMapValue) CompareTo(o Value) int {
	if o == nil {
		return 0
	}
	if o.GetValueType() != this.GetValueType() {
		return int(this.GetValueType() - o.GetValueType())
	}
	that := o.(*IntMapValue)
	if this.table.Size() != that.table.Size() {
		return this.table.Size() - that.table.Size()
	}
	keys := this.Keys()
	for keys.HasMoreElements() {
		key := keys.NextInt()
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

func (this *IntMapValue) Equals(o Value) bool {
	if o == nil || o.GetValueType() != this.GetValueType() {
		return false
	}
	that := o.(*IntMapValue)
	if this.table.Size() != that.table.Size() {
		return false
	}
	keys := this.Keys()
	for keys.HasMoreElements() {
		key := keys.NextInt()
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

func (this *IntMapValue) Get(key int32) Value {
	o := this.table.Get(key)
	if o == nil {
		return nil
	}
	return o.(Value)
}
func (this *IntMapValue) GetString(key int32) string {
	o := this.table.Get(key)
	if o != nil && o.(Value).GetValueType() == VALUE_TEXT {
		t := o.(*TextValue)
		return t.Val
	}
	return ""
}
func (this *IntMapValue) GetBool(key int32) bool {
	o := this.table.Get(key)
	if o != nil && o.(Value).GetValueType() == VALUE_BOOLEAN {
		t := o.(*BoolValue)
		return t.Val
	}
	return false
}

func (this *IntMapValue) PutString(key int32, value string) {
	this.Put(key, NewTextValue(value))
}
func (this *IntMapValue) PutLong(key int32, value int64) {
	this.Put(key, NewDecimalValue(value))
}
func (this *IntMapValue) Put(key int32, value Value) {
	this.table.Put(key, value)
}
func (this *IntMapValue) Clear() {
	this.table.Clear()
}
func (this *IntMapValue) Size() int {
	return this.table.Size()
}
func (this *IntMapValue) Keys() hmap.IntEnumer {
	return this.table.Keys()
}

func (this *IntMapValue) GetValueType() byte {
	return INT_VALUE_MAP
}
func (this *IntMapValue) Write(dout *io.DataOutputX) {
	dout.WriteDecimal(int64(this.table.Size()))
	keys := this.Keys()
	for keys.HasMoreElements() {
		key := keys.NextInt()
		value := this.table.Get(key).(Value)

		dout.WriteInt(key)
		WriteValue(dout, value)
	}
}
func (this *IntMapValue) Read(din *io.DataInputX) {
	count := int(din.ReadDecimal())
	for t := 0; t < count; t++ {
		key := din.ReadInt()
		value := ReadValue(din)
		this.table.Put(key, value)
	}
}
func (this *IntMapValue) WriteValue(out *io.DataOutputX) *io.DataOutputX {
	out.WriteByte(this.GetValueType())
	this.Write(out)
	return out
}

func (this *IntMapValue) NewList(key int32) *ListValue {
	list := NewListValue(nil)
	this.Put(key, list)
	return list
}
