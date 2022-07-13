package value

import (
	"github.com/whatap/golib/io"
)

type TextHashValue struct {
	Val int32
}

func NewTextHashValue(v int32) *TextHashValue {
	m := new(TextHashValue)
	m.Val = v
	return m
}

func (this *TextHashValue) CompareTo(o Value) int {
	if o != nil && o.GetValueType() == this.GetValueType() {
		if this.Val == o.(*TextHashValue).Val {
			return 0
		}
		if this.Val < o.(*TextHashValue).Val {
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

func (this *TextHashValue) Equals(o Value) bool {
	if o != nil && o.GetValueType() == this.GetValueType() {
		return this.Val == o.(*TextHashValue).Val
	}
	return false
}

func (this *TextHashValue) GetValueType() byte {
	return VALUE_TEXT_HASH
}

func (this *TextHashValue) Write(out *io.DataOutputX) {
	out.WriteInt(this.Val)
}

func (this *TextHashValue) Read(in *io.DataInputX) {
	this.Val = in.ReadInt()
}
