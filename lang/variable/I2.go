package variable

import (
	"github.com/whatap/golib/io"
	"github.com/whatap/golib/util/compare"
	"github.com/whatap/golib/util/hmap"
)

// implements LinkedKey
type I2 struct {
	V1 int32
	V2 int32
}

func NewI2Defuault() *I2 {
	p := new(I2)

	return p
}

func NewI2(v1, v2 int32) *I2 {
	p := new(I2)
	p.V1 = v1
	p.V2 = v2

	return p
}

// override LinkedKey Hash() uint
func (this *I2) Hash() uint {

	prime := int32(31)
	result := int32(1)
	result = prime*result + this.V2
	result = prime*result + this.V1

	return uint(result)
}

// override LinkedKey Equals(h LinkedKey) bool
func (this *I2) Equals(h hmap.LinkedKey) bool {
	if this == h {
		return true
	}
	if h == nil {
		return false
	}
	//	if (getClass() != obj.getClass())
	//		return false;
	other := h.(*I2)
	return (this.V2 == other.V2 && this.V1 == other.V1)
}

func (this *I2) CompareTo(o *I2) int {
	if this.V1 != o.V1 {
		return compare.CompareToInt(int(this.V1), int(o.V1))
	}
	return compare.CompareToInt(int(this.V2), int(o.V2))
}

func (this *I2) ToBytes() []byte {
	b := make([]byte, 8)
	io.SetBytesInt(b, 0, this.V1)
	io.SetBytesInt(b, 4, this.V2)
	return b
}

func (this *I2) ToObject(b []byte) *I2 {
	this.V1 = io.ToInt(b, 0)
	this.V2 = io.ToInt(b, 4)
	return this
}
