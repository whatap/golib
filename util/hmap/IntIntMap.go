package hmap

import (
	"bytes"
	"sort"
	"sync"

	"github.com/whatap/golib/io"
)

type IntIntMap struct {
	table      []*IntIntEntry
	count      int
	threshold  int
	loadFactor float32
	NONE       int32
	lock       sync.Mutex
	max        int
}

func NewIntIntMap(initCapacity int, loadFactor float32) *IntIntMap {

	this := new(IntIntMap)
	this.loadFactor = float32(loadFactor)
	this.table = make([]*IntIntEntry, initCapacity)
	this.threshold = int(float32(initCapacity) * loadFactor)

	//fmt.Printf("this=%p, threshold=%d \r\n", this, this.threshold)

	return this
}

func NewIntIntMapDefault() *IntIntMap {

	return NewIntIntMap(DEFAULT_CAPACITY, DEFAULT_LOAD_FACTOR)
}

func (this *IntIntMap) Size() int {
	return this.count
}

func (this *IntIntMap) Keys() IntEnumer {
	//return &IntIntMapEnumer{isKey: true, parent: this, entry: this.header.link_next}
	en := NewIntIntMapEnumer(ELEMENT_TYPE_KEYS)
	en.table = this.table
	en.index = len(en.table)

	return en
}
func (this *IntIntMap) Values() IntEnumer {
	//return &IntIntMapEnumer{isKey: false, parent: this, entry: this.header.link_next}
	en := NewIntIntMapEnumer(ELEMENT_TYPE_VALUES)
	en.table = this.table
	en.index = len(en.table)

	return en
}
func (this *IntIntMap) Entries() Enumeration {
	//return &IntIntMapEnumer{parent: this, entry: this.header.link_next}
	en := NewIntIntMapEnumer(ELEMENT_TYPE_ENTRIES)
	en.table = this.table
	en.index = len(en.table)

	return en
}

func (this *IntIntMap) ContainsValue(value int32) bool {
	this.lock.Lock()
	defer this.lock.Unlock()

	tab := this.table
	for i := len(tab); i > 0; i-- {
		for e := tab[i]; e != nil; e = e.next {
			if e.value == value {
				return true
			}
		}
	}
	return false
}
func (this *IntIntMap) ContainsKey(key int32) bool {
	this.lock.Lock()
	defer this.lock.Unlock()

	tab := this.table
	index := this.hash(key) % uint(len(tab))
	for e := tab[index]; e != nil; e = e.next {
		if e.key == key {
			return true
		}
	}
	return false

}
func (this *IntIntMap) hash(key int32) uint {
	return uint(key)
}

func (this *IntIntMap) Get(key int32) int32 {
	this.lock.Lock()
	defer this.lock.Unlock()

	tab := this.table
	index := this.hash(key) % uint(len(tab))
	for e := tab[index]; e != nil; e = e.next {
		if e.key == key {
			return e.value
		}
	}
	return this.NONE
}

func (this *IntIntMap) rehash() {
	oldCapacity := len(this.table)
	oldMap := this.table
	newCapacity := oldCapacity*2 + 1
	newMap := make([]*IntIntEntry, newCapacity)

	this.threshold = int(float32(newCapacity) * this.loadFactor)
	this.table = newMap
	for i := oldCapacity; i > 0; i-- {
		for old := oldMap[i-1]; old != nil; {
			e := old
			old = old.next
			key := e.key
			index := int(this.hash(key) % uint(newCapacity))
			e.next = newMap[index]
			newMap[index] = e
		}
	}
}

func (this *IntIntMap) KeyArray() []int32 {
	this.lock.Lock()
	defer this.lock.Unlock()

	_keys := make([]int32, this.Size())
	en := this.Keys()
	for i := 0; i < len(_keys); i++ {
		_keys[i] = en.NextInt()
	}
	return _keys
}

func (this *IntIntMap) ValueArray() []int32 {
	this.lock.Lock()
	defer this.lock.Unlock()

	_values := make([]int32, this.Size())
	en := this.Values()
	for i := 0; i < len(_values); i++ {
		_values[i] = en.NextInt()
	}
	return _values
}

func (this *IntIntMap) SetMax(max int) *IntIntMap {
	this.max = max
	return this
}
func (this *IntIntMap) Put(key int32, value int32) int32 {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.put(key, value)
}

func (this *IntIntMap) put(key int32, value int32) int32 {

	tab := this.table
	index := this.hash(key) % uint(len(tab))
	for e := tab[index]; e != nil; e = e.next {
		if e.key == key {
			old := e.value
			e.value = value
			return old
		}
	}

	if this.count >= this.threshold {
		this.rehash()
		tab = this.table
		index = this.hash(key) % uint(len(tab))
	}
	e := &IntIntEntry{key: key, value: value, next: tab[index]}
	tab[index] = e
	this.count++
	return this.NONE
}

func (this *IntIntMap) Add(key int32, value int32) int32 {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.add(key, value)
}

func (this *IntIntMap) add(key int32, value int32) int32 {
	tab := this.table
	index := this.hash(key) % uint(len(tab))
	for e := tab[index]; e != nil; e = e.next {
		if e.key == key {
			old := e.value
			e.value += value
			return old
		}
	}

	if this.count >= this.threshold {
		this.rehash()
		tab = this.table
		index = this.hash(key) % uint(len(tab))
	}
	e := &IntIntEntry{key: key, value: value, next: tab[index]}
	tab[index] = e
	this.count++
	return value
}

func (this *IntIntMap) AddIfExist(key int32, value int32) int32 {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.addIfExist(key, value)
}

func (this *IntIntMap) addIfExist(key int32, value int32) int32 {
	tab := this.table
	index := this.hash(key) % uint(len(tab))
	for e := tab[index]; e != nil; e = e.next {
		if e.key == key {
			e.value += value
			return e.value
		}
	}
	return 0
}

// 여기까지 작업

func (this *IntIntMap) Remove(key int32) int32 {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.remove(key)
}

func (this *IntIntMap) remove(key int32) int32 {

	tab := this.table
	index := this.hash(key) % uint(len(tab))
	e := tab[index]
	var prev *IntIntEntry = nil
	for e != nil {
		if e.key == key {
			if prev != nil {
				prev.next = e.next
			} else {
				tab[index] = e.next
			}
			this.count--
			oldValue := e.value
			e.value = this.NONE
			return oldValue
		}
		prev = e
		e = e.next
	}
	return this.NONE
}

func (this *IntIntMap) IsEmpty() bool {
	return this.count == 0
}
func (this *IntIntMap) IsFull() bool {
	return this.max > 0 && this.max <= this.count
}

func (this *IntIntMap) Clear() {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.clear()
}
func (this *IntIntMap) clear() {
	if this.count == 0 {
		return
	}

	tab := this.table
	for index := len(tab) - 1; index >= 0; index-- {
		tab[index] = nil
	}
	this.count = 0
}

type IntIntMapSortable struct {
	compare func(a, b int32) bool
	data    []*IntIntEntry
}

func (this IntIntMapSortable) Len() int {
	return len(this.data)
}
func (this IntIntMapSortable) Less(i, j int) bool {
	return this.compare(this.data[i].GetKey(), this.data[j].GetKey())
}

func (this IntIntMapSortable) Swap(i, j int) {
	this.data[i], this.data[j] = this.data[j], this.data[i]
}

func (this *IntIntMap) Sort(c func(k1, k2 int32) bool) {
	this.lock.Lock()
	defer this.lock.Unlock()

	sz := this.Size()
	list := make([]*IntIntEntry, sz)
	en := this.Entries()
	for i := 0; i < sz; i++ {
		list[i] = en.NextElement().(*IntIntEntry)
	}
	sort.Sort(IntIntMapSortable{compare: c, data: list})

	this.clear()
	for i := 0; i < sz; i++ {
		this.put(list[i].GetKey(), list[i].GetValue())
	}
}

func (this *IntIntMap) ToBytes(dout *io.DataOutputX) {
	dout.WriteDecimal(int64(this.Size()))
	en := this.Entries()
	for en.HasMoreElements() {
		e := en.NextElement().(*IntIntEntry)
		dout.WriteDecimal(int64(e.GetKey()))
		dout.WriteDecimal(int64(e.GetValue()))
	}
}

func (this *IntIntMap) ToObject(din *io.DataInputX) *IntIntMap {
	cnt := int(din.ReadDecimal())
	for i := 0; i < cnt; i++ {
		key := int32(din.ReadDecimal())
		value := int32(din.ReadDecimal())
		this.Put(key, value)
	}
	return this
}
func (this *IntIntMap) valueSum() int {
	sum := 0
	en := this.Values()
	for en.HasMoreElements() {
		sum += int(en.NextInt())
	}
	return sum
}

func (this *IntIntMap) ToString() string {
	this.lock.Lock()
	defer this.lock.Unlock()

	var buffer bytes.Buffer
	x := this.Entries()
	buffer.WriteString("{")
	for i := 0; x.HasMoreElements(); i++ {
		if i > 0 {
			buffer.WriteString(", ")
		}
		e := x.NextElement().(*IntIntEntry)
		buffer.WriteString(e.ToString())
	}
	buffer.WriteString("}")
	return buffer.String()
}

type IntIntMapEnumer struct {
	table        []*IntIntEntry
	index        int
	entry        *IntIntEntry
	lastReturned *IntIntEntry
	Type         int

	//	isKey  bool
	//	parent *IntIntMap
	//	entry  *IntIntEntry
}

func NewIntIntMapEnumer(t int) *IntIntMapEnumer {
	p := new(IntIntMapEnumer)
	p.Type = t
	return p
}

func (this *IntIntMapEnumer) HasMoreElements() bool {
	for this.entry == nil && this.index > 0 {
		this.index--
		this.entry = this.table[this.index]
	}
	return this.entry != nil
}

func (this *IntIntMapEnumer) NextElement() interface{} {
	for this.entry == nil && this.index > 0 {
		this.index--
		this.entry = this.table[this.index]
	}

	if this.entry != nil {
		this.lastReturned = this.entry
		e := this.lastReturned

		this.entry = e.next
		switch this.Type {
		case ELEMENT_TYPE_KEYS:
			return e.key
		case ELEMENT_TYPE_VALUES:
			return e.value
		default:
			return e
		}
	}
	panic("no more next")

}

func (this *IntIntMapEnumer) NextInt() int32 {
	for this.entry == nil && this.index > 0 {
		this.index--
		this.entry = this.table[this.index]
	}

	if this.entry != nil {
		this.lastReturned = this.entry
		e := this.lastReturned

		this.entry = e.next
		switch this.Type {
		case ELEMENT_TYPE_KEYS:
			return e.key
		case ELEMENT_TYPE_VALUES:
			return e.value
		default:
			return 0
		}
	}
	panic("no more next int")
}
