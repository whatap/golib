package value

import (
	"github.com/whatap/golib/io"
	"github.com/whatap/golib/util/compare"
)

type TextArray struct {
	Val []string
}

func NewTextArray(v []string) *TextArray {
	m := new(TextArray)
	m.Val = v
	return m
}

func (this *TextArray) CompareTo(o Value) int {
	if o != nil && o.GetValueType() == this.GetValueType() {
		return compare.CompareToStrings(this.Val, o.(*TextArray).Val)
	}
	if o == nil {
		return 1
	} else {
		return int(this.GetValueType() - o.GetValueType())
	}
}

func (this *TextArray) Equals(o Value) bool {
	if o != nil && o.GetValueType() == this.GetValueType() {
		return compare.EqualStrings(this.Val, o.(*TextArray).Val)
	}
	return false
}

func (this *TextArray) GetValueType() byte {
	return ARRAY_TEXT
}
func (this *TextArray) Write(out *io.DataOutputX) {
	out.WriteTextArray(this.Val)
}
func (this *TextArray) Read(in *io.DataInputX) {
	this.Val = in.ReadTextArray()
}
