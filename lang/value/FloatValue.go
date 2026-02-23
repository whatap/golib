package value

import (
	"strconv"

	"github.com/whatap/golib/io"
)

type FloatValue struct {
	Val float32
}

func NewFloatValue(v float32) *FloatValue {
	m := new(FloatValue)
	m.Val = v
	return m
}

func (this *FloatValue) CompareTo(o Value) int {
	if o != nil && o.GetValueType() == this.GetValueType() {
		if this.Val == o.(*FloatValue).Val {
			return 0
		}
		if this.Val < o.(*FloatValue).Val {
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

func (this *FloatValue) Equals(o Value) bool {
	if o != nil && o.GetValueType() == this.GetValueType() {
		return this.Val == o.(*FloatValue).Val
	}
	return false
}

func (this *FloatValue) GetValueType() byte {
	return VALUE_FLOAT
}
func (this *FloatValue) Write(out *io.DataOutputX) {
	out.WriteFloat(this.Val)
}
func (this *FloatValue) Read(in *io.DataInputX) {
	this.Val = in.ReadFloat()
}

func (this *FloatValue) String() string {
	return strconv.FormatFloat(float64(this.Val), 'f', 7, 32)
}
