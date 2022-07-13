package hmap

import (
	"bytes"
	"fmt"
	"sort"
	"sync"
)

type LinkedSet struct {
	table      []*LinkedSetry
	header     *LinkedSetry
	count      int
	threshold  int
	loadFactor float32
	lock       sync.Mutex
	max        int
}

func NewLinkedSet() *LinkedSet {

	initCapacity := DEFAULT_CAPACITY
	loadFactor := DEFAULT_LOAD_FACTOR

	this := new(LinkedSet)
	this.loadFactor = float32(loadFactor)
	this.table = make([]*LinkedSetry, initCapacity)
	this.header = &LinkedSetry{}
	this.header.link_next = this.header
	this.header.link_prev = this.header
	this.threshold = int(float64(initCapacity) * loadFactor)
	return this
}

func (this *LinkedSet) Size() int {
	return this.count
}

func (this *LinkedSet) KeyArray() []LinkedKey {
	this.lock.Lock()
	defer this.lock.Unlock()

	_keys := make([]LinkedKey, this.Size())
	en := this.Keys()
	for i := 0; i < len(_keys); i++ {
		_keys[i] = en.NextElement().(LinkedKey)
	}
	return _keys
}

type EnumerSetImpl struct {
	parent *LinkedSet
	entry  *LinkedSetry
	rtype  int
}

func (this *EnumerSetImpl) HasMoreElements() bool {
	return this.entry != nil && this.parent.header != this.entry
}

func (this *EnumerSetImpl) NextElement() interface{} {
	if this.HasMoreElements() {
		e := this.entry
		this.entry = e.link_next
		return e.Get()
	}
	return nil
}
func (this *LinkedSet) Keys() Enumeration {
	return &EnumerSetImpl{parent: this, entry: this.header.link_next}
}

func (this *LinkedSet) Contains(key LinkedKey) bool {
	this.lock.Lock()
	defer this.lock.Unlock()

	tab := this.table
	index := this.hash(key) % uint(len(tab))
	for e := tab[index]; e != nil; e = e.hash_next {
		if e.key.Equals(key) {
			return true
		}
	}
	return false

}

func (this *LinkedSet) GetFirst() LinkedKey {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.header.link_next.key
}

func (this *LinkedSet) GetLast() LinkedKey {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.header.link_prev.key
}
func (this *LinkedSet) hash(key LinkedKey) uint {
	return key.Hash()
}

func (this *LinkedSet) rehash() {
	oldCapacity := len(this.table)
	oldMap := this.table
	newCapacity := oldCapacity*2 + 1
	newMap := make([]*LinkedSetry, newCapacity)
	this.threshold = int(float32(newCapacity) * this.loadFactor)
	this.table = newMap
	for i := oldCapacity; i > 0; i-- {
		for old := oldMap[i-1]; old != nil; {
			e := old
			old = old.hash_next
			index := uint(e.keyHash % uint(newCapacity))
			e.hash_next = newMap[index]
			newMap[index] = e
		}
	}
}

func (this *LinkedSet) SetMax(max int) *LinkedSet {
	this.max = max
	return this
}
func (this *LinkedSet) Put(key LinkedKey) interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.put(key, PUT_LAST)
}
func (this *LinkedSet) PutLast(key LinkedKey) interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.put(key, PUT_FORCE_LAST)
}
func (this *LinkedSet) PutFirst(key LinkedKey) interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.put(key, PUT_FORCE_FIRST)
}
func (this *LinkedSet) put(key LinkedKey, m PUT_MODE) interface{} {
	tab := this.table
	keyHash := this.hash(key)
	index := keyHash % uint(len(tab))
	for e := tab[index]; e != nil; e = e.hash_next {
		if e.key.Equals(key) {
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
	e := &LinkedSetry{key: key, keyHash: keyHash, hash_next: tab[index]}
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

func (this *LinkedSet) Remove(key LinkedKey) interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.remove(key)
}
func (this *LinkedSet) RemoveFirst() interface{} {
	if this.IsEmpty() {
		return 0
	}
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.remove(this.header.link_next.key)
}

func (this *LinkedSet) RemoveLast() interface{} {
	if this.IsEmpty() {
		return 0
	}
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.remove(this.header.link_prev.key)
}

func (this *LinkedSet) remove(key LinkedKey) interface{} {

	tab := this.table
	index := this.hash(key) % uint(len(tab))
	e := tab[index]
	var prev *LinkedSetry = nil
	for e != nil {
		if e.key.Equals(key) {
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

func (this *LinkedSet) IsEmpty() bool {
	return this.count == 0
}
func (this *LinkedSet) IsFull() bool {
	return this.max > 0 && this.max <= this.count
}

func (this *LinkedSet) Clear() {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.clear()
}
func (this *LinkedSet) clear() {
	tab := this.table
	for index := len(tab) - 1; index >= 0; index-- {
		tab[index] = nil
	}
	this.header.link_next = this.header
	this.header.link_prev = this.header
	this.count = 0
}

func (this *LinkedSet) chain(link_prev *LinkedSetry, link_next *LinkedSetry, e *LinkedSetry) {
	e.link_prev = link_prev
	e.link_next = link_next
	link_prev.link_next = e
	link_next.link_prev = e
}

func (this *LinkedSet) unchain(e *LinkedSetry) {
	e.link_prev.link_next = e.link_next
	e.link_next.link_prev = e.link_prev
	e.link_prev = nil
	e.link_next = nil
}

func (this *LinkedSet) ToString() string {
	this.lock.Lock()
	defer this.lock.Unlock()

	var buffer bytes.Buffer
	x := this.Keys()
	buffer.WriteString("{")
	for i := 0; x.HasMoreElements(); i++ {
		if i > 0 {
			buffer.WriteString(", ")
		}
		e := x.NextElement().(LinkedKey)
		buffer.WriteString(fmt.Sprintf("%v", e))
	}
	buffer.WriteString("}")
	return buffer.String()
}

type setSortable struct {
	compare func(a, b LinkedKey) bool
	data    []*LinkedSetry
}

func (this setSortable) Len() int {
	return len(this.data)
}
func (this setSortable) Less(i, j int) bool {
	return this.compare(this.data[i].Get(), this.data[j].Get())
}

func (this setSortable) Swap(i, j int) {
	this.data[i], this.data[j] = this.data[j], this.data[i]
}

func (this *LinkedSet) Sort(c func(k1, k2 LinkedKey) bool) {
	this.lock.Lock()
	defer this.lock.Unlock()

	sz := this.Size()
	list := make([]*LinkedSetry, sz)
	en := this.Keys()
	for i := 0; i < sz; i++ {
		list[i] = en.NextElement().(*LinkedSetry)
	}
	sort.Sort(setSortable{compare: c, data: list})

	this.clear()
	for i := 0; i < sz; i++ {
		this.put(list[i].Get(), PUT_LAST)
	}
}
