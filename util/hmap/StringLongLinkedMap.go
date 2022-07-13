package hmap

import (
	"bytes"
	"sort"
	"sync"

	"github.com/whatap/golib/util/hash"
)

type StringLongLinkedMap struct {
	table      []*StringLongLinkedEntry
	header     *StringLongLinkedEntry
	count      int
	threshold  int
	loadFactor float32
	lock       sync.Mutex
	max        int
	NONE       int64
}

func NewStringLongLinkedMap() *StringLongLinkedMap {

	initCapacity := DEFAULT_CAPACITY
	loadFactor := DEFAULT_LOAD_FACTOR

	this := new(StringLongLinkedMap)
	this.loadFactor = float32(loadFactor)
	this.table = make([]*StringLongLinkedEntry, initCapacity)
	this.header = &StringLongLinkedEntry{}
	this.header.link_next = this.header
	this.header.link_prev = this.header
	this.threshold = int(float64(initCapacity) * loadFactor)

	this.NONE = 0

	return this
}

func (this *StringLongLinkedMap) SetNullValue(none int64) *StringLongLinkedMap {
	this.NONE = none
	return this
}

func (this *StringLongLinkedMap) Size() int {
	return this.count
}

func (this *StringLongLinkedMap) KeyArray() []string {
	this.lock.Lock()
	defer this.lock.Unlock()

	_keys := make([]string, this.Size())
	en := this.Keys()
	for i := 0; i < len(_keys); i++ {
		_keys[i] = en.NextString()
	}
	return _keys
}

func (this *StringLongLinkedMap) Keys() StringEnumer {
	return &StringLongLinkedEnumer{parent: this, entry: this.header.link_next}
}
func (this *StringLongLinkedMap) Values() Enumeration {
	return &StringLongLinkedEnumer{parent: this, entry: this.header.link_next, isEntry: false}
}
func (this *StringLongLinkedMap) Entries() Enumeration {
	return &StringLongLinkedEnumer{parent: this, entry: this.header.link_next, isEntry: true}
}

func (this *StringLongLinkedMap) ContainsValue(value int64) bool {
	this.lock.Lock()
	defer this.lock.Unlock()

	tab := this.table

	for i := len(tab) - 1; i > 0; i-- {
		for e := tab[i]; e != nil; e = e.hash_next {
			if e.value == value {
				return true
			}
		}
	}
	return false
}

func (this *StringLongLinkedMap) ContainsKey(key string) bool {
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
func (this *StringLongLinkedMap) Get(key string) int64 {
	this.lock.Lock()
	defer this.lock.Unlock()

	tab := this.table
	index := this.hash(key) % uint(len(tab))
	for e := tab[index]; e != nil; e = e.hash_next {
		if e.key == key {
			return e.value
		}
	}
	return this.NONE
}

func (this *StringLongLinkedMap) GetFirstKey() string {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.header.link_next.key
}

func (this *StringLongLinkedMap) GetLastKey() string {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.header.link_prev.key
}

func (this *StringLongLinkedMap) GetFirstValue() interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()

	if this.IsEmpty() {
		return this.NONE
	}

	return this.header.link_next.value
}

func (this *StringLongLinkedMap) GetLastValue() interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()

	if this.IsEmpty() {
		return this.NONE
	}

	return this.header.link_prev.value
}
func (this *StringLongLinkedMap) hash(key string) uint {
	return uint(hash.HashStr(key))
}

func (this *StringLongLinkedMap) rehash() {
	oldCapacity := len(this.table)
	oldMap := this.table
	newCapacity := oldCapacity*2 + 1
	newMap := make([]*StringLongLinkedEntry, newCapacity)
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

func (this *StringLongLinkedMap) SetMax(max int) *StringLongLinkedMap {
	this.max = max
	return this
}
func (this *StringLongLinkedMap) Put(key string, value int64) int64 {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.put(key, value, PUT_LAST)
}
func (this *StringLongLinkedMap) PutLast(key string, value int64) int64 {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.put(key, value, PUT_FORCE_LAST)
}
func (this *StringLongLinkedMap) PutFirst(key string, value int64) int64 {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.put(key, value, PUT_FORCE_FIRST)
}
func (this *StringLongLinkedMap) put(key string, value int64, m PUT_MODE) int64 {
	if key == "" {
		return this.NONE
	}

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
	e := &StringLongLinkedEntry{key: key, keyHash: keyHash, value: value, hash_next: tab[index]}
	tab[index] = e
	switch m {
	case PUT_FORCE_FIRST, PUT_FIRST:
		this.chain(this.header, this.header.link_next, e)
	case PUT_FORCE_LAST, PUT_LAST:
		this.chain(this.header.link_prev, this.header, e)
	}
	this.count++
	return this.NONE
}

func (this *StringLongLinkedMap) Add(key string, value int64) int64 {
	return this._add(key, value, PUT_LAST)
}

func (this *StringLongLinkedMap) AddLast(key string, value int64) int64 {
	return this._add(key, value, PUT_FORCE_LAST)
}

func (this *StringLongLinkedMap) AddFirst(key string, value int64) int64 {
	return this._add(key, value, PUT_FORCE_FIRST)
}
func (this *StringLongLinkedMap) _add(key string, value int64, m PUT_MODE) int64 {
	this.lock.Lock()
	defer this.lock.Unlock()

	if key == "" {
		return this.NONE
	}

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
	e := &StringLongLinkedEntry{key: key, keyHash: keyHash, value: value, hash_next: tab[index]}
	tab[index] = e
	switch m {
	case PUT_FORCE_FIRST, PUT_FIRST:
		this.chain(this.header, this.header.link_next, e)
	case PUT_FORCE_LAST, PUT_LAST:
		this.chain(this.header.link_prev, this.header, e)
	}
	this.count++

	return this.NONE
}

func (this *StringLongLinkedMap) Remove(key string) interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.remove(key)
}
func (this *StringLongLinkedMap) RemoveFirst() interface{} {
	if this.IsEmpty() {
		return this.NONE
	}
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.remove(this.header.link_next.key)
}

func (this *StringLongLinkedMap) RemoveLast() interface{} {
	if this.IsEmpty() {
		return this.NONE
	}
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.remove(this.header.link_prev.key)
}

func (this *StringLongLinkedMap) remove(key string) interface{} {

	tab := this.table
	index := this.hash(key) % uint(len(tab))
	e := tab[index]
	var prev *StringLongLinkedEntry = nil
	for e != nil {
		if e.key == key {
			if prev != nil {
				prev.hash_next = e.hash_next
			} else {
				tab[index] = e.hash_next
			}
			this.count--
			oldValue := e.value
			e.value = this.NONE
			//
			this.unchain(e)
			return oldValue
		}
		prev = e
		e = e.hash_next
	}
	return this.NONE
}

func (this *StringLongLinkedMap) IsEmpty() bool {
	return this.count == 0
}
func (this *StringLongLinkedMap) IsFull() bool {
	return this.max > 0 && this.max <= this.count
}

func (this *StringLongLinkedMap) Clear() {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.clear()
}
func (this *StringLongLinkedMap) clear() {
	tab := this.table
	for index := len(tab) - 1; index >= 0; index-- {
		tab[index] = nil
	}
	this.header.link_next = this.header
	this.header.link_prev = this.header
	this.count = 0
}

func (this *StringLongLinkedMap) chain(link_prev *StringLongLinkedEntry, link_next *StringLongLinkedEntry, e *StringLongLinkedEntry) {
	e.link_prev = link_prev
	e.link_next = link_next
	link_prev.link_next = e
	link_next.link_prev = e
}

func (this *StringLongLinkedMap) unchain(e *StringLongLinkedEntry) {
	e.link_prev.link_next = e.link_next
	e.link_next.link_prev = e.link_prev
	e.link_prev = nil
	e.link_next = nil
}

type StringLongSortable struct {
	compare func(a, b string) bool
	data    []*StringLongLinkedEntry
}

func (this StringLongSortable) Len() int {
	return len(this.data)
}
func (this StringLongSortable) Less(i, j int) bool {
	return this.compare(this.data[i].GetKey(), this.data[j].GetKey())
}

func (this StringLongSortable) Swap(i, j int) {
	this.data[i], this.data[j] = this.data[j], this.data[i]
}

func (this *StringLongLinkedMap) Sort(c func(k1, k2 string) bool) {
	this.lock.Lock()
	defer this.lock.Unlock()

	sz := this.Size()
	list := make([]*StringLongLinkedEntry, sz)
	en := this.Entries()
	for i := 0; i < sz; i++ {
		list[i] = en.NextElement().(*StringLongLinkedEntry)
	}
	sort.Sort(StringLongSortable{compare: c, data: list})

	this.clear()
	for i := 0; i < sz; i++ {
		this.put(list[i].GetKey(), list[i].GetValue(), PUT_LAST)
	}
}

func (this *StringLongLinkedMap) ToString() string {
	this.lock.Lock()
	defer this.lock.Unlock()

	var buffer bytes.Buffer
	x := this.Entries()
	buffer.WriteString("{")
	for i := 0; x.HasMoreElements(); i++ {
		if i > 0 {
			buffer.WriteString(", ")
		}
		e := x.NextElement().(*StringLongLinkedEntry)
		buffer.WriteString(e.ToString())
	}
	buffer.WriteString("}")
	return buffer.String()
}

type StringLongLinkedEnumer struct {
	parent  *StringLongLinkedMap
	entry   *StringLongLinkedEntry
	isEntry bool
	Type    int
}

func NewStringLongLinkedEnumer(parent *StringLongLinkedMap, entry *StringLongLinkedEntry, Type int) *StringLongLinkedEnumer {
	p := new(StringLongLinkedEnumer)
	p.parent = parent
	p.entry = entry
	p.Type = Type

	return p
}
func (this *StringLongLinkedEnumer) HasMoreElements() bool {
	return this.entry != nil && this.parent.header != this.entry
}
func (this *StringLongLinkedEnumer) NextElement() interface{} {
	if this.HasMoreElements() {
		e := this.entry
		this.entry = e.link_next

		switch this.Type {
		case ELEMENT_TYPE_KEYS:
			return e.key
		case ELEMENT_TYPE_VALUES:
			return e.value
		default:
			return e
		}
	}
	return nil
}

func (this *StringLongLinkedEnumer) NextLong() int64 {
	if this.HasMoreElements() {
		e := this.entry
		this.entry = e.link_next
		return e.value
	}
	return 0
}

func (this *StringLongLinkedEnumer) NextString() string {
	if this.HasMoreElements() {
		e := this.entry
		this.entry = e.link_next
		return e.key
	}
	return ""
}
