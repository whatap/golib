package hmap

import (
	"fmt"
)

type IntIntLinkedEntry struct {
	key       int32
	value     int32
	hash_next *IntIntLinkedEntry
	link_next *IntIntLinkedEntry
	link_prev *IntIntLinkedEntry
}

func (this *IntIntLinkedEntry) GetKey() int32 {
	return this.key
}
func (this *IntIntLinkedEntry) GetValue() int32 {
	return this.value
}
func (this *IntIntLinkedEntry) SetValue(v int32) int32 {
	old := this.value
	this.value = v
	return old
}
func (this *IntIntLinkedEntry) Equals(o *IntIntLinkedEntry) bool {
	return this.key == o.key && this.value == o.value
}
func (this *IntIntLinkedEntry) HashCode() uint {
	return uint(this.key) ^ uint(this.value)
}
func (this *IntIntLinkedEntry) ToString() string {
	return fmt.Sprintf("%d=%d", this.key, this.value)
}
