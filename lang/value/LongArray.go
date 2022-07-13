package value

import (
	"github.com/whatap/golib/io"
	"github.com/whatap/golib/util/compare"
)

type LongArray struct {
	Val []int64
}

func NewLongArray(v []int64) *LongArray {
	m := new(LongArray)
	m.Val = v
	return m
}

func (this *LongArray) CompareTo(o Value) int {
	if o != nil && o.GetValueType() == this.GetValueType() {
		return compare.CompareToLongs(this.Val, o.(*LongArray).Val)
	}
	if o == nil {
		return 1
	} else {
		return int(this.GetValueType() - o.GetValueType())
	}
}

func (this *LongArray) Equals(o Value) bool {
	if o != nil && o.GetValueType() == this.GetValueType() {
		return compare.EqualLongs(this.Val, o.(*LongArray).Val)
	}
	return false
}

func (this *LongArray) GetValueType() byte {
	return ARRAY_LONG
}
func (this *LongArray) Write(out *io.DataOutputX) {
	out.WriteLongArray(this.Val)
}
func (this *LongArray) Read(in *io.DataInputX) {
	this.Val = in.ReadLongArray()
}
