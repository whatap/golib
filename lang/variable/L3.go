package variable

import (
	"github.com/whatap/golib/io"
	"github.com/whatap/golib/util/compare"
	"github.com/whatap/golib/util/hmap"
)

// implements LinkedKey
type L3 struct {
	V1 int64
	V2 int64
	V3 int64
}

func NewL3Defuault() *L3 {
	p := new(L3)

	return p
}

func NewL3(v1, v2, v3 int64) *L3 {
	p := new(L3)
	p.V1 = v1
	p.V2 = v2
	p.V3 = v3

	return p
}

// override LinkedKey Hash() uint
func (this *L3) Hash() uint {

	prime := int64(31)
	result := int64(1)
	result = prime*result + this.V1 ^ int64(uint64(this.V1)>>32)
	result = prime*result + this.V2 ^ int64(uint64(this.V2)>>32)
	result = prime*result + this.V3 ^ int64(uint64(this.V3)>>32)

	return uint(result)
}

// override LinkedKey Equals(h LinkedKey) bool
func (this *L3) Equals(h hmap.LinkedKey) bool {
	if this == h {
		return true
	}
	if h == nil {
		return false
	}
	//	if (getClass() != obj.getClass())
	//		return false;
	other := h.(*L3)
	return (this.V1 == other.V1 && this.V2 == other.V2 && this.V3 == other.V3)
}

func (this *L3) CompareTo(o *L3) int {
	if this.V1 != o.V1 {
		return compare.CompareToLong(this.V1, o.V1)
	}
	if this.V2 != o.V2 {
		return compare.CompareToLong(this.V2, o.V2)
	}
	return compare.CompareToLong(this.V3, o.V3)
}

func (this *L3) ToBytes() []byte {
	b := make([]byte, 24)
	io.SetBytesLong(b, 0, this.V1)
	io.SetBytesLong(b, 8, this.V2)
	io.SetBytesLong(b, 16, this.V3)
	return b
}

func (this *L3) ToObject(b []byte) *L3 {
	this.V1 = io.ToLong(b, 0)
	this.V2 = io.ToLong(b, 8)
	this.V3 = io.ToLong(b, 16)
	return this
}
