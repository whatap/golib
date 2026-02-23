package value

import (
	"strconv"

	"github.com/whatap/golib/io"
)

type IntValue struct {
	Val int32
}

func NewIntValue(v int32) *IntValue {
	m := new(IntValue)
	m.Val = v
	return m
}

func (this *IntValue) CompareTo(o Value) int {
	if o != nil && o.GetValueType() == this.GetValueType() {
		if this.Val == o.(*IntValue).Val {
			return 0
		}
		if this.Val < o.(*IntValue).Val {
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

func (this *IntValue) Equals(o Value) bool {
	if o != nil && o.GetValueType() == this.GetValueType() {
		return this.Val == o.(*IntValue).Val
	}
	return false
}

func (this *IntValue) GetValueType() byte {
	return VALUE_DECIMAL_INT
}

func (this *IntValue) Write(out *io.DataOutputX) {
	out.WriteInt(this.Val)
}

func (this *IntValue) Read(in *io.DataInputX) {
	this.Val = in.ReadInt()
}

func (this *IntValue) String() string {
	return strconv.FormatInt(int64(this.Val), 10)
}
