package hmap

import (
	"fmt"
)

type LinkedKey interface {
	Hash() uint
	Equals(h LinkedKey) bool
}
type LinkedEntry struct {
	key       LinkedKey
	keyHash   uint
	value     interface{}
	hash_next *LinkedEntry
	link_next *LinkedEntry
	link_prev *LinkedEntry
}

func (this *LinkedEntry) GetKey() LinkedKey {
	return this.key
}
func (this *LinkedEntry) GetValue() interface{} {
	return this.value
}
func (this *LinkedEntry) SetValue(v interface{}) interface{} {
	old := this.value
	this.value = v
	return old
}
func (this *LinkedEntry) Equals(o *LinkedEntry) bool {
	return this.key.Equals(o.key)
}

func (this *LinkedEntry) HashCode() uint {
	return this.key.Hash()
}

func (this *LinkedEntry) ToString() string {
	return fmt.Sprintf("%v=%v", this.key, this.value)
}
