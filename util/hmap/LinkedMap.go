package hmap

import (
	"bytes"
	"sort"
	"sync"
)

type LinkedMap struct {
	table      []*LinkedEntry
	header     *LinkedEntry
	count      int
	threshold  int
	loadFactor float32
	lock       sync.Mutex
	max        int
}

func NewLinkedMapDefault() *LinkedMap {
	return NewLinkedMap(DEFAULT_CAPACITY, DEFAULT_LOAD_FACTOR)
}

func NewLinkedMap(initCapacity int, loadFactor float32) *LinkedMap {

	this := new(LinkedMap)
	this.loadFactor = float32(loadFactor)
	this.table = make([]*LinkedEntry, initCapacity)
	this.header = &LinkedEntry{}
	this.header.link_next = this.header
	this.header.link_prev = this.header
	this.threshold = int(float32(initCapacity) * loadFactor)
	return this
}

func (this *LinkedMap) Size() int {
	return this.count
}

func (this *LinkedMap) KeyArray() []LinkedKey {
	this.lock.Lock()
	defer this.lock.Unlock()

	_keys := make([]LinkedKey, this.Size())
	en := this.Keys()
	for i := 0; i < len(_keys); i++ {
		_keys[i] = en.NextElement().(LinkedKey)
	}
	return _keys
}

type EnumerImpl struct {
	parent *LinkedMap
	entry  *LinkedEntry
	rtype  int
}

func (this *EnumerImpl) HasMoreElements() bool {
	return this.entry != nil && this.parent.header != this.entry
}

func (this *EnumerImpl) NextElement() interface{} {
	if this.HasMoreElements() {
		e := this.entry
		this.entry = e.link_next
		switch this.rtype {
		case 1:
			return e.GetKey()
		case 2:
			return e.GetValue()
		case 3:
			return e
		}
	}
	return ""
}
func (this *LinkedMap) Keys() Enumeration {
	return &EnumerImpl{parent: this, entry: this.header.link_next, rtype: 1}
}
func (this *LinkedMap) Values() Enumeration {
	return &EnumerImpl{parent: this, entry: this.header.link_next, rtype: 2}
}
func (this *LinkedMap) Entries() Enumeration {
	return &EnumerImpl{parent: this, entry: this.header.link_next, rtype: 3}
}

func (this *LinkedMap) ContainsKey(key LinkedKey) bool {
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
func (this *LinkedMap) Get(key LinkedKey) interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()

	tab := this.table
	index := this.hash(key) % uint(len(tab))
	for e := tab[index]; e != nil; e = e.hash_next {
		if e.key.Equals(key) {
			return e.value
		}
	}
	return nil
}
func (this *LinkedMap) GetFirstKey() LinkedKey {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.header.link_next.key
}

func (this *LinkedMap) GetLastKey() LinkedKey {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.header.link_prev.key
}

func (this *LinkedMap) GetFirstValue() interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.header.link_next.value
}

func (this *LinkedMap) GetLastValue() interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.header.link_prev.value
}
func (this *LinkedMap) hash(key LinkedKey) uint {
	return key.Hash()
}

func (this *LinkedMap) rehash() {
	oldCapacity := len(this.table)
	oldMap := this.table
	newCapacity := oldCapacity*2 + 1
	newMap := make([]*LinkedEntry, newCapacity)
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

func (this *LinkedMap) SetMax(max int) *LinkedMap {
	this.max = max
	return this
}
func (this *LinkedMap) Put(key LinkedKey, value interface{}) interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.put(key, value, PUT_LAST)
}
func (this *LinkedMap) PutLast(key LinkedKey, value interface{}) interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.put(key, value, PUT_FORCE_LAST)
}
func (this *LinkedMap) PutFirst(key LinkedKey, value interface{}) interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.put(key, value, PUT_FORCE_FIRST)
}
func (this *LinkedMap) put(key LinkedKey, value interface{}, m PUT_MODE) interface{} {

	tab := this.table
	keyHash := this.hash(key)
	index := keyHash % uint(len(tab))
	for e := tab[index]; e != nil; e = e.hash_next {
		if e.key.Equals(key) {
			old := e.value
			e.value = value
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
			return old
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
				// removeFirst();
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
	e := &LinkedEntry{key: key, keyHash: keyHash, value: value, hash_next: tab[index]}
	tab[index] = e
	switch m {
	case PUT_FORCE_FIRST, PUT_FIRST:
		this.chain(this.header, this.header.link_next, e)
	case PUT_FORCE_LAST, PUT_LAST:
		this.chain(this.header.link_prev, this.header, e)
	}
	this.count++
	return ""
}

func (this *LinkedMap) Remove(key LinkedKey) interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.remove(key)
}
func (this *LinkedMap) RemoveFirst() interface{} {
	if this.IsEmpty() {
		return 0
	}
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.remove(this.header.link_next.key)
}

func (this *LinkedMap) RemoveLast() interface{} {
	if this.IsEmpty() {
		return 0
	}
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.remove(this.header.link_prev.key)
}

func (this *LinkedMap) remove(key LinkedKey) interface{} {

	tab := this.table
	index := this.hash(key) % uint(len(tab))
	e := tab[index]
	var prev *LinkedEntry = nil
	for e != nil {
		if e.key.Equals(key) {
			if prev != nil {
				prev.hash_next = e.hash_next
			} else {
				tab[index] = e.hash_next
			}
			this.count--
			oldValue := e.value
			e.value = nil
			//
			this.unchain(e)
			return oldValue
		}
		prev = e
		e = e.hash_next
	}
	return nil
}

func (this *LinkedMap) IsEmpty() bool {
	return this.count == 0
}
func (this *LinkedMap) IsFull() bool {
	return this.max > 0 && this.max <= this.count
}

func (this *LinkedMap) Clear() {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.clear()
}
func (this *LinkedMap) clear() {
	tab := this.table
	for index := len(tab) - 1; index >= 0; index-- {
		tab[index] = nil
	}
	this.header.link_next = this.header
	this.header.link_prev = this.header
	this.count = 0
}

func (this *LinkedMap) chain(link_prev *LinkedEntry, link_next *LinkedEntry, e *LinkedEntry) {
	e.link_prev = link_prev
	e.link_next = link_next
	link_prev.link_next = e
	link_next.link_prev = e
}

func (this *LinkedMap) unchain(e *LinkedEntry) {
	e.link_prev.link_next = e.link_next
	e.link_next.link_prev = e.link_prev
	e.link_prev = nil
	e.link_next = nil
}

type mapSortable struct {
	compare func(a, b LinkedKey) bool
	data    []*LinkedEntry
}

func (this mapSortable) Len() int {
	return len(this.data)
}
func (this mapSortable) Less(i, j int) bool {
	return this.compare(this.data[i].GetKey(), this.data[j].GetKey())
}

func (this mapSortable) Swap(i, j int) {
	this.data[i], this.data[j] = this.data[j], this.data[i]
}

func (this *LinkedMap) Sort(c func(k1, k2 LinkedKey) bool) {
	this.lock.Lock()
	defer this.lock.Unlock()

	sz := this.Size()
	list := make([]*LinkedEntry, sz)
	en := this.Entries()
	for i := 0; i < sz; i++ {
		list[i] = en.NextElement().(*LinkedEntry)
	}
	sort.Sort(mapSortable{compare: c, data: list})

	this.clear()
	for i := 0; i < sz; i++ {
		this.put(list[i].GetKey(), list[i].GetValue(), PUT_LAST)
	}
}

func (this *LinkedMap) ToString() string {
	this.lock.Lock()
	defer this.lock.Unlock()

	var buffer bytes.Buffer
	x := this.Entries()
	buffer.WriteString("{")
	for i := 0; x.HasMoreElements(); i++ {
		if i > 0 {
			buffer.WriteString(", ")
		}
		e := x.NextElement().(*LinkedEntry)
		buffer.WriteString(e.ToString())
	}
	buffer.WriteString("}")
	return buffer.String()
}
