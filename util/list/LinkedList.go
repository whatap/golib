package list

import (
	"bytes"
	"sync"
)

type LinkedList struct {
	size  int
	first *LinkedListEntity
	last  *LinkedListEntity
	lock  sync.Mutex
}

func NewLinkedList() *LinkedList {
	o := new(LinkedList)
	return o
}

func (o *LinkedList) AddFirst(v interface{}) {
	o.lock.Lock()
	defer o.lock.Unlock()

	newNode := &LinkedListEntity{prev: nil, Value: v, next: o.first}

	f := o.first
	o.first = newNode
	if f == nil {
		o.last = newNode
	} else {
		f.prev = newNode
	}
	o.size++
}

func (o *LinkedList) AddLast(v interface{}) {
	o.lock.Lock()
	defer o.lock.Unlock()

	newNode := &LinkedListEntity{prev: o.last, Value: v, next: nil}
	l := o.last
	o.last = newNode
	if l == nil {
		o.first = newNode
	} else {
		l.next = newNode
	}
	o.size++
}

func (o *LinkedList) PutBefore(v interface{}, succ *LinkedListEntity) *LinkedListEntity {
	o.lock.Lock()
	defer o.lock.Unlock()
	prev := succ.prev

	newNode := &LinkedListEntity{prev: prev, Value: v, next: succ}

	succ.prev = newNode
	if prev == nil {
		o.first = newNode
	} else {
		prev.next = newNode
	}
	o.size++
	return newNode
}

func (o *LinkedList) Remove(x *LinkedListEntity) interface{} {
	o.lock.Lock()
	defer o.lock.Unlock()
	return o.remove(x)
}

func (o *LinkedList) remove(x *LinkedListEntity) interface{} {

	v := x.Value

	if x.prev == nil {
		o.first = x.next
	} else {
		x.prev.next = x.next
	}
	if x.next == nil {
		o.last = x.prev
	} else {
		x.next.prev = x.prev
	}
	o.size--
	
	// avoid memory leaks
	x.next = nil
	x.prev = nil
	x.Value = nil
	x = nil
	
	return v
}

func (o *LinkedList) GetFirst() *LinkedListEntity {
	return o.first
}

func (o *LinkedList) GetLast() *LinkedListEntity {
	return o.last
}

func (o *LinkedList) GetNext(e *LinkedListEntity) *LinkedListEntity {
	return e.next
}

func (o *LinkedList) RemoveFirst() interface{} {
	o.lock.Lock()
	defer o.lock.Unlock()

	if o.first != nil {
		return o.remove(o.first)
	}
	return nil
}

func (o *LinkedList) RemoveLast() interface{} {
	o.lock.Lock()
	defer o.lock.Unlock()

	if o.last != nil {
		return o.remove(o.last)
	}
	return nil
}

func (o *LinkedList) Size() int {
	return o.size
}

func (o *LinkedList) Add(v interface{}) bool {
	o.AddLast(v)
	return true
}

func (o *LinkedList) Clear() {
	o.lock.Lock()
	defer o.lock.Unlock()

	o.first = nil
	o.last = nil
	o.size = 0
}

func (o *LinkedList) ToArray() []interface{} {
	o.lock.Lock()
	defer o.lock.Unlock()

	result := make([]interface{}, o.size)
	x := o.first
	for i := 0; i < o.size; i++ {
		result[i] = x.Value
		x = x.next
	}
	return result
}

func (o *LinkedList) ToString() string {
	o.lock.Lock()
	defer o.lock.Unlock()

	var buffer bytes.Buffer
	x := o.first
	for i := 0; i < o.size; i++ {
		if i > 0 {
			buffer.WriteString(",")
		}
		buffer.WriteString(x.ToString())
		x = x.next
	}
	return buffer.String()
}
