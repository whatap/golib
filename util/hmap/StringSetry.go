package hmap

import (
	"github.com/whatap/golib/util/stringutil"
	//"fmt"
)

type StringSetry struct {
	hash uint
	key  string
	next *StringSetry
}

func NewStringSetry(hash uint, key string, next *StringSetry) *StringSetry {
	p := new(StringSetry)
	p.hash = hash
	p.key = key
	p.next = next
	return p
}
func (this *StringSetry) Clone() *StringSetry {
	var n *StringSetry
	if this.next == nil {
		n = nil
	} else {
		n = this.next.Clone()
	}
	return NewStringSetry(this.hash, this.key, n)
}

func (this *StringSetry) GetKey() string {
	return this.key
}

func (this *StringSetry) Get() string {
	return this.key
}

//func (this *StringSetry) Equals(o interface{}) bool{
//	var e *StringSetry
//	switch o.(type){
//		case StringSetry:
//			e = o.(*StringSetry)
//		default:
//			return false
//	}
//	if e.GetKey() == this.key {
//		return true
//	}else {
//		return false
//	}
//}

func (this *StringSetry) Equals(o *StringSetry) bool {
	return this.key == o.key
}

func (this *StringSetry) HashCode() uint {
	return uint(stringutil.HashCode(this.key))
}

func (this *StringSetry) ToString() string {
	return this.key
}
