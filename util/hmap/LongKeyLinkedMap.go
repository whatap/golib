package hmap

import (
	"bytes"
	"sort"
	"sync"
)

type LongKeyLinkedMap struct {
	table      []*LongKeyLinkedEntry
	header     *LongKeyLinkedEntry
	count      int
	threshold  int
	loadFactor float32
	lock       sync.Mutex
	max        int
}

func NewLongKeyLinkedMapDefault() *LongKeyLinkedMap {

	initCapacity := DEFAULT_CAPACITY
	loadFactor := DEFAULT_LOAD_FACTOR

	this := new(LongKeyLinkedMap)
	this.loadFactor = float32(loadFactor)
	this.table = make([]*LongKeyLinkedEntry, initCapacity)
	this.header = &LongKeyLinkedEntry{}
	this.header.link_next = this.header
	this.header.link_prev = this.header
	this.threshold = int(float64(initCapacity) * loadFactor)
	return this
}

func NewLongKeyLinkedMap(initCapacity int, loadFactor float32) *LongKeyLinkedMap {
	this := new(LongKeyLinkedMap)
	this.loadFactor = float32(loadFactor)
	this.table = make([]*LongKeyLinkedEntry, initCapacity)
	this.header = &LongKeyLinkedEntry{}
	this.header.link_next = this.header
	this.header.link_prev = this.header
	this.threshold = int(float32(initCapacity) * loadFactor)
	return this
}

func (this *LongKeyLinkedMap) Size() int {
	return this.count
}

func (this *LongKeyLinkedMap) KeyArray() []int64 {
	this.lock.Lock()
	defer this.lock.Unlock()

	_keys := make([]int64, this.Size())
	en := this.Keys()
	for i := 0; i < len(_keys); i++ {
		_keys[i] = en.NextLong()
	}
	return _keys
}

type LongKeyEnumerImpl struct {
	parent  *LongKeyLinkedMap
	entry   *LongKeyLinkedEntry
	isEntry bool
}

func (this *LongKeyEnumerImpl) HasMoreElements() bool {
	return this.entry != nil && this.parent.header != this.entry
}
func (this *LongKeyEnumerImpl) NextLong() int64 {
	if this.HasMoreElements() {
		e := this.entry
		this.entry = e.link_next
		return e.key
	}
	return 0
}

func (this *LongKeyEnumerImpl) NextElement() interface{} {
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
func (this *LongKeyLinkedMap) Keys() LongEnumer {
	return &LongKeyEnumerImpl{parent: this, entry: this.header.link_next}
}
func (this *LongKeyLinkedMap) Values() Enumeration {
	return &LongKeyEnumerImpl{parent: this, entry: this.header.link_next, isEntry: false}
}
func (this *LongKeyLinkedMap) Entries() Enumeration {
	return &LongKeyEnumerImpl{parent: this, entry: this.header.link_next, isEntry: true}
}

func (this *LongKeyLinkedMap) ContainsKey(key int64) bool {
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
func (this *LongKeyLinkedMap) Get(key int64) interface{} {
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
func (this *LongKeyLinkedMap) GetFirstKey() int64 {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.header.link_next.key
}

func (this *LongKeyLinkedMap) GetLastKey() int64 {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.header.link_prev.key
}

func (this *LongKeyLinkedMap) GetFirstValue() interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.header.link_next.value
}

func (this *LongKeyLinkedMap) GetLastValue() interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.header.link_prev.value
}
func (this *LongKeyLinkedMap) hash(key int64) uint {
	return uint(key ^ key>>32)
}

func (this *LongKeyLinkedMap) rehash() {
	oldCapacity := len(this.table)
	oldMap := this.table
	newCapacity := oldCapacity*2 + 1
	newMap := make([]*LongKeyLinkedEntry, newCapacity)
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

func (this *LongKeyLinkedMap) SetMax(max int) *LongKeyLinkedMap {
	this.max = max
	return this
}
func (this *LongKeyLinkedMap) Put(key int64, value interface{}) interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.put(key, value, PUT_LAST)
}
func (this *LongKeyLinkedMap) PutLast(key int64, value interface{}) interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.put(key, value, PUT_FORCE_LAST)
}
func (this *LongKeyLinkedMap) PutFirst(key int64, value interface{}) interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.put(key, value, PUT_FORCE_FIRST)
}
func (this *LongKeyLinkedMap) put(key int64, value interface{}, m PUT_MODE) interface{} {
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
			//log.Println(">>>>", "Put Dup Txid=", key)
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
		}
	}
	if this.count >= this.threshold {
		this.rehash()
		tab = this.table
		index = keyHash % uint(len(tab))
	}
	e := &LongKeyLinkedEntry{key: key, keyHash: keyHash, value: value, hash_next: tab[index]}
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

func (this *LongKeyLinkedMap) Remove(key int64) interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.remove(key)
}
func (this *LongKeyLinkedMap) RemoveFirst() interface{} {
	if this.IsEmpty() {
		return 0
	}
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.remove(this.header.link_next.key)
}

func (this *LongKeyLinkedMap) RemoveLast() interface{} {
	if this.IsEmpty() {
		return 0
	}
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.remove(this.header.link_prev.key)
}

func (this *LongKeyLinkedMap) remove(key int64) interface{} {
	tab := this.table
	index := this.hash(key) % uint(len(tab))
	e := tab[index]
	var prev *LongKeyLinkedEntry = nil
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

func (this *LongKeyLinkedMap) IsEmpty() bool {
	return this.count == 0
}
func (this *LongKeyLinkedMap) IsFull() bool {
	return this.max > 0 && this.max <= this.count
}

func (this *LongKeyLinkedMap) Clear() {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.clear()
}
func (this *LongKeyLinkedMap) clear() {
	tab := this.table
	for index := len(tab) - 1; index >= 0; index-- {
		tab[index] = nil
	}
	this.header.link_next = this.header
	this.header.link_prev = this.header
	this.count = 0
}

func (this *LongKeyLinkedMap) chain(link_prev *LongKeyLinkedEntry, link_next *LongKeyLinkedEntry, e *LongKeyLinkedEntry) {
	e.link_prev = link_prev
	e.link_next = link_next
	link_prev.link_next = e
	link_next.link_prev = e
}

func (this *LongKeyLinkedMap) unchain(e *LongKeyLinkedEntry) {
	e.link_prev.link_next = e.link_next
	e.link_next.link_prev = e.link_prev
	e.link_prev = nil
	e.link_next = nil
}

type LongKeySortable struct {
	compare func(a, b int64) bool
	data    []*LongKeyLinkedEntry
}

func (this LongKeySortable) Len() int {
	return len(this.data)
}
func (this LongKeySortable) Less(i, j int) bool {
	return this.compare(this.data[i].GetKey(), this.data[j].GetKey())
}

func (this LongKeySortable) Swap(i, j int) {
	this.data[i], this.data[j] = this.data[j], this.data[i]
}

func (this *LongKeyLinkedMap) Sort(c func(k1, k2 int64) bool) {
	this.lock.Lock()
	defer this.lock.Unlock()

	sz := this.Size()
	list := make([]*LongKeyLinkedEntry, sz)
	en := this.Entries()
	for i := 0; i < sz; i++ {
		list[i] = en.NextElement().(*LongKeyLinkedEntry)
	}
	sort.Sort(LongKeySortable{compare: c, data: list})

	this.clear()
	for i := 0; i < sz; i++ {
		this.put(list[i].GetKey(), list[i].GetValue(), PUT_LAST)
	}
}

func (this *LongKeyLinkedMap) ToString() string {
	this.lock.Lock()
	defer this.lock.Unlock()

	var buffer bytes.Buffer
	x := this.Entries()
	buffer.WriteString("{")
	for i := 0; x.HasMoreElements(); i++ {
		if i > 0 {
			buffer.WriteString(", ")
		}
		e := x.NextElement().(*LongKeyLinkedEntry)
		buffer.WriteString(e.ToString())
	}
	buffer.WriteString("}")
	return buffer.String()
}
