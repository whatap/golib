package variable

import (
	"github.com/whatap/golib/io"
	"github.com/whatap/golib/util/compare"
	"github.com/whatap/golib/util/hmap"
)

// implements LinkedKey
type I3 struct {
	V1 int32
	V2 int32
	V3 int32
}

func NewI3Defuault() *I3 {
	p := new(I3)

	return p
}

func NewI3(v1, v2, v3 int32) *I3 {
	p := new(I3)
	p.V1 = v1
	p.V2 = v2
	p.V3 = v3

	return p
}

// override LinkedKey Hash() uint
func (this *I3) Hash() uint {

	prime := int32(31)
	result := int32(1)
	result = prime*result + this.V1
	result = prime*result + this.V2
	result = prime*result + this.V3
	return uint(result)
}

// override LinkedKey Equals(h LinkedKey) bool
func (this *I3) Equals(h hmap.LinkedKey) bool {
	if this == h {
		return true
	}
	if h == nil {
		return false
	}
	//	if (getClass() != obj.getClass())
	//		return false;
	other := h.(*I3)
	return (this.V2 == other.V2 && this.V3 == other.V3 && this.V1 == other.V1)
}

func (this *I3) CompareTo(o *I3) int {
	if this.V1 != o.V1 {
		return compare.CompareToInt(int(this.V1), int(o.V1))
	}
	if this.V2 != o.V2 {
		return compare.CompareToInt(int(this.V2), int(o.V2))
	}
	return compare.CompareToInt(int(this.V3), int(o.V3))
}

func (this *I3) ToBytes() []byte {
	b := make([]byte, 12)
	io.SetBytesInt(b, 0, this.V1)
	io.SetBytesInt(b, 4, this.V2)
	io.SetBytesInt(b, 8, this.V2)
	return b
}

func (this *I3) ToObject(b []byte) *I3 {
	this.V1 = io.ToInt(b, 0)
	this.V2 = io.ToInt(b, 4)
	this.V3 = io.ToInt(b, 8)
	return this
}
