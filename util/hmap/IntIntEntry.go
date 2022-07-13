package hmap

import (
	"fmt"
)

type IntIntEntry struct {
	key       int32
	value     int32
	next *IntIntEntry
}

func (this *IntIntEntry) GetKey() int32 {
	return this.key
}
func (this *IntIntEntry) GetValue() int32 {
	return this.value
}
func (this *IntIntEntry) SetValue(v int32) int32 {
	old := this.value
	this.value = v
	return old
}
func (this *IntIntEntry) Equals(o *IntIntEntry) bool {
	return this.key == o.key && this.value == o.value
}
func (this *IntIntEntry) HashCode() uint {
	return uint(this.key) ^ uint(this.value)
}
func (this *IntIntEntry) ToString() string {
	return fmt.Sprintf("%d=%d", this.key, this.value)
}
