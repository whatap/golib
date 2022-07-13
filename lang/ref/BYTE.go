package ref

import ()

type BYTE struct {
	Value byte
}

func NewBYTE() *BYTE {
	p := new(BYTE)
	return p
}

func (this *BYTE) HashCode() int32 {
	return int32(this.Value)
}

func (this *BYTE) Equals(obj *BYTE) bool {
	if obj != nil {
		other := obj
		return this.Value == other.Value
	}
	return false
}
