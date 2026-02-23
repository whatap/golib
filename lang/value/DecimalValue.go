package value

import (
	"strconv"

	"github.com/whatap/golib/io"
)

type DecimalValue struct {
	Val int64
}

func NewDecimalValue(v int64) *DecimalValue {
	m := new(DecimalValue)
	m.Val = v
	return m
}

func (this *DecimalValue) CompareTo(o Value) int {
	if o != nil && o.GetValueType() == this.GetValueType() {
		if this.Val == o.(*DecimalValue).Val {
			return 0
		}
		if this.Val < o.(*DecimalValue).Val {
			return 1
		} else {
			return -1
		}
	}
	if o == nil {
		return 1
	} else {
		return int(this.GetValueType() - o.GetValueType())
	}
}

func (this *DecimalValue) Equals(o Value) bool {
	if o != nil && o.GetValueType() == this.GetValueType() {
		return this.Val == o.(*DecimalValue).Val
	}
	return false
}

func (this *DecimalValue) GetValueType() byte {
	return VALUE_DECIMAL
}

func (this *DecimalValue) Write(out *io.DataOutputX) {
	out.WriteDecimal(this.Val)
}

func (this *DecimalValue) Read(in *io.DataInputX) {
	this.Val = in.ReadDecimal()
}
func (this *DecimalValue) String() string {
	return strconv.FormatInt(this.Val, 10)
}
