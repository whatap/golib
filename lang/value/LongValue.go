package value

import (
	"strconv"

	"github.com/whatap/golib/io"
)

type LongValue struct {
	Val int64
}

func NewLongValue(v int64) *LongValue {
	m := new(LongValue)
	m.Val = v
	return m
}

func (this *LongValue) CompareTo(o Value) int {
	if o != nil && o.GetValueType() == this.GetValueType() {
		if this.Val == o.(*LongValue).Val {
			return 0
		}
		if this.Val < o.(*LongValue).Val {
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

func (this *LongValue) Equals(o Value) bool {
	if o != nil && o.GetValueType() == this.GetValueType() {
		return this.Val == o.(*LongValue).Val
	}
	return false
}

func (this *LongValue) GetValueType() byte {
	return VALUE_DECIMAL_LONG
}

func (this *LongValue) Write(out *io.DataOutputX) {
	out.WriteLong(this.Val)
}

func (this *LongValue) Read(in *io.DataInputX) {
	this.Val = in.ReadLong()
}

func (this *LongValue) String() string {
	return strconv.FormatInt(this.Val, 10)
}
