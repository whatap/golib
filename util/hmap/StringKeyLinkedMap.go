package hmap

import (
	"bytes"
	"sort"
	"sync"

	"github.com/whatap/golib/util/hash"
)

type StringKeyLinkedMap struct {
	table      []*StringKeyLinkedEntry
	header     *StringKeyLinkedEntry
	count      int
	threshold  int
	loadFactor float32
	lock       sync.Mutex
	max        int
}

func NewStringKeyLinkedMap() *StringKeyLinkedMap {

	initCapacity := DEFAULT_CAPACITY
	loadFactor := DEFAULT_LOAD_FACTOR

	this := new(StringKeyLinkedMap)
	this.loadFactor = float32(loadFactor)
	this.table = make([]*StringKeyLinkedEntry, initCapacity)
	this.header = &StringKeyLinkedEntry{}
	this.header.link_next = this.header
	this.header.link_prev = this.header
	this.threshold = int(float64(initCapacity) * loadFactor)
	return this
}

func (this *StringKeyLinkedMap) Size() int {
	return this.count
}

func (this *StringKeyLinkedMap) KeyArray() []string {
	this.lock.Lock()
	defer this.lock.Unlock()

	_keys := make([]string, this.Size())
	en := this.Keys()
	for i := 0; i < len(_keys); i++ {
		_keys[i] = en.NextString()
	}
	return _keys
}

type StringKeyEnumerImpl struct {
	parent  *StringKeyLinkedMap
	entry   *StringKeyLinkedEntry
	isEntry bool
}

func (this *StringKeyEnumerImpl) HasMoreElements() bool {
	return this.entry != nil && this.parent.header != this.entry
}
func (this *StringKeyEnumerImpl) NextString() string {
	if this.HasMoreElements() {
		e := this.entry
		this.entry = e.link_next
		return e.key
	}
	return ""
}

func (this *StringKeyEnumerImpl) NextElement() interface{} {
	if this.HasMoreElements() {
		e := this.entry
		this.entry = e.link_next
		if this.isEntry {
			return e
		} else {
			return e.GetValue()
		}
	}
	return ""
}
func (this *StringKeyLinkedMap) Keys() StringEnumer {
	return &StringKeyEnumerImpl{parent: this, entry: this.header.link_next}
}
func (this *StringKeyLinkedMap) Values() Enumeration {
	return &StringKeyEnumerImpl{parent: this, entry: this.header.link_next, isEntry: false}
}
func (this *StringKeyLinkedMap) Entries() Enumeration {
	return &StringKeyEnumerImpl{parent: this, entry: this.header.link_next, isEntry: true}
}

func (this *StringKeyLinkedMap) ContainsKey(key string) bool {
	this.lock.Lock()
	defer this.lock.Unlock()

	tab := this.table
	index := this.hash(key) % uint(len(tab))
	for e := tab[index]; e != nil; e = e.hash_next {
		if e.key == key {
			return true
		}
	}
	return false

}
func (this *StringKeyLinkedMap) Get(key string) interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()

	tab := this.table
	index := this.hash(key) % uint(len(tab))
	for e := tab[index]; e != nil; e = e.hash_next {
		if e.key == key {
			return e.value
		}
	}
	return nil
}
func (this *StringKeyLinkedMap) GetFirstKey() string {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.header.link_next.key
}

func (this *StringKeyLinkedMap) GetLastKey() string {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.header.link_prev.key
}

func (this *StringKeyLinkedMap) GetFirstValue() interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.header.link_next.value
}

func (this *StringKeyLinkedMap) GetLastValue() interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.header.link_prev.value
}
func (this *StringKeyLinkedMap) hash(key string) uint {
	return uint(hash.HashStr(key))
}

func (this *StringKeyLinkedMap) rehash() {
	oldCapacity := len(this.table)
	oldMap := this.table
	newCapacity := oldCapacity*2 + 1
	newMap := make([]*StringKeyLinkedEntry, newCapacity)
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

func (this *StringKeyLinkedMap) SetMax(max int) *StringKeyLinkedMap {
	this.max = max
	return this
}
func (this *StringKeyLinkedMap) Put(key string, value interface{}) interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.put(key, value, PUT_LAST)
}
func (this *StringKeyLinkedMap) PutLast(key string, value interface{}) interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.put(key, value, PUT_FORCE_LAST)
}
func (this *StringKeyLinkedMap) PutFirst(key string, value interface{}) interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.put(key, value, PUT_FORCE_FIRST)
}
func (this *StringKeyLinkedMap) put(key string, value interface{}, m PUT_MODE) interface{} {

	tab := this.table
	keyHash := this.hash(key)
	index := keyHash % uint(len(tab))
	for e := tab[index]; e != nil; e = e.hash_next {
		if e.key == key {
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
	e := &StringKeyLinkedEntry{key: key, keyHash: keyHash, value: value, hash_next: tab[index]}
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

func (this *StringKeyLinkedMap) Remove(key string) interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.remove(key)
}
func (this *StringKeyLinkedMap) RemoveFirst() interface{} {
	if this.IsEmpty() {
		return 0
	}
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.remove(this.header.link_next.key)
}

func (this *StringKeyLinkedMap) RemoveLast() interface{} {
	if this.IsEmpty() {
		return 0
	}
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.remove(this.header.link_prev.key)
}

func (this *StringKeyLinkedMap) remove(key string) interface{} {

	tab := this.table
	index := this.hash(key) % uint(len(tab))
	e := tab[index]
	var prev *StringKeyLinkedEntry = nil
	for e != nil {
		if e.key == key {
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

func (this *StringKeyLinkedMap) IsEmpty() bool {
	return this.count == 0
}
func (this *StringKeyLinkedMap) IsFull() bool {
	return this.max > 0 && this.max <= this.count
}

func (this *StringKeyLinkedMap) Clear() {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.clear()
}
func (this *StringKeyLinkedMap) clear() {
	tab := this.table
	for index := len(tab) - 1; index >= 0; index-- {
		tab[index] = nil
	}
	this.header.link_next = this.header
	this.header.link_prev = this.header
	this.count = 0
}

func (this *StringKeyLinkedMap) chain(link_prev *StringKeyLinkedEntry, link_next *StringKeyLinkedEntry, e *StringKeyLinkedEntry) {
	e.link_prev = link_prev
	e.link_next = link_next
	link_prev.link_next = e
	link_next.link_prev = e
}

func (this *StringKeyLinkedMap) unchain(e *StringKeyLinkedEntry) {
	e.link_prev.link_next = e.link_next
	e.link_next.link_prev = e.link_prev
	e.link_prev = nil
	e.link_next = nil
}

type StringKeySortable struct {
	compare func(a, b string) bool
	data    []*StringKeyLinkedEntry
}

func (this StringKeySortable) Len() int {
	return len(this.data)
}
func (this StringKeySortable) Less(i, j int) bool {
	return this.compare(this.data[i].GetKey(), this.data[j].GetKey())
}

func (this StringKeySortable) Swap(i, j int) {
	this.data[i], this.data[j] = this.data[j], this.data[i]
}

func (this *StringKeyLinkedMap) Sort(c func(k1, k2 string) bool) {
	this.lock.Lock()
	defer this.lock.Unlock()

	sz := this.Size()
	list := make([]*StringKeyLinkedEntry, sz)
	en := this.Entries()
	for i := 0; i < sz; i++ {
		list[i] = en.NextElement().(*StringKeyLinkedEntry)
	}
	sort.Sort(StringKeySortable{compare: c, data: list})

	this.clear()
	for i := 0; i < sz; i++ {
		this.put(list[i].GetKey(), list[i].GetValue(), PUT_LAST)
	}
}

func (this *StringKeyLinkedMap) ToString() string {
	this.lock.Lock()
	defer this.lock.Unlock()

	var buffer bytes.Buffer
	x := this.Entries()
	buffer.WriteString("{")
	for i := 0; x.HasMoreElements(); i++ {
		if i > 0 {
			buffer.WriteString(", ")
		}
		e := x.NextElement().(*StringKeyLinkedEntry)
		buffer.WriteString(e.ToString())
	}
	buffer.WriteString("}")
	return buffer.String()
}
