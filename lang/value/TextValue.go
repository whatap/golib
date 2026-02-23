package value

import (
	"github.com/whatap/golib/io"
)

type TextValue struct {
	Val string
}

func NewTextValue(v string) *TextValue {
	return &TextValue{Val: v}
}

func (this *TextValue) CompareTo(o Value) int {
	if o != nil && o.GetValueType() == this.GetValueType() {
		if this.Val == o.(*TextValue).Val {
			return 0
		}
		if this.Val < o.(*TextValue).Val {
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

func (this *TextValue) Equals(o Value) bool {
	if o != nil && o.GetValueType() == this.GetValueType() {
		return this.Val == o.(*TextValue).Val
	}
	return false
}

func (this *TextValue) GetValueType() byte {
	return VALUE_TEXT
}
func (this *TextValue) Write(out *io.DataOutputX) {
	out.WriteText(this.Val)
}
func (this *TextValue) Read(in *io.DataInputX) {
	this.Val = in.ReadText()
}
func (this *TextValue) String() string {
	return this.Val
}
