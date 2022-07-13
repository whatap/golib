package hmap

import (
	"fmt"
	"math"
)

type IntFloatLinkedEntry struct {
	key       int32
	value     float32
	hash_next *IntFloatLinkedEntry
	link_next *IntFloatLinkedEntry
	link_prev *IntFloatLinkedEntry
}

func (this *IntFloatLinkedEntry) GetKey() int32 {
	return this.key
}
func (this *IntFloatLinkedEntry) GetValue() float32 {
	return this.value
}
func (this *IntFloatLinkedEntry) SetValue(v float32) float32 {
	old := this.value
	this.value = v
	return old
}
func (this *IntFloatLinkedEntry) Equals(o *IntFloatLinkedEntry) bool {
	return this.key == o.key && this.value == o.value
}
func (this *IntFloatLinkedEntry) HashCode() uint {
	return uint(this.key) ^ uint(math.Float32bits(this.value))
}
func (this *IntFloatLinkedEntry) ToString() string {
	return fmt.Sprintf("%d=%f", this.key, this.value)
}
