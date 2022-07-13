package hmap

import (
	//"log"
	"fmt"
)

type IntKeyEntry struct {
	Key int32
	Value interface{}
	Next *IntKeyEntry
}

func NewIntKeyEntry(key int32, value interface{}, next *IntKeyEntry) *IntKeyEntry {
	p := new(IntKeyEntry)
	p.Key = key
	p.Value = value
	p.Next = next
	
	return p
}

func (this * IntKeyEntry) clone() *IntKeyEntry {
	if this.Next == nil {
		return NewIntKeyEntry(this.Key, this.Value, nil) 
	} else {
		return NewIntKeyEntry(this.Key, this.Value, this.Next.clone())
	}
}

func (this * IntKeyEntry) GetKey() int32 {
	return this.Key
}

func (this * IntKeyEntry) GetValue() interface{} {
	return this.Value
}

func (this * IntKeyEntry) SetValue(value interface{}) interface{} {
	defer func() {
		if r := recover(); r != nil {
			// TODO 추후 hmap 에서 recover 는 없애고 panic 처리. 호출하는 쪽에서 recover 할 것
			//logutil.Println("WA821", r)
		}
	}()
	
	if value == nil {
		panic("Error value is Nil")
	}
	
	oldValue := this.Value
	this.Value = value
	
	return oldValue
}

// *hmap.IntKeyEntry
//func (this * IntKeyEntry) Equals(o interface{}) bool{
//	
//	// type assert
//	e, ok := o.(*IntKeyEntry)
//	if ok {
//		
//	}  else {
//		return false
//	}
//	//return
//	//return (this.key == e.getKey()) && (this.value == nil ? e.getValue() == null : value.equals(e.getValue()))
//	if( this.Value == nil ) {
//		return (this.Key == e.GetKey()) && (e.GetValue() == nil)
//	} else {
//		return (this.Key == e.GetKey()) && (this.Value.Equals(e.GetValue()))	
//	}
//}

func (this *IntKeyEntry) Equals(o *IntKeyEntry) bool{
	return this.Key == o.Key //&& this.Value == o.Value
}

func (this * IntKeyEntry) HashCode() int32 {
//	if this.Value == nil {
//		return this.Key ^ 0 
//	} else {
//		return this.Key ^ this.value.hashCode()
//	}
	//return key ^ (value == null ? 0 : value.hashCode())
	return this.Key 
}

func (this * IntKeyEntry) ToString() string{
			return fmt.Sprintf("%d=%v", this.Key, this.Value)
		}
