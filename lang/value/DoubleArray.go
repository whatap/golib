package value

import (
	"fmt"
	"strings"

	"github.com/whatap/golib/io"
)

const ARRAY_DOUBLE = 75

type DoubleArray struct {
	Val []float64
}

func NewDoubleArray(v []float64) *DoubleArray {
	m := new(DoubleArray)
	m.Val = v
	return m
}

func (this *DoubleArray) Length() int {
	if this.Val == nil {
		return 0
	}
	return len(this.Val)
}

func (this *DoubleArray) Get(i int) Value {
	return NewDoubleValue(this.Val[i])
}

func (this *DoubleArray) CompareTo(o Value) int {
	if o != nil && o.GetValueType() == this.GetValueType() {
		other := o.(*DoubleArray)
		if len(this.Val) != len(other.Val) {
			return len(this.Val) - len(other.Val)
		}
		for i := 0; i < len(this.Val); i++ {
			if this.Val[i] < other.Val[i] {
				return -1
			} else if this.Val[i] > other.Val[i] {
				return 1
			}
		}
		return 0
	}
	if o == nil {
		return 1
	}
	return int(this.GetValueType() - o.GetValueType())
}

func (this *DoubleArray) Equals(o Value) bool {
	if o != nil && o.GetValueType() == this.GetValueType() {
		other := o.(*DoubleArray)
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

func (this *DoubleArray) GetValueType() byte {
	return ARRAY_DOUBLE
}

func (this *DoubleArray) Write(out *io.DataOutputX) {
	out.WriteDecimal(int64(len(this.Val)))
	for _, v := range this.Val {
		out.WriteDouble(v)
	}
}

func (this *DoubleArray) Read(in *io.DataInputX) {
	count := int(in.ReadDecimal())
	this.Val = make([]float64, count)
	for i := 0; i < count; i++ {
		this.Val[i] = in.ReadDouble()
	}
}

func (this *DoubleArray) String() string {
	parts := make([]string, len(this.Val))
	for i, v := range this.Val {
		parts[i] = fmt.Sprintf("%f", v)
	}
	return "[" + strings.Join(parts, ", ") + "]"
}
