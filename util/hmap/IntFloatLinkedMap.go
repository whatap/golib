package hmap

import (
	"bytes"
	"sort"
	"sync"

	"github.com/whatap/golib/io"
)

type IntFloatLinkedMap struct {
	table      []*IntFloatLinkedEntry
	header     *IntFloatLinkedEntry
	count      int
	threshold  int
	loadFactor float32
	NONE       float32
	lock       sync.Mutex
	max        int
}

func NewIntFloatLinkedMap() *IntFloatLinkedMap {

	initCapacity := DEFAULT_CAPACITY
	loadFactor := DEFAULT_LOAD_FACTOR

	this := new(IntFloatLinkedMap)
	this.loadFactor = float32(loadFactor)
	this.table = make([]*IntFloatLinkedEntry, initCapacity)
	this.header = &IntFloatLinkedEntry{}
	this.header.link_next = this.header
	this.header.link_prev = this.header
	this.threshold = int(float64(initCapacity) * loadFactor)
	return this
}

func (this *IntFloatLinkedMap) Size() int {
	return this.count
}

func (this *IntFloatLinkedMap) KeyArray() []int32 {
	this.lock.Lock()
	defer this.lock.Unlock()

	_keys := make([]int32, this.Size())
	en := this.Keys()
	for i := 0; i < len(_keys); i++ {
		_keys[i] = en.NextInt()
	}
	return _keys
}

type Enumer struct {
	parent *IntFloatLinkedMap
	entry  *IntFloatLinkedEntry
}

func (this *Enumer) HasMoreElements() bool {
	return this.entry != nil && this.parent.header != this.entry
}
func (this *Enumer) NextInt() int32 {
	if this.HasMoreElements() {
		e := this.entry
		this.entry = e.link_next
		return e.key
	}
	return 0
}
func (this *Enumer) NextFloat() float32 {
	if this.HasMoreElements() {
		e := this.entry
		this.entry = e.link_next
		return e.value
	}
	return this.parent.NONE
}
func (this *Enumer) NextElement() interface{} {
	if this.HasMoreElements() {
		e := this.entry
		this.entry = e.link_next
		return e
	}
	return nil
}
func (this *IntFloatLinkedMap) Keys() IntEnumer {
	return &Enumer{parent: this, entry: this.header.link_next}
}
func (this *IntFloatLinkedMap) Values() FloatEnumer {
	return &Enumer{parent: this, entry: this.header.link_next}
}
func (this *IntFloatLinkedMap) Entries() Enumeration {
	return &Enumer{parent: this, entry: this.header.link_next}
}

func (this *IntFloatLinkedMap) ContainsValue(value float32) bool {
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
func (this *IntFloatLinkedMap) ContainsKey(key int32) bool {
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
func (this *IntFloatLinkedMap) Get(key int32) float32 {
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
func (this *IntFloatLinkedMap) GetFirstKey() int32 {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.header.link_next.key
}

func (this *IntFloatLinkedMap) GetLastKey() int32 {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.header.link_prev.key
}

func (this *IntFloatLinkedMap) GetFirstValue() float32 {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.header.link_next.value
}

func (this *IntFloatLinkedMap) GetLastValue() float32 {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.header.link_prev.value
}
func (this *IntFloatLinkedMap) hash(key int32) uint {
	return uint(key)
}

func (this *IntFloatLinkedMap) rehash() {
	oldCapacity := len(this.table)
	oldMap := this.table
	newCapacity := oldCapacity*2 + 1
	newMap := make([]*IntFloatLinkedEntry, newCapacity)
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

func (this *IntFloatLinkedMap) SetMax(max int) *IntFloatLinkedMap {
	this.max = max
	return this
}
func (this *IntFloatLinkedMap) Put(key int32, value float32) float32 {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.put(key, value, PUT_LAST)
}
func (this *IntFloatLinkedMap) PutLast(key int32, value float32) float32 {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.put(key, value, PUT_FORCE_LAST)
}
func (this *IntFloatLinkedMap) PutFirst(key int32, value float32) float32 {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.put(key, value, PUT_FORCE_FIRST)
}
func (this *IntFloatLinkedMap) put(key int32, value float32, m PUT_MODE) float32 {

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
	e := &IntFloatLinkedEntry{key: key, value: value, hash_next: tab[index]}
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

func (this *IntFloatLinkedMap) Add(key int32, value float32) float32 {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.add(key, value, PUT_LAST)
}
func (this *IntFloatLinkedMap) AddLast(key int32, value float32) float32 {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.add(key, value, PUT_FORCE_LAST)
}
func (this *IntFloatLinkedMap) AddFirst(key int32, value float32) float32 {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.add(key, value, PUT_FORCE_FIRST)
}

func (this *IntFloatLinkedMap) add(key int32, value float32, m PUT_MODE) float32 {
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
	e := &IntFloatLinkedEntry{key: key, value: value, hash_next: tab[index]}
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
func (this *IntFloatLinkedMap) Remove(key int32) float32 {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.remove(key)
}
func (this *IntFloatLinkedMap) RemoveFirst() float32 {
	if this.IsEmpty() {
		return 0
	}
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.remove(this.header.link_next.key)
}

func (this *IntFloatLinkedMap) RemoveLast() float32 {
	if this.IsEmpty() {
		return 0
	}
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.remove(this.header.link_prev.key)
}

func (this *IntFloatLinkedMap) remove(key int32) float32 {

	tab := this.table
	index := this.hash(key) % uint(len(tab))
	e := tab[index]
	var prev *IntFloatLinkedEntry = nil
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

func (this *IntFloatLinkedMap) IsEmpty() bool {
	return this.count == 0
}
func (this *IntFloatLinkedMap) IsFull() bool {
	return this.max > 0 && this.max <= this.count
}

func (this *IntFloatLinkedMap) Clear() {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.clear()
}
func (this *IntFloatLinkedMap) clear() {
	tab := this.table
	for index := len(tab) - 1; index >= 0; index-- {
		tab[index] = nil
	}
	this.header.link_next = this.header
	this.header.link_prev = this.header
	this.count = 0
}

func (this *IntFloatLinkedMap) chain(link_prev *IntFloatLinkedEntry, link_next *IntFloatLinkedEntry, e *IntFloatLinkedEntry) {
	e.link_prev = link_prev
	e.link_next = link_next
	link_prev.link_next = e
	link_next.link_prev = e
}

func (this *IntFloatLinkedMap) unchain(e *IntFloatLinkedEntry) {
	e.link_prev.link_next = e.link_next
	e.link_next.link_prev = e.link_prev
	e.link_prev = nil
	e.link_next = nil
}

type intFloatSortable struct {
	compare func(a, b int32) bool
	data    []*IntFloatLinkedEntry
}

func (this intFloatSortable) Len() int {
	return len(this.data)
}
func (this intFloatSortable) Less(i, j int) bool {
	return this.compare(this.data[i].GetKey(), this.data[j].GetKey())
}

func (this intFloatSortable) Swap(i, j int) {
	this.data[i], this.data[j] = this.data[j], this.data[i]
}

func (this *IntFloatLinkedMap) Sort(c func(k1, k2 int32) bool) {
	this.lock.Lock()
	defer this.lock.Unlock()

	sz := this.Size()
	list := make([]*IntFloatLinkedEntry, sz)
	en := this.Entries()
	for i := 0; i < sz; i++ {
		list[i] = en.NextElement().(*IntFloatLinkedEntry)
	}
	sort.Sort(intFloatSortable{compare: c, data: list})

	this.clear()
	for i := 0; i < sz; i++ {
		this.put(list[i].GetKey(), list[i].GetValue(), PUT_LAST)
	}
}

func (this *IntFloatLinkedMap) ToBytes(dout *io.DataOutputX) {
	dout.WriteDecimal(int64(this.Size()))
	en := this.Entries()
	for en.HasMoreElements() {
		e := en.NextElement().(IntFloatLinkedEntry)
		dout.WriteDecimal(int64(e.GetKey()))
		dout.WriteFloat(e.GetValue())
	}
}

func (this *IntFloatLinkedMap) ToObject(din *io.DataInputX) *IntFloatLinkedMap {
	cnt := int(din.ReadDecimal())
	for i := 0; i < cnt; i++ {
		key := int32(din.ReadDecimal())
		value := din.ReadFloat()
		this.Put(key, value)
	}
	return this
}
func (this *IntFloatLinkedMap) ToString() string {
	this.lock.Lock()
	defer this.lock.Unlock()

	var buffer bytes.Buffer
	x := this.Entries()
	buffer.WriteString("{")
	for i := 0; x.HasMoreElements(); i++ {
		if i > 0 {
			buffer.WriteString(", ")
		}
		e := x.NextElement().(*IntFloatLinkedEntry)
		buffer.WriteString(e.ToString())
	}
	buffer.WriteString("}")
	return buffer.String()
}
