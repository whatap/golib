package value

import (
	"fmt"
	"strings"

	"github.com/whatap/golib/io"
)

const ARRAY_BOOLEAN = 76

type BooleanArray struct {
	Val []bool
}

func NewBooleanArray(v []bool) *BooleanArray {
	m := new(BooleanArray)
	m.Val = v
	return m
}

func (this *BooleanArray) Length() int {
	if this.Val == nil {
		return 0
	}
	return len(this.Val)
}

func (this *BooleanArray) Get(i int) Value {
	return NewBoolValue(this.Val[i])
}

func (this *BooleanArray) CompareTo(o Value) int {
	if o != nil && o.GetValueType() == this.GetValueType() {
		other := o.(*BooleanArray)
		if len(this.Val) != len(other.Val) {
			return len(this.Val) - len(other.Val)
		}
		return 0
	}
	if o == nil {
		return 1
	}
	return int(this.GetValueType() - o.GetValueType())
}

func (this *BooleanArray) Equals(o Value) bool {
	if o != nil && o.GetValueType() == this.GetValueType() {
		other := o.(*BooleanArray)
		if len(this.Val) != len(other.Val) {
			return false
		}
		for i := 0; i < len(this.Val); i++ {
			if this.Val[i] != other.Val[i] {
				return false
			}
		}
		return true
	}
	return false
}

func (this *BooleanArray) GetValueType() byte {
	return ARRAY_BOOLEAN
}

func (this *BooleanArray) Write(out *io.DataOutputX) {
	out.WriteDecimal(int64(len(this.Val)))
	for _, v := range this.Val {
		out.WriteBool(v)
	}
}

func (this *BooleanArray) Read(in *io.DataInputX) {
	count := int(in.ReadDecimal())
	this.Val = make([]bool, count)
	for i := 0; i < count; i++ {
		this.Val[i] = in.ReadBool()
	}
}

func (this *BooleanArray) String() string {
	parts := make([]string, len(this.Val))
	for i, v := range this.Val {
		parts[i] = fmt.Sprintf("%v", v)
	}
	return "[" + strings.Join(parts, ", ") + "]"
}
