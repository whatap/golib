package variable

import (
	"github.com/whatap/golib/io"
	"github.com/whatap/golib/util/compare"
	"github.com/whatap/golib/util/hmap"
)

// implements LinkedKey
type L2 struct {
	V1 int64
	V2 int64
}

func NewL2Defuault() *L2 {
	p := new(L2)

	return p
}

func NewL2(v1, v2 int64) *L2 {
	p := new(L2)
	p.V1 = v1
	p.V2 = v2

	return p
}

// override LinkedKey Hash() uint
func (this *L2) Hash() uint {

	prime := int64(31)
	result := int64(1)
	result = prime*result + this.V1 ^ int64(uint64(this.V1)>>32)
	result = prime*result + this.V2 ^ int64(uint64(this.V2)>>32)

	return uint(result)
}

// override LinkedKey Equals(h LinkedKey) bool
func (this *L2) Equals(h hmap.LinkedKey) bool {
	if this == h {
		return true
	}
	if h == nil {
		return false
	}
	//	if (getClass() != obj.getClass())
	//		return false;
	other := h.(*L2)
	return (this.V1 == other.V1 && this.V2 == other.V2)
}

func (this *L2) CompareTo(o *L2) int {
	if this.V1 != o.V1 {
		return compare.CompareToLong(this.V1, o.V1)
	}
	return compare.CompareToLong(this.V2, o.V2)
}

func (this *L2) ToBytes() []byte {
	b := make([]byte, 16)
	io.SetBytesLong(b, 0, this.V1)
	io.SetBytesLong(b, 8, this.V2)
	return b
}

func (this *L2) ToObject(b []byte) *L2 {
	this.V1 = io.ToLong(b, 0)
	this.V2 = io.ToLong(b, 8)
	return this
}
