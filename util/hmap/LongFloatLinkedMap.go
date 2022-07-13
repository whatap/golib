package hmap

import (
	"bytes"
	"sort"
	"sync"

	"github.com/whatap/golib/io"
)

type LongFloatLinkedMap struct {
	table      []*LongFloatLinkedEntry
	header     *LongFloatLinkedEntry
	count      int
	threshold  int
	loadFactor float32
	NONE       float32
	lock       sync.Mutex
	max        int
}

func NewLongFloatLinkedMap() *LongFloatLinkedMap {

	initCapacity := DEFAULT_CAPACITY
	loadFactor := DEFAULT_LOAD_FACTOR

	this := new(LongFloatLinkedMap)
	this.loadFactor = float32(loadFactor)
	this.table = make([]*LongFloatLinkedEntry, initCapacity)
	this.header = &LongFloatLinkedEntry{}
	this.header.link_next = this.header
	this.header.link_prev = this.header
	this.threshold = int(float64(initCapacity) * loadFactor)
	return this
}

func (this *LongFloatLinkedMap) Size() int {
	return this.count
}

func (this *LongFloatLinkedMap) KeyArray() []int64 {
	this.lock.Lock()
	defer this.lock.Unlock()

	_keys := make([]int64, this.Size())
	en := this.Keys()
	for i := 0; i < len(_keys); i++ {
		_keys[i] = en.NextLong()
	}
	return _keys
}

type LongFloatEnumerImpl struct {
	parent *LongFloatLinkedMap
	entry  *LongFloatLinkedEntry
}

func (this *LongFloatEnumerImpl) HasMoreElements() bool {
	return this.entry != nil && this.parent.header != this.entry
}
func (this *LongFloatEnumerImpl) NextLong() int64 {
	if this.HasMoreElements() {
		e := this.entry
		this.entry = e.link_next
		return e.key
	}
	return 0
}
func (this *LongFloatEnumerImpl) NextFloat() float32 {
	if this.HasMoreElements() {
		e := this.entry
		this.entry = e.link_next
		return e.value
	}
	return this.parent.NONE
}
func (this *LongFloatEnumerImpl) NextElement() interface{} {
	if this.HasMoreElements() {
		e := this.entry
		this.entry = e.link_next
		return e
	}
	return nil
}
func (this *LongFloatLinkedMap) Keys() LongEnumer {
	return &LongFloatEnumerImpl{parent: this, entry: this.header.link_next}
}
func (this *LongFloatLinkedMap) Values() FloatEnumer {
	return &LongFloatEnumerImpl{parent: this, entry: this.header.link_next}
}
func (this *LongFloatLinkedMap) Entries() Enumeration {
	return &LongFloatEnumerImpl{parent: this, entry: this.header.link_next}
}

func (this *LongFloatLinkedMap) ContainsValue(value float32) bool {
	this.lock.Lock()
	defer this.lock.Unlock()

	tab := this.table
	for i := len(tab); i > 0; i-- {
		for e := tab[i]; e != nil; e = e.hash_next {
			if e.value == value {
				return true
			}
		}
	}
	return false
}
func (this *LongFloatLinkedMap) ContainsKey(key int64) bool {
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
func (this *LongFloatLinkedMap) Get(key int64) float32 {
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
func (this *LongFloatLinkedMap) GetFirstKey() int64 {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.header.link_next.key
}

func (this *LongFloatLinkedMap) GetLastKey() int64 {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.header.link_prev.key
}

func (this *LongFloatLinkedMap) GetFirstValue() float32 {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.header.link_next.value
}

func (this *LongFloatLinkedMap) GetLastValue() float32 {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.header.link_prev.value
}
func (this *LongFloatLinkedMap) hash(key int64) uint {
	return uint(key)
}

func (this *LongFloatLinkedMap) rehash() {
	oldCapacity := len(this.table)
	oldMap := this.table
	newCapacity := oldCapacity*2 + 1
	newMap := make([]*LongFloatLinkedEntry, newCapacity)
	this.threshold = int(float32(newCapacity) * this.loadFactor)
	this.table = newMap
	for i := oldCapacity; i > 0; i-- {
		for old := oldMap[i-1]; old != nil; {
			e := old
			old = old.hash_next
			key := e.key
			index := int(this.hash(key) % uint(newCapacity))
			e.hash_next = newMap[index]
			newMap[index] = e
		}
	}
}

func (this *LongFloatLinkedMap) SetMax(max int) *LongFloatLinkedMap {
	this.max = max
	return this
}
func (this *LongFloatLinkedMap) Put(key int64, value float32) float32 {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.put(key, value, PUT_LAST)
}
func (this *LongFloatLinkedMap) PutLast(key int64, value float32) float32 {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.put(key, value, PUT_FORCE_LAST)
}
func (this *LongFloatLinkedMap) PutFirst(key int64, value float32) float32 {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.put(key, value, PUT_FORCE_FIRST)
}
func (this *LongFloatLinkedMap) put(key int64, value float32, m PUT_MODE) float32 {

	tab := this.table
	index := this.hash(key) % uint(len(tab))
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
		index = this.hash(key) % uint(len(tab))
	}
	e := &LongFloatLinkedEntry{key: key, value: value, hash_next: tab[index]}
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

func (this *LongFloatLinkedMap) Add(key int64, value float32) float32 {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.add(key, value, PUT_LAST)
}
func (this *LongFloatLinkedMap) AddLast(key int64, value float32) float32 {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.add(key, value, PUT_FORCE_LAST)
}
func (this *LongFloatLinkedMap) AddFirst(key int64, value float32) float32 {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.add(key, value, PUT_FORCE_FIRST)
}

func (this *LongFloatLinkedMap) add(key int64, value float32, m PUT_MODE) float32 {
	tab := this.table
	index := this.hash(key) % uint(len(tab))
	for e := tab[index]; e != nil; e = e.hash_next {
		if e.key == key {
			old := e.value
			e.value += value
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
				k := this.header.link_next.key
				this.remove(k)
			}
			break
		}
	}
	if this.count >= this.threshold {
		this.rehash()
		tab = this.table
		index = this.hash(key) % uint(len(tab))
	}
	e := &LongFloatLinkedEntry{key: key, value: value, hash_next: tab[index]}
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
func (this *LongFloatLinkedMap) Remove(key int64) float32 {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.remove(key)
}
func (this *LongFloatLinkedMap) RemoveFirst() float32 {
	if this.IsEmpty() {
		return 0
	}
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.remove(this.header.link_next.key)
}

func (this *LongFloatLinkedMap) RemoveLast() float32 {
	if this.IsEmpty() {
		return 0
	}
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.remove(this.header.link_prev.key)
}

func (this *LongFloatLinkedMap) remove(key int64) float32 {

	tab := this.table
	index := this.hash(key) % uint(len(tab))
	e := tab[index]
	var prev *LongFloatLinkedEntry = nil
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

func (this *LongFloatLinkedMap) IsEmpty() bool {
	return this.count == 0
}
func (this *LongFloatLinkedMap) IsFull() bool {
	return this.max > 0 && this.max <= this.count
}

func (this *LongFloatLinkedMap) Clear() {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.clear()
}
func (this *LongFloatLinkedMap) clear() {
	tab := this.table
	for index := len(tab) - 1; index >= 0; index-- {
		tab[index] = nil
	}
	this.header.link_next = this.header
	this.header.link_prev = this.header
	this.count = 0
}

func (this *LongFloatLinkedMap) chain(link_prev *LongFloatLinkedEntry, link_next *LongFloatLinkedEntry, e *LongFloatLinkedEntry) {
	e.link_prev = link_prev
	e.link_next = link_next
	link_prev.link_next = e
	link_next.link_prev = e
}

func (this *LongFloatLinkedMap) unchain(e *LongFloatLinkedEntry) {
	e.link_prev.link_next = e.link_next
	e.link_next.link_prev = e.link_prev
	e.link_prev = nil
	e.link_next = nil
}

type longFloatSortable struct {
	compare func(a, b int64) bool
	data    []*LongFloatLinkedEntry
}

func (this longFloatSortable) Len() int {
	return len(this.data)
}
func (this longFloatSortable) Less(i, j int) bool {
	return this.compare(this.data[i].GetKey(), this.data[j].GetKey())
}

func (this longFloatSortable) Swap(i, j int) {
	this.data[i], this.data[j] = this.data[j], this.data[i]
}

func (this *LongFloatLinkedMap) Sort(c func(k1, k2 int64) bool) {
	this.lock.Lock()
	defer this.lock.Unlock()

	sz := this.Size()
	list := make([]*LongFloatLinkedEntry, sz)
	en := this.Entries()
	for i := 0; i < sz; i++ {
		list[i] = en.NextElement().(*LongFloatLinkedEntry)
	}
	sort.Sort(longFloatSortable{compare: c, data: list})

	this.clear()
	for i := 0; i < sz; i++ {
		this.put(list[i].GetKey(), list[i].GetValue(), PUT_LAST)
	}
}

func (this *LongFloatLinkedMap) ToBytes(dout *io.DataOutputX) {
	dout.WriteDecimal(int64(this.Size()))
	en := this.Entries()
	for en.HasMoreElements() {
		e := en.NextElement().(LongFloatLinkedEntry)
		dout.WriteDecimal(int64(e.GetKey()))
		dout.WriteFloat(e.GetValue())
	}
}

func (this *LongFloatLinkedMap) ToObject(din *io.DataInputX) *LongFloatLinkedMap {
	cnt := int(din.ReadDecimal())
	for i := 0; i < cnt; i++ {
		key := int64(din.ReadDecimal())
		value := din.ReadFloat()
		this.Put(key, value)
	}
	return this
}
func (this *LongFloatLinkedMap) ToString() string {
	this.lock.Lock()
	defer this.lock.Unlock()

	var buffer bytes.Buffer
	x := this.Entries()
	buffer.WriteString("{")
	for i := 0; x.HasMoreElements(); i++ {
		if i > 0 {
			buffer.WriteString(", ")
		}
		e := x.NextElement().(*LongFloatLinkedEntry)
		buffer.WriteString(e.ToString())
	}
	buffer.WriteString("}")
	return buffer.String()
}
