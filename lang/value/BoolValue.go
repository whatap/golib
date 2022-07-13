package value

import (
	"github.com/whatap/golib/io"
)

type BoolValue struct {
	Val bool
}

func NewBoolValue(v bool) *BoolValue {
	m := new(BoolValue)
	m.Val = v
	return m
}

func (this *BoolValue) CompareTo(o Value) int {
	if o != nil && o.GetValueType() == this.GetValueType() {
		if this.Val == o.(*BoolValue).Val {
			return 0
		}
		if this.Val {
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

func (this *BoolValue) Equals(o Value) bool {
	if o != nil && o.GetValueType() == this.GetValueType() {
		return this.Val == o.(*BoolValue).Val
	}
	return false
}

func (this *BoolValue) GetValueType() byte {
	return VALUE_BOOLEAN
}
func (this *BoolValue) Write(out *io.DataOutputX) {
	out.WriteBool(this.Val)
}
func (this *BoolValue) Read(in *io.DataInputX) {
	this.Val = in.ReadBool()
}
