package hmap

import (
	"fmt"
)

type LongLongLinkedEntry struct {
	key       int64
	value     int64
	hash_next *LongLongLinkedEntry
	link_next *LongLongLinkedEntry
	link_prev *LongLongLinkedEntry
}
func NewLongLongLinkedEntry(key, value int64, next *LongLongLinkedEntry) *LongLongLinkedEntry {
	p := new(LongLongLinkedEntry)
	
	p.key = key
	p.value = value
	p.hash_next = next
	
	return p
}
func (this *LongLongLinkedEntry) GetKey() int64 {
	return this.key
}
func (this *LongLongLinkedEntry) GetValue() int64 {
	return this.value
}
func (this *LongLongLinkedEntry) SetValue(v int64) int64 {
	old := this.value
	this.value = v
	return old
}
func (this *LongLongLinkedEntry) Equals(o *LongLongLinkedEntry) bool {
	return this.key == o.key && this.value == o.value
}
func (this *LongLongLinkedEntry) HashCode() uint {
	return uint(this.key) ^ uint(this.value)
	//return (int) (key ^ (key >>> 32)) ^ (int) (value ^ (value >>> 32));
}
func (this *LongLongLinkedEntry) ToString() string {
	return fmt.Sprintf("%d=%d", this.key, this.value)
}

