package value

import (
	"strconv"

	"github.com/whatap/golib/io"
)

type DoubleValue struct {
	Val float64
}

func NewDoubleValue(v float64) *DoubleValue {
	m := new(DoubleValue)
	m.Val = v
	return m
}

func (this *DoubleValue) CompareTo(o Value) int {
	if o != nil && o.GetValueType() == this.GetValueType() {
		if this.Val == o.(*DoubleValue).Val {
			return 0
		}
		if this.Val < o.(*DoubleValue).Val {
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

func (this *DoubleValue) Equals(o Value) bool {
	if o != nil && o.GetValueType() == this.GetValueType() {
		return this.Val == o.(*DoubleValue).Val
	}
	return false
}

func (this *DoubleValue) GetValueType() byte {
	return VALUE_DOUBLE
}
func (this *DoubleValue) Write(out *io.DataOutputX) {
	out.WriteDouble(this.Val)
}
func (this *DoubleValue) Read(in *io.DataInputX) {
	this.Val = in.ReadDouble()
}

func (this *DoubleValue) String() string {
	return strconv.FormatFloat(this.Val, 'f', 7, 64)
}
