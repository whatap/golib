package hmap

import (
	"bytes"
	"sort"
	"sync"

	"github.com/whatap/golib/io"
)

type IntIntLinkedMap struct {
	table      []*IntIntLinkedEntry
	header     *IntIntLinkedEntry
	count      int
	threshold  int
	loadFactor float32
	NONE       int32
	lock       sync.Mutex
	max        int
}

func NewIntIntLinkedMap() *IntIntLinkedMap {

	initCapacity := DEFAULT_CAPACITY
	loadFactor := DEFAULT_LOAD_FACTOR

	this := new(IntIntLinkedMap)
	this.loadFactor = float32(loadFactor)
	this.table = make([]*IntIntLinkedEntry, initCapacity)
	this.header = &IntIntLinkedEntry{}
	this.header.link_next = this.header
	this.header.link_prev = this.header
	this.threshold = int(float64(initCapacity) * loadFactor)

	//fmt.Printf("this=%p, threshold=%d \r\n", this, this.threshold)

	return this
}

func (this *IntIntLinkedMap) Size() int {
	return this.count
}

func (this *IntIntLinkedMap) KeyArray() []int32 {
	this.lock.Lock()
	defer this.lock.Unlock()

	_keys := make([]int32, this.Size())
	en := this.Keys()
	for i := 0; i < len(_keys); i++ {
		_keys[i] = en.NextInt()
	}
	return _keys
}

type IntIntEnumer struct {
	isKey  bool
	parent *IntIntLinkedMap
	entry  *IntIntLinkedEntry
}

func (this *IntIntEnumer) HasMoreElements() bool {
	return this.entry != nil && this.parent.header != this.entry
}
func (this *IntIntEnumer) NextInt() int32 {
	if this.HasMoreElements() {
		e := this.entry
		this.entry = e.link_next
		if this.isKey {
			return e.key
		} else {
			return e.value
		}
	}
	return 0
}
func (this *IntIntEnumer) NextElement() interface{} {
	if this.HasMoreElements() {
		e := this.entry
		this.entry = e.link_next
		return e
	}
	return nil
}
func (this *IntIntLinkedMap) Keys() IntEnumer {
	return &IntIntEnumer{isKey: true, parent: this, entry: this.header.link_next}
}
func (this *IntIntLinkedMap) Values() IntEnumer {
	return &IntIntEnumer{isKey: false, parent: this, entry: this.header.link_next}
}
func (this *IntIntLinkedMap) Entries() Enumeration {
	return &IntIntEnumer{parent: this, entry: this.header.link_next}
}

func (this *IntIntLinkedMap) ContainsValue(value int32) bool {
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
func (this *IntIntLinkedMap) ContainsKey(key int32) bool {
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
func (this *IntIntLinkedMap) Get(key int32) int32 {
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
func (this *IntIntLinkedMap) GetFirstKey() int32 {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.header.link_next.key
}

func (this *IntIntLinkedMap) GetLastKey() int32 {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.header.link_prev.key
}

func (this *IntIntLinkedMap) GetFirstValue() int32 {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.header.link_next.value
}

func (this *IntIntLinkedMap) GetLastValue() int32 {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.header.link_prev.value
}
func (this *IntIntLinkedMap) hash(key int32) uint {
	return uint(key)
}

func (this *IntIntLinkedMap) rehash() {
	oldCapacity := len(this.table)
	oldMap := this.table
	newCapacity := oldCapacity*2 + 1
	newMap := make([]*IntIntLinkedEntry, newCapacity)

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

func (this *IntIntLinkedMap) SetMax(max int) *IntIntLinkedMap {
	this.max = max
	return this
}
func (this *IntIntLinkedMap) Put(key int32, value int32) int32 {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.put(key, value, PUT_LAST)
}
func (this *IntIntLinkedMap) PutLast(key int32, value int32) int32 {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.put(key, value, PUT_FORCE_LAST)
}
func (this *IntIntLinkedMap) PutFirst(key int32, value int32) int32 {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.put(key, value, PUT_FORCE_FIRST)
}
func (this *IntIntLinkedMap) put(key int32, value int32, m PUT_MODE) int32 {

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
	e := &IntIntLinkedEntry{key: key, value: value, hash_next: tab[index]}
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

func (this *IntIntLinkedMap) Add(key int32, value int32) int32 {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.add(key, value, PUT_LAST)
}

func (this *IntIntLinkedMap) AddNoOver(key int32, value int32) int32 {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.addNoOver(key, value, PUT_LAST)
}

func (this *IntIntLinkedMap) AddLast(key int32, value int32) int32 {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.add(key, value, PUT_FORCE_LAST)
}
func (this *IntIntLinkedMap) AddFirst(key int32, value int32) int32 {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.add(key, value, PUT_FORCE_FIRST)
}

func (this *IntIntLinkedMap) add(key int32, value int32, m PUT_MODE) int32 {
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
	e := &IntIntLinkedEntry{key: key, value: value, hash_next: tab[index]}
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

func (this *IntIntLinkedMap) addNoOver(key int32, value int32, m PUT_MODE) int32 {
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

	//MAX에 도달하면 더 이상 키를 추가하지 않는다.
	if this.max > 0 && this.count >= this.max {
		return this.NONE
	}

	if this.count >= this.threshold {
		this.rehash()
		tab = this.table
		index = this.hash(key) % uint(len(tab))
	}
	e := &IntIntLinkedEntry{key: key, value: value, hash_next: tab[index]}
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

func (this *IntIntLinkedMap) Remove(key int32) int32 {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.remove(key)
}
func (this *IntIntLinkedMap) RemoveFirst() int32 {
	if this.IsEmpty() {
		return 0
	}
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.remove(this.header.link_next.key)
}

func (this *IntIntLinkedMap) RemoveLast() int32 {
	if this.IsEmpty() {
		return 0
	}
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.remove(this.header.link_prev.key)
}

func (this *IntIntLinkedMap) remove(key int32) int32 {

	tab := this.table
	index := this.hash(key) % uint(len(tab))
	e := tab[index]
	var prev *IntIntLinkedEntry = nil
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

func (this *IntIntLinkedMap) IsEmpty() bool {
	return this.count == 0
}
func (this *IntIntLinkedMap) IsFull() bool {
	return this.max > 0 && this.max <= this.count
}

func (this *IntIntLinkedMap) Clear() {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.clear()
}
func (this *IntIntLinkedMap) clear() {
	tab := this.table
	for index := len(tab) - 1; index >= 0; index-- {
		tab[index] = nil
	}
	this.header.link_next = this.header
	this.header.link_prev = this.header
	this.count = 0
}

func (this *IntIntLinkedMap) chain(link_prev *IntIntLinkedEntry, link_next *IntIntLinkedEntry, e *IntIntLinkedEntry) {
	e.link_prev = link_prev
	e.link_next = link_next
	link_prev.link_next = e
	link_next.link_prev = e
}

func (this *IntIntLinkedMap) unchain(e *IntIntLinkedEntry) {
	e.link_prev.link_next = e.link_next
	e.link_next.link_prev = e.link_prev
	e.link_prev = nil
	e.link_next = nil
}

type intIntSortable struct {
	compare func(a, b int32) bool
	data    []*IntIntLinkedEntry
}

func (this intIntSortable) Len() int {
	return len(this.data)
}
func (this intIntSortable) Less(i, j int) bool {
	return this.compare(this.data[i].GetKey(), this.data[j].GetKey())
}

func (this intIntSortable) Swap(i, j int) {
	this.data[i], this.data[j] = this.data[j], this.data[i]
}

func (this *IntIntLinkedMap) Sort(c func(k1, k2 int32) bool) {
	this.lock.Lock()
	defer this.lock.Unlock()

	sz := this.Size()
	list := make([]*IntIntLinkedEntry, sz)
	en := this.Entries()
	for i := 0; i < sz; i++ {
		list[i] = en.NextElement().(*IntIntLinkedEntry)
	}
	sort.Sort(intIntSortable{compare: c, data: list})

	this.clear()
	for i := 0; i < sz; i++ {
		this.put(list[i].GetKey(), list[i].GetValue(), PUT_LAST)
	}
}

func (this *IntIntLinkedMap) ToBytes(dout *io.DataOutputX) {
	dout.WriteDecimal(int64(this.Size()))
	en := this.Entries()
	for en.HasMoreElements() {
		e := en.NextElement().(*IntIntLinkedEntry)
		dout.WriteDecimal(int64(e.GetKey()))
		dout.WriteDecimal(int64(e.GetValue()))
	}
}

func (this *IntIntLinkedMap) ToObject(din *io.DataInputX) *IntIntLinkedMap {
	cnt := int(din.ReadDecimal())
	for i := 0; i < cnt; i++ {
		key := int32(din.ReadDecimal())
		value := int32(din.ReadDecimal())
		this.Put(key, value)
	}
	return this
}
func (this *IntIntLinkedMap) ToString() string {
	this.lock.Lock()
	defer this.lock.Unlock()

	var buffer bytes.Buffer
	x := this.Entries()
	buffer.WriteString("{")
	for i := 0; x.HasMoreElements(); i++ {
		if i > 0 {
			buffer.WriteString(", ")
		}
		e := x.NextElement().(*IntIntLinkedEntry)
		buffer.WriteString(e.ToString())
	}
	buffer.WriteString("}")
	return buffer.String()
}
