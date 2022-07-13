package hmap

import (
	"fmt"
)

type LongKeyLinkedEntry struct {
	key       int64
	keyHash   uint
	value     interface{}
	hash_next *LongKeyLinkedEntry
	link_next *LongKeyLinkedEntry
	link_prev *LongKeyLinkedEntry
}

func (this *LongKeyLinkedEntry) GetKey() int64 {
	return this.key
}
func (this *LongKeyLinkedEntry) GetValue() interface{} {
	return this.value
}
func (this *LongKeyLinkedEntry) SetValue(v interface{}) interface{} {
	old := this.value
	this.value = v
	return old
}
func (this *LongKeyLinkedEntry) Equals(o *LongKeyLinkedEntry) bool {
	return this.key == o.key
}

func (this *LongKeyLinkedEntry) HashCode() uint {
	return uint(this.key ^ this.key>>32)
}

func (this *LongKeyLinkedEntry) ToString() string {
	return fmt.Sprintf("%d=%v", this.key, this.value)
}
