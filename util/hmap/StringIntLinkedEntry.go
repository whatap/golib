package hmap

import (
	"fmt"

	"github.com/whatap/golib/util/hash"
)

type StringIntLinkedEntry struct {
	key       string
	keyHash   uint
	value     int32
	hash_next *StringIntLinkedEntry
	link_next *StringIntLinkedEntry
	link_prev *StringIntLinkedEntry
}

func (this *StringIntLinkedEntry) GetKey() string {
	return this.key
}
func (this *StringIntLinkedEntry) GetValue() int32 {
	return this.value
}
func (this *StringIntLinkedEntry) SetValue(v int32) int32 {
	old := this.value
	this.value = v
	return old
}
func (this *StringIntLinkedEntry) Equals(o *StringIntLinkedEntry) bool {
	return this.key == o.key
}

func (this *StringIntLinkedEntry) HashCode() uint {
	return uint(hash.Hash([]byte(this.key)))
}

func (this *StringIntLinkedEntry) ToString() string {
	return fmt.Sprintf("%s=%v", this.key, this.value)
}
