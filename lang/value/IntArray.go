package value

import (
	"github.com/whatap/golib/io"
	"github.com/whatap/golib/util/compare"
)

type IntArray struct {
	Val []int32
}

func NewIntArray(v []int32) *IntArray {
	m := new(IntArray)
	m.Val = v
	return m
}

func (this *IntArray) CompareTo(o Value) int {
	if o != nil && o.GetValueType() == this.GetValueType() {
		return compare.CompareToInts(this.Val, o.(*IntArray).Val)
	}
	if o == nil {
		return 1
	} else {
		return int(this.GetValueType() - o.GetValueType())
	}
}

func (this *IntArray) Equals(o Value) bool {
	if o != nil && o.GetValueType() == this.GetValueType() {
		return compare.EqualInts(this.Val, o.(*IntArray).Val)
	}
	return false
}

func (this *IntArray) GetValueType() byte {
	return ARRAY_INT
}
func (this *IntArray) Write(out *io.DataOutputX) {
	out.WriteIntArray(this.Val)
}
func (this *IntArray) Read(in *io.DataInputX) {
	this.Val = in.ReadIntArray()
}
