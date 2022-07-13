package hmap

import (
	"fmt"
)

type IntLinkedSetry struct {
	key       int32
	hash_next *IntLinkedSetry
	link_next *IntLinkedSetry
	link_prev *IntLinkedSetry
}

func (this *IntLinkedSetry) Get() int32 {
	return this.key
}

func (this *IntLinkedSetry) Equals(o *IntLinkedSetry) bool {
	return this.key == o.key
}

func (this *IntLinkedSetry) HashCode() uint {
	return uint(this.key)
}

func (this *IntLinkedSetry) ToString() string {
	return fmt.Sprintf("%d", this.key)
}
