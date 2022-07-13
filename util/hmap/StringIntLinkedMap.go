package hmap

import (
	"bytes"
	"sort"
	"sync"

	"github.com/whatap/golib/util/hash"
)

type StringIntLinkedMap struct {
	table      []*StringIntLinkedEntry
	header     *StringIntLinkedEntry
	count      int
	threshold  int
	loadFactor float32
	lock       sync.Mutex
	max        int
	NONE       int32
}

func NewStringIntLinkedMap() *StringIntLinkedMap {

	initCapacity := DEFAULT_CAPACITY
	loadFactor := DEFAULT_LOAD_FACTOR

	this := new(StringIntLinkedMap)
	this.loadFactor = float32(loadFactor)
	this.table = make([]*StringIntLinkedEntry, initCapacity)
	this.header = &StringIntLinkedEntry{}
	this.header.link_next = this.header
	this.header.link_prev = this.header
	this.threshold = int(float64(initCapacity) * loadFactor)

	this.NONE = 0

	return this
}

func (this *StringIntLinkedMap) SetNullValue(none int32) *StringIntLinkedMap {
	this.NONE = none
	return this
}

func (this *StringIntLinkedMap) Size() int {
	return this.count
}

func (this *StringIntLinkedMap) KeyArray() []string {
	this.lock.Lock()
	defer this.lock.Unlock()

	_keys := make([]string, this.Size())
	en := this.Keys()
	for i := 0; i < len(_keys); i++ {
		_keys[i] = en.NextString()
	}
	return _keys
}

func (this *StringIntLinkedMap) Keys() StringEnumer {
	return &StringIntLinkedEnumer{parent: this, entry: this.header.link_next}
}
func (this *StringIntLinkedMap) Values() Enumeration {
	return &StringIntLinkedEnumer{parent: this, entry: this.header.link_next, isEntry: false}
}
func (this *StringIntLinkedMap) Entries() Enumeration {
	return &StringIntLinkedEnumer{parent: this, entry: this.header.link_next, isEntry: true}
}

func (this *StringIntLinkedMap) ContainsValue(value int32) bool {
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

func (this *StringIntLinkedMap) ContainsKey(key string) bool {
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
func (this *StringIntLinkedMap) Get(key string) int32 {
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

func (this *StringIntLinkedMap) GetFirstKey() string {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.header.link_next.key
}

func (this *StringIntLinkedMap) GetLastKey() string {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.header.link_prev.key
}

func (this *StringIntLinkedMap) GetFirstValue() interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()

	if this.IsEmpty() {
		return this.NONE
	}

	return this.header.link_next.value
}

func (this *StringIntLinkedMap) GetLastValue() interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()

	if this.IsEmpty() {
		return this.NONE
	}

	return this.header.link_prev.value
}
func (this *StringIntLinkedMap) hash(key string) uint {
	return uint(hash.HashStr(key))
}

func (this *StringIntLinkedMap) rehash() {
	oldCapacity := len(this.table)
	oldMap := this.table
	newCapacity := oldCapacity*2 + 1
	newMap := make([]*StringIntLinkedEntry, newCapacity)
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

func (this *StringIntLinkedMap) SetMax(max int) *StringIntLinkedMap {
	this.max = max
	return this
}
func (this *StringIntLinkedMap) Put(key string, value int32) int32 {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.put(key, value, PUT_LAST)
}
func (this *StringIntLinkedMap) PutLast(key string, value int32) int32 {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.put(key, value, PUT_FORCE_LAST)
}
func (this *StringIntLinkedMap) PutFirst(key string, value int32) int32 {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.put(key, value, PUT_FORCE_FIRST)
}
func (this *StringIntLinkedMap) put(key string, value int32, m PUT_MODE) int32 {
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
	e := &StringIntLinkedEntry{key: key, keyHash: keyHash, value: value, hash_next: tab[index]}
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

func (this *StringIntLinkedMap) Add(key string, value int32) int32 {
	return this._add(key, value, PUT_LAST)
}

func (this *StringIntLinkedMap) AddLast(key string, value int32) int32 {
	return this._add(key, value, PUT_FORCE_LAST)
}

func (this *StringIntLinkedMap) AddFirst(key string, value int32) int32 {
	return this._add(key, value, PUT_FORCE_FIRST)
}
func (this *StringIntLinkedMap) _add(key string, value int32, m PUT_MODE) int32 {
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
	e := &StringIntLinkedEntry{key: key, keyHash: keyHash, value: value, hash_next: tab[index]}
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

func (this *StringIntLinkedMap) Remove(key string) interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.remove(key)
}
func (this *StringIntLinkedMap) RemoveFirst() interface{} {
	if this.IsEmpty() {
		return this.NONE
	}
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.remove(this.header.link_next.key)
}

func (this *StringIntLinkedMap) RemoveLast() interface{} {
	if this.IsEmpty() {
		return this.NONE
	}
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.remove(this.header.link_prev.key)
}

func (this *StringIntLinkedMap) remove(key string) interface{} {

	tab := this.table
	index := this.hash(key) % uint(len(tab))
	e := tab[index]
	var prev *StringIntLinkedEntry = nil
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

func (this *StringIntLinkedMap) IsEmpty() bool {
	return this.count == 0
}
func (this *StringIntLinkedMap) IsFull() bool {
	return this.max > 0 && this.max <= this.count
}

func (this *StringIntLinkedMap) Clear() {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.clear()
}
func (this *StringIntLinkedMap) clear() {
	tab := this.table
	for index := len(tab) - 1; index >= 0; index-- {
		tab[index] = nil
	}
	this.header.link_next = this.header
	this.header.link_prev = this.header
	this.count = 0
}

func (this *StringIntLinkedMap) chain(link_prev *StringIntLinkedEntry, link_next *StringIntLinkedEntry, e *StringIntLinkedEntry) {
	e.link_prev = link_prev
	e.link_next = link_next
	link_prev.link_next = e
	link_next.link_prev = e
}

func (this *StringIntLinkedMap) unchain(e *StringIntLinkedEntry) {
	e.link_prev.link_next = e.link_next
	e.link_next.link_prev = e.link_prev
	e.link_prev = nil
	e.link_next = nil
}

type StringIntSortable struct {
	compare func(a, b string) bool
	data    []*StringIntLinkedEntry
}

func (this StringIntSortable) Len() int {
	return len(this.data)
}
func (this StringIntSortable) Less(i, j int) bool {
	return this.compare(this.data[i].GetKey(), this.data[j].GetKey())
}

func (this StringIntSortable) Swap(i, j int) {
	this.data[i], this.data[j] = this.data[j], this.data[i]
}

func (this *StringIntLinkedMap) Sort(c func(k1, k2 string) bool) {
	this.lock.Lock()
	defer this.lock.Unlock()

	sz := this.Size()
	list := make([]*StringIntLinkedEntry, sz)
	en := this.Entries()
	for i := 0; i < sz; i++ {
		list[i] = en.NextElement().(*StringIntLinkedEntry)
	}
	sort.Sort(StringIntSortable{compare: c, data: list})

	this.clear()
	for i := 0; i < sz; i++ {
		this.put(list[i].GetKey(), list[i].GetValue(), PUT_LAST)
	}
}

func (this *StringIntLinkedMap) ToString() string {
	this.lock.Lock()
	defer this.lock.Unlock()

	var buffer bytes.Buffer
	x := this.Entries()
	buffer.WriteString("{")
	for i := 0; x.HasMoreElements(); i++ {
		if i > 0 {
			buffer.WriteString(", ")
		}
		e := x.NextElement().(*StringIntLinkedEntry)
		buffer.WriteString(e.ToString())
	}
	buffer.WriteString("}")
	return buffer.String()
}

type StringIntLinkedEnumer struct {
	parent  *StringIntLinkedMap
	entry   *StringIntLinkedEntry
	isEntry bool
	Type    int
}

func NewStringIntLinkedEnumer(parent *StringIntLinkedMap, entry *StringIntLinkedEntry, Type int) *StringIntLinkedEnumer {
	p := new(StringIntLinkedEnumer)
	p.parent = parent
	p.entry = entry
	p.Type = Type

	return p
}
func (this *StringIntLinkedEnumer) HasMoreElements() bool {
	return this.entry != nil && this.parent.header != this.entry
}
func (this *StringIntLinkedEnumer) NextElement() interface{} {
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

func (this *StringIntLinkedEnumer) NextInt() int32 {
	if this.HasMoreElements() {
		e := this.entry
		this.entry = e.link_next
		return e.value
	}
	return 0
}

func (this *StringIntLinkedEnumer) NextString() string {
	if this.HasMoreElements() {
		e := this.entry
		this.entry = e.link_next
		return e.key
	}
	return ""
}
