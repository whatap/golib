package hmap

import (
	"github.com/whatap/golib/util/stringutil"
	//"fmt"
)

type StringLinkedSetry struct {
	key       string
	hash_next *StringLinkedSetry
	link_next *StringLinkedSetry
	link_prev *StringLinkedSetry
}

func (this *StringLinkedSetry) Get() string {
	return this.key
}

func (this *StringLinkedSetry) Equals(o *StringLinkedSetry) bool {
	return this.key == o.key
}

func (this *StringLinkedSetry) HashCode() uint {
	return uint(stringutil.HashCode(this.key))
}

func (this *StringLinkedSetry) ToString() string {
	return this.key
}
