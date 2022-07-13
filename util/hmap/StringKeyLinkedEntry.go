package hmap

import (
	"fmt"

	"github.com/whatap/golib/util/hash"
)

type StringKeyLinkedEntry struct {
	key       string
	keyHash   uint
	value     interface{}
	hash_next *StringKeyLinkedEntry
	link_next *StringKeyLinkedEntry
	link_prev *StringKeyLinkedEntry
}

func (this *StringKeyLinkedEntry) GetKey() string {
	return this.key
}
func (this *StringKeyLinkedEntry) GetValue() interface{} {
	return this.value
}
func (this *StringKeyLinkedEntry) SetValue(v interface{}) interface{} {
	old := this.value
	this.value = v
	return old
}
func (this *StringKeyLinkedEntry) Equals(o *StringKeyLinkedEntry) bool {
	return this.key == o.key
}

func (this *StringKeyLinkedEntry) HashCode() uint {
	return uint(hash.Hash([]byte(this.key)))
}

func (this *StringKeyLinkedEntry) ToString() string {
	return fmt.Sprintf("%s=%v", this.key, this.value)
}
