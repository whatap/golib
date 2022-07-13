package value

import (
	"github.com/whatap/golib/io"
	"github.com/whatap/golib/util/compare"
)

type BlobValue struct {
	Val []byte
}

func NewBlobValue(v []byte) *BlobValue {
	m := new(BlobValue)
	m.Val = v
	return m
}

func (this *BlobValue) CompareTo(o Value) int {
	if o != nil && o.GetValueType() == this.GetValueType() {
		return compare.CompareToBytes(this.Val, o.(*BlobValue).Val)
	}
	if o == nil {
		return 1
	} else {
		return int(this.GetValueType() - o.GetValueType())
	}
}

func (this *BlobValue) Equals(o Value) bool {
	if o != nil && o.GetValueType() == this.GetValueType() {
		return compare.EqualBytes(this.Val, o.(*BlobValue).Val)
	}
	return false
}

func (this *BlobValue) GetValueType() byte {
	return VALUE_BLOB
}
func (this *BlobValue) Write(out *io.DataOutputX) {
	out.WriteBlob(this.Val)
}
func (this *BlobValue) Read(in *io.DataInputX) {
	this.Val = in.ReadBlob()
}
