package hmap

import (
	"fmt"
	"math"
)

type LongFloatLinkedEntry struct {
	key       int64
	value     float32
	hash_next *LongFloatLinkedEntry
	link_next *LongFloatLinkedEntry
	link_prev *LongFloatLinkedEntry
}

func (this *LongFloatLinkedEntry) GetKey() int64 {
	return this.key
}
func (this *LongFloatLinkedEntry) GetValue() float32 {
	return this.value
}
func (this *LongFloatLinkedEntry) SetValue(v float32) float32 {
	old := this.value
	this.value = v
	return old
}
func (this *LongFloatLinkedEntry) Equals(o *LongFloatLinkedEntry) bool {
	return this.key == o.key && this.value == o.value
}
func (this *LongFloatLinkedEntry) HashCode() uint {
	return uint(this.key) ^ uint(math.Float32bits(this.value))
}
func (this *LongFloatLinkedEntry) ToString() string {
	return fmt.Sprintf("%d=%f", this.key, this.value)
}
