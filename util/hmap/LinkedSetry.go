package hmap

import (
	"fmt"
)

type LinkedSetry struct {
	key       LinkedKey
	keyHash   uint
	hash_next *LinkedSetry
	link_next *LinkedSetry
	link_prev *LinkedSetry
}

func (this *LinkedSetry) Get() LinkedKey {
	return this.key
}

func (this *LinkedSetry) Equals(o *LinkedSetry) bool {
	return this.key.Equals(o.key)
}

func (this *LinkedSetry) HashCode() uint {
	return this.key.Hash()
}

func (this *LinkedSetry) ToString() string {
	return fmt.Sprintf("%v", this.key)
}
