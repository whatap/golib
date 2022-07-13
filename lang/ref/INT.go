package ref

import ()

type INT struct {
	Value int
}

func NewINT() *INT {
	p := new(INT)
	return p
}

func (this *INT) HashCode() int32 {
	return int32(this.Value)
}

func (this *INT) Equals(obj *INT) bool {
	if obj != nil {
		other := obj
		return this.Value == other.Value
	}
	return false
}
