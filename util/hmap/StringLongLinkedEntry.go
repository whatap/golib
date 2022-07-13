package hmap

import (
	"fmt"

	"github.com/whatap/golib/util/hash"
)

type StringLongLinkedEntry struct {
	key       string
	keyHash   uint
	value     int64
	hash_next *StringLongLinkedEntry
	link_next *StringLongLinkedEntry
	link_prev *StringLongLinkedEntry
}

func (this *StringLongLinkedEntry) GetKey() string {
	return this.key
}
func (this *StringLongLinkedEntry) GetValue() int64 {
	return this.value
}
func (this *StringLongLinkedEntry) SetValue(v int64) int64 {
	old := this.value
	this.value = v
	return old
}
func (this *StringLongLinkedEntry) Equals(o *StringLongLinkedEntry) bool {
	return this.key == o.key
}

func (this *StringLongLinkedEntry) HashCode() uint {
	return uint(hash.Hash([]byte(this.key)))
}

func (this *StringLongLinkedEntry) ToString() string {
	return fmt.Sprintf("%s=%v", this.key, this.value)
}
