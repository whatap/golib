package value

import (
	"github.com/whatap/golib/io"
)

type NullValue struct {
}

var NULL_VALUE *NullValue = new(NullValue)

func NewNullValue() *NullValue {
	return NULL_VALUE
}

func (this *NullValue) CompareTo(o Value) int {
	if o != nil && o.GetValueType() == this.GetValueType() {
		return 0
	}
	return 1
}

func (this *NullValue) Equals(o Value) bool {
	return o != nil && o.GetValueType() == this.GetValueType()
}

func (this *NullValue) GetValueType() byte {
	return VALUE_NULL
}
func (this *NullValue) Write(out *io.DataOutputX) {}
func (this *NullValue) Read(in *io.DataInputX)    {}
