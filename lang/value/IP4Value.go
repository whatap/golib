package value

import (
	"github.com/whatap/golib/io"
	"github.com/whatap/golib/util/compare"
	"github.com/whatap/golib/util/iputil"
)

type IP4Value struct {
	Val []byte
}

func NewIP4ValueString(v string) *IP4Value {
	return NewIP4Value(iputil.ToBytes(v))
}
func NewIP4Value(v []byte) *IP4Value {
	m := new(IP4Value)
	if v != nil && len(v) == 4 {
		m.Val = v
	} else {
		m.Val = []byte{0, 0, 0, 0}
	}
	return m
}

func (this *IP4Value) CompareTo(o Value) int {
	if o != nil && o.GetValueType() == this.GetValueType() {
		return compare.CompareToBytes(this.Val, o.(*IP4Value).Val)
	}
	if o == nil {
		return 1
	} else {
		return int(this.GetValueType() - o.GetValueType())
	}
}

func (this *IP4Value) Equals(o Value) bool {
	if o != nil && o.GetValueType() == this.GetValueType() {
		return compare.EqualBytes(this.Val, o.(*IP4Value).Val)
	}
	return false
}

func (this *IP4Value) GetValueType() byte {
	return VALUE_IP4ADDR
}
func (this *IP4Value) Write(out *io.DataOutputX) {
	out.WriteBytes(this.Val)
}
func (this *IP4Value) Read(in *io.DataInputX) {
	this.Val = in.ReadBytes(4)
}
