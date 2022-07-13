package value

import (
	"github.com/whatap/golib/io"
)

type ListValue struct {
	table []interface{}
}

func NewListValue(value []interface{}) *ListValue {
	v := new(ListValue)
	if value == nil {
		v.table = []interface{}{}
	} else {
		v.table = value
	}
	return v
}

func (this *ListValue) CompareTo(o Value) int {
	if o == nil {
		return 0
	}
	if o.GetValueType() != this.GetValueType() {
		return int(this.GetValueType() - o.GetValueType())
	}
	that := o.(*ListValue)
	if len(this.table) != len(that.table) {
		return len(this.table) - len(that.table)
	}
	for i := 0; i < len(this.table); i++ {
		v1 := this.table[i].(Value)
		v2 := that.table[i].(Value)
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

func (this *ListValue) Equals(o Value) bool {
	if o == nil || o.GetValueType() != this.GetValueType() {
		return false
	}
	that := o.(*ListValue)
	if len(this.table) != len(that.table) {
		return false
	}
	for i := 0; i < len(this.table); i++ {
		v1 := this.table[i].(Value)
		v2 := that.table[i].(Value)
		if v2 == nil {
			return false
		}
		if v1.Equals(v2) == false {
			return false
		}
	}
	return true
}

func (this *ListValue) Get(i int) Value {
	o := this.table[i]
	if o == nil {
		return nil
	}
	return o.(Value)
}
func (this *ListValue) GetString(i int) string {
	o := this.table[i]
	if o != nil && o.(Value).GetValueType() == VALUE_TEXT {
		t := o.(*TextValue)
		return t.Val
	}
	return ""
}
func (this *ListValue) GetBool(i int) bool {
	o := this.table[i]
	if o != nil && o.(Value).GetValueType() == VALUE_BOOLEAN {
		t := o.(*BoolValue)
		return t.Val
	}
	return false
}

func (this *ListValue) AddString(value string) {
	this.table = append(this.table, NewTextValue(value))
}
func (this *ListValue) AddLong(value int64) {
	this.table = append(this.table, NewDecimalValue(value))
}
func (this *ListValue) Add(value Value) {
	this.table = append(this.table, value)
}
func (this *ListValue) Set(idx int, value Value) {
	this.table[idx] = value
}
func (this *ListValue) Clear() {
	this.table = []interface{}{}
}
func (this *ListValue) Size() int {
	return len(this.table)
}

func (this *ListValue) GetValueType() byte {
	return VALUE_LIST
}
func (this *ListValue) Write(dout *io.DataOutputX) {
	if this.table == nil || len(this.table) == 0 {
		dout.WriteDecimal(0)
		return
	}
	sz := len(this.table)
	dout.WriteDecimal(int64(sz))
	for i := 0; i < sz; i++ {
		value := this.table[i].(Value)
		WriteValue(dout, value)
	}
}
func (this *ListValue) Read(din *io.DataInputX) {
	count := int(din.ReadDecimal())
	if count == 0 {
		return
	}
	this.table = make([]interface{}, count)
	for t := 0; t < count; t++ {
		this.table[t] = ReadValue(din)
	}
}
