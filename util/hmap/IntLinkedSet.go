package hmap

import (
	"bytes"
	"fmt"
	"sort"
	"sync"
)

type IntLinkedSet struct {
	table      []*IntLinkedSetry
	header     *IntLinkedSetry
	count      int
	threshold  int
	loadFactor float32
	lock       sync.Mutex
	max        int
}

func NewIntLinkedSet() *IntLinkedSet {

	initCapacity := DEFAULT_CAPACITY
	loadFactor := DEFAULT_LOAD_FACTOR

	this := new(IntLinkedSet)
	this.loadFactor = float32(loadFactor)
	this.table = make([]*IntLinkedSetry, initCapacity)
	this.header = &IntLinkedSetry{}
	this.header.link_next = this.header
	this.header.link_prev = this.header
	this.threshold = int(float64(initCapacity) * loadFactor)
	return this
}

func (this *IntLinkedSet) Size() int {
	return this.count
}

func (this *IntLinkedSet) KeyArray() []int32 {
	this.lock.Lock()
	defer this.lock.Unlock()

	_keys := make([]int32, this.Size())
	en := this.Keys()
	for i := 0; i < len(_keys); i++ {
		_keys[i] = en.NextInt()
	}
	return _keys
}

type IntEnumerSetImpl struct {
	parent *IntLinkedSet
	entry  *IntLinkedSetry
	rtype  int
}

func (this *IntEnumerSetImpl) HasMoreElements() bool {
	return this.entry != nil && this.parent.header != this.entry
}

func (this *IntEnumerSetImpl) NextInt() int32 {
	if this.HasMoreElements() {
		e := this.entry
		this.entry = e.link_next
		return e.Get()
	}
	return 0
}
func (this *IntLinkedSet) Keys() IntEnumer {
	return &IntEnumerSetImpl{parent: this, entry: this.header.link_next}
}

func (this *IntLinkedSet) Contains(key int32) bool {
	this.lock.Lock()
	defer this.lock.Unlock()

	tab := this.table
	index := this.hash(key) % uint(len(tab))
	for e := tab[index]; e != nil; e = e.hash_next {
		if e.key==key {
			return true
		}
	}
	return false

}

func (this *IntLinkedSet) GetFirst() int32 {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.header.link_next.key
}

func (this *IntLinkedSet) GetLast() int32 {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.header.link_prev.key
}
func (this *IntLinkedSet) hash(key int32) uint {
	return uint(key)
}

func (this *IntLinkedSet) rehash() {
	oldCapacity := len(this.table)
	oldMap := this.table
	newCapacity := oldCapacity*2 + 1
	newMap := make([]*IntLinkedSetry, newCapacity)
	this.threshold = int(float32(newCapacity) * this.loadFactor)
	this.table = newMap
	for i := oldCapacity; i > 0; i-- {
		for old := oldMap[i-1]; old != nil; {
			e := old
			old = old.hash_next
			index := uint(this.hash(e.key) % uint(newCapacity))
			e.hash_next = newMap[index]
			newMap[index] = e
		}
	}
}

func (this *IntLinkedSet) SetMax(max int) *IntLinkedSet {
	this.max = max
	return this
}
func (this *IntLinkedSet) Put(key int32) interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.put(key, PUT_LAST)
}
func (this *IntLinkedSet) PutLast(key int32) interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.put(key, PUT_FORCE_LAST)
}
func (this *IntLinkedSet) PutFirst(key int32) interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.put(key, PUT_FORCE_FIRST)
}
func (this *IntLinkedSet) put(key int32, m PUT_MODE) interface{} {
	tab := this.table
	keyHash := this.hash(key)
	index := keyHash % uint(len(tab))
	for e := tab[index]; e != nil; e = e.hash_next {
		if e.key==key {
			switch m {
			case PUT_FORCE_FIRST:
				if this.header.link_next != e {
					this.unchain(e)
					this.chain(this.header, this.header.link_next, e)
				}
			case PUT_FORCE_LAST:
				if this.header.link_prev != e {
					this.unchain(e)
					this.chain(this.header.link_prev, this.header, e)
				}
			}
			return key
		}
	}
	if this.max > 0 {
		switch m {
		case PUT_FORCE_FIRST, PUT_FIRST:
			for this.count >= this.max {
				k := this.header.link_prev.key
				this.remove(k)
			}
		case PUT_FORCE_LAST, PUT_LAST:
			for this.count >= this.max {
				k := this.header.link_next.key
				this.remove(k)
			}
			break
		}
	}
	if this.count >= this.threshold {
		this.rehash()
		tab = this.table
		index = keyHash % uint(len(tab))
	}
	e := &IntLinkedSetry{key: key,  hash_next: tab[index]}
	tab[index] = e
	switch m {
	case PUT_FORCE_FIRST, PUT_FIRST:
		this.chain(this.header, this.header.link_next, e)
	case PUT_FORCE_LAST, PUT_LAST:
		this.chain(this.header.link_prev, this.header, e)
	}
	this.count++
	return nil
}

func (this *IntLinkedSet) Remove(key int32) interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.remove(key)
}
func (this *IntLinkedSet) RemoveFirst() interface{} {
	if this.IsEmpty() {
		return 0
	}
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.remove(this.header.link_next.key)
}

func (this *IntLinkedSet) RemoveLast() interface{} {
	if this.IsEmpty() {
		return 0
	}
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.remove(this.header.link_prev.key)
}

func (this *IntLinkedSet) remove(key int32) interface{} {

	tab := this.table
	index := this.hash(key) % uint(len(tab))
	e := tab[index]
	var prev *IntLinkedSetry = nil
	for e != nil {
		if e.key==key {
			if prev != nil {
				prev.hash_next = e.hash_next
			} else {
				tab[index] = e.hash_next
			}
			this.count--
			//
			this.unchain(e)
			return key
		}
		prev = e
		e = e.hash_next
	}
	return nil
}

func (this *IntLinkedSet) IsEmpty() bool {
	return this.count == 0
}
func (this *IntLinkedSet) IsFull() bool {
	return this.max > 0 && this.max <= this.count
}

func (this *IntLinkedSet) Clear() {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.clear()
}
func (this *IntLinkedSet) clear() {
	tab := this.table
	for index := len(tab) - 1; index >= 0; index-- {
		tab[index] = nil
	}
	this.header.link_next = this.header
	this.header.link_prev = this.header
	this.count = 0
}

func (this *IntLinkedSet) chain(link_prev *IntLinkedSetry, link_next *IntLinkedSetry, e *IntLinkedSetry) {
	e.link_prev = link_prev
	e.link_next = link_next
	link_prev.link_next = e
	link_next.link_prev = e
}

func (this *IntLinkedSet) unchain(e *IntLinkedSetry) {
	e.link_prev.link_next = e.link_next
	e.link_next.link_prev = e.link_prev
	e.link_prev = nil
	e.link_next = nil
}

func (this *IntLinkedSet) ToString() string {
	this.lock.Lock()
	defer this.lock.Unlock()

	var buffer bytes.Buffer
	x := this.Keys()
	buffer.WriteString("{")
	for i := 0; x.HasMoreElements(); i++ {
		if i > 0 {
			buffer.WriteString(", ")
		}
		e := x.NextInt()
		buffer.WriteString(fmt.Sprintf("%d", e))
	}
	buffer.WriteString("}")
	return buffer.String()
}

type setIntSortable struct {
	compare func(a, b int32) bool
	data    []int32
}

func (this setIntSortable) Len() int {
	return len(this.data)
}
func (this setIntSortable) Less(i, j int) bool {
	return this.compare(this.data[i], this.data[j])
}

func (this setIntSortable) Swap(i, j int) {
	this.data[i], this.data[j] = this.data[j], this.data[i]
}

func (this *IntLinkedSet) Sort(c func(k1, k2 int32) bool) {
	this.lock.Lock()
	defer this.lock.Unlock()

	sz := this.Size()
	list := make([]int32, sz)
	en := this.Keys()
	for i := 0; i < sz; i++ {
		list[i] = en.NextInt()
	}
	sort.Sort(setIntSortable{compare: c, data: list})

	this.clear()
	for i := 0; i < sz; i++ {
		this.put(list[i], PUT_LAST)
	}
}
