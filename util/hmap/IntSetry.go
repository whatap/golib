package hmap

import (
//"fmt"
	"strconv"
)

type IntSetry struct {
	key  int32
	next *IntSetry
}

func NewIntSetry(key int32, next *IntSetry) *IntSetry {
	p := new(IntSetry)
	p.key = key
	p.next = next
	return p
}
func (this *IntSetry) Clone() *IntSetry {
	var n *IntSetry
	if this.next == nil {
		n = nil
	} else {
		n = this.next.Clone()
	}
	return NewIntSetry(this.key, n)
}

func (this *IntSetry) GetKey() int32 {
	return this.key
}

func (this *IntSetry) Get() int32 {
	return this.key
}

//func (this *IntSetry) Equals(o interface{}) bool{
//	var e *IntSetry
//	switch o.(type){
//		case IntSetry:
//			e = o.(*IntSetry)
//		default:
//			return false
//	}
//	if e.GetKey() == this.key {
//		return true
//	}else {
//		return false
//	}
//}

func (this *IntSetry) Equals(o *IntSetry) bool {
	return this.key == o.key
}

func (this *IntSetry) HashCode() int32 {
	return this.key
}

func (this *IntSetry) ToString() string {
	return strconv.Itoa(int(this.key))
}
