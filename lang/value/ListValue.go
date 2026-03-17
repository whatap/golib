package value

import (
	"bytes"
	"fmt"
	"strconv"

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

func (this *ListValue) ToJsonString() string {
	var buf bytes.Buffer
	buf.WriteString("[")
	for i := 0; i < len(this.table); i++ {
		if i > 0 {
			buf.WriteString(",")
		}
		value := this.table[i].(Value)
		buf.WriteString(valueToJsonString(value))
	}
	buf.WriteString("]")
	return buf.String()
}

func (this *ListValue) String() string {
	var buf bytes.Buffer
	buf.WriteString("[")
	for i := 0; i < len(this.table); i++ {
		if i > 0 {
			buf.WriteString(",")
		}
		value := this.table[i].(Value)
		buf.WriteString(fmt.Sprintf("%s", value))
	}
	buf.WriteString("]")
	return buf.String()
}

func valueToJsonString(value Value) string {
	if value == nil {
		return "null"
	}
	switch value.GetValueType() {
	case VALUE_TEXT:
		return quoteString(value.(*TextValue).Val)
	case VALUE_MAP:
		return value.(*MapValue).ToJsonString()
	case VALUE_LIST:
		return value.(*ListValue).ToJsonString()
	case VALUE_BOOLEAN:
		if value.(*BoolValue).Val {
			return "true"
		}
		return "false"
	case VALUE_NULL:
		return "null"
	case VALUE_DECIMAL:
		return strconv.FormatInt(value.(*DecimalValue).Val, 10)
	case VALUE_FLOAT:
		return strconv.FormatFloat(float64(value.(*FloatValue).Val), 'f', -1, 32)
	case VALUE_DOUBLE:
		return strconv.FormatFloat(value.(*DoubleValue).Val, 'f', -1, 64)
	default:
		return fmt.Sprintf("%v", value)
	}
}

func quoteString(s string) string {
	var buf bytes.Buffer
	buf.WriteString("\"")
	for _, r := range s {
		switch r {
		case '"':
			buf.WriteString("\\\"")
		case '\\':
			buf.WriteString("\\\\")
		case '\n':
			buf.WriteString("\\n")
		case '\r':
			buf.WriteString("\\r")
		case '\t':
			buf.WriteString("\\t")
		default:
			if r < 32 {
				buf.WriteString(fmt.Sprintf("\\u%04x", r))
			} else {
				buf.WriteRune(r)
			}
		}
	}
	buf.WriteString("\"")
	return buf.String()
}
