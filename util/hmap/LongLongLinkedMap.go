package hmap

import (
	"fmt"
	//"log"
	"bytes"
	"sort"
	"sync"

	"github.com/whatap/golib/io"
)

type LongLongLinkedMap struct {
	table  []*LongLongLinkedEntry
	header *LongLongLinkedEntry

	count      int
	threshold  int
	loadFactor float32
	NONE       int64

	max int

	lock sync.Mutex
}

func NewLongLongLinkedMap(initCapacity int, loadFactor float32) *LongLongLinkedMap {
	defer func() {
		if r := recover(); r != nil {
			// TODO 추후 hmap 에서 recover 는 없애고 panic 처리. 호출하는 쪽에서 recover 할 것
			//logutil.Println("WA824", "NewLongLongLinkedMap Recover", r)
			//return NewIntKeyMap(DEFAULT_CAPACITY, DEFAULT_LOAD_FACTOR)
		}
	}()

	p := new(LongLongLinkedMap)

	if initCapacity < 0 {
		panic(fmt.Sprintf("Capacity Error: %d", initCapacity))
		//throw new RuntimeException("Capacity Error: " + initCapacity);
	}
	if loadFactor <= 0 {
		panic(fmt.Sprintf("Load Count Error: %f", loadFactor))
		//throw new RuntimeException("Load Count Error: " + loadFactor);
	}
	if initCapacity == 0 {
		initCapacity = 1
	}
	p.loadFactor = loadFactor
	p.table = make([]*LongLongLinkedEntry, initCapacity)
	p.header = NewLongLongLinkedEntry(0, 0, nil)
	p.header.link_prev = p.header
	p.header.link_next = p.header.link_prev

	p.threshold = int(float32(initCapacity) * loadFactor)

	p.lock = sync.Mutex{}
	return p
}

func NewLongLongLinkedMapDefault() *LongLongLinkedMap {
	p := NewLongLongLinkedMap(DEFAULT_CAPACITY, DEFAULT_LOAD_FACTOR)
	return p
}

func (this *LongLongLinkedMap) SetNullValue(none int64) *LongLongLinkedMap {
	this.NONE = none
	return this
}

func (this *LongLongLinkedMap) Size() int {
	return this.count
}

func (this *LongLongLinkedMap) KeyArray() []int64 {
	_keys := make([]int64, this.Size())
	en := this.Keys()
	for i := 0; i < len(_keys); i++ {
		_keys[i] = en.NextLong()
	}
	return _keys
}

func (this *LongLongLinkedMap) Keys() LongEnumer {
	return &LongLongLinkedEnumer{isKey: true, parent: this, entry: this.header.link_next}
}
func (this *LongLongLinkedMap) Values() LongEnumer {
	return &LongLongLinkedEnumer{isKey: false, parent: this, entry: this.header.link_next}
}
func (this *LongLongLinkedMap) Entries() Enumeration {
	return &LongLongLinkedEnumer{parent: this, entry: this.header.link_next}
}

//func (this *LongLongLinkedMap) Keys() LongEnumer {
//	this.lock.Lock()
//	defer this.lock.Unlock()
//	return NewLongLongLinkedEnumer(this, this.header.link_next, ELEMENT_TYPE_KEYS)
//	//return &LongLongLinkedEnumer{isKey: true, parent: this, entry: this.header.link_next}
//}
//
//func (this *LongLongLinkedMap) Values() LongEnumer {
//	this.lock.Lock()
//	defer this.lock.Unlock()
//	return NewLongLongLinkedEnumer(this, this.header.link_next, ELEMENT_TYPE_VALUES)
//	//return &LongLongLinkedEnumer{isKey: false, parent: this, entry: this.header.link_next}
//}
//
//func (this *LongLongLinkedMap) Entries() Enumeration {
//	this.lock.Lock()
//	defer this.lock.Unlock()
//	return NewLongLongLinkedEnumer(this, this.header.link_next, ELEMENT_TYPE_ENTRIES)
//	////return &LongLongLinkedEnumer{parent: this, entry: this.header.link_next}
//}

func (this *LongLongLinkedMap) ContainsValue(value int64) bool {
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

func (this *LongLongLinkedMap) ContainsKey(key int64) bool {
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

func (this *LongLongLinkedMap) Get(key int64) int64 {
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

func (this *LongLongLinkedMap) GetFirstKey() int64 {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.header.link_next.key
}

func (this *LongLongLinkedMap) GetLastKey() int64 {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.header.link_prev.key
}
func (this *LongLongLinkedMap) GetFirstValue() int64 {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.header.link_next.value
}

func (this *LongLongLinkedMap) GetLastValue() int64 {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.header.link_prev.value
}

func (this *LongLongLinkedMap) hash(key int64) uint {
	return uint(key)
	//return uint(key & math.MaxInt32)
	//return (int) (key ^ (key >>> 32)) & Integer.MAX_VALUE;
}

func (this *LongLongLinkedMap) rehash() {
	oldCapacity := len(this.table)
	oldMap := this.table
	newCapacity := oldCapacity*2 + 1
	newMap := make([]*LongLongLinkedEntry, newCapacity)
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

func (this *LongLongLinkedMap) SetMax(max int) *LongLongLinkedMap {
	this.max = max
	return this
}

func (this *LongLongLinkedMap) Put(key int64, value int64) int64 {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.put(key, value, PUT_LAST)
}
func (this *LongLongLinkedMap) PutLast(key int64, value int64) int64 {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.put(key, value, PUT_FORCE_LAST)
}
func (this *LongLongLinkedMap) PutFirst(key int64, value int64) int64 {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.put(key, value, PUT_FORCE_FIRST)
}

func (this *LongLongLinkedMap) put(key int64, value int64, m PUT_MODE) int64 {

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
				//overflowed(k, v);
			}
		case PUT_FORCE_LAST, PUT_LAST:
			for this.count >= this.max {
				// removeFirst();
				k := this.header.link_next.key
				this.remove(k)
				//overflowed(k, v);
			}
			break
		}
	}
	if this.count >= this.threshold {
		this.rehash()
		tab = this.table
		index = this.hash(key) % uint(len(tab))
	}
	e := &LongLongLinkedEntry{key: key, value: value, hash_next: tab[index]}
	//fmt.Println("LongLongLinkedMap Put======================>", e)
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
func (this *LongLongLinkedMap) Add(key int64, value int64) int64 {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.add(key, value, PUT_LAST)
}
func (this *LongLongLinkedMap) AddLast(key int64, value int64) int64 {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.add(key, value, PUT_FORCE_LAST)
}

func (this *LongLongLinkedMap) AddFirst(key int64, value int64) int64 {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.add(key, value, PUT_FORCE_FIRST)
}

func (this *LongLongLinkedMap) add(key int64, value int64, m PUT_MODE) int64 {
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
				//removeLast();
			}
		case PUT_FORCE_LAST, PUT_LAST:
			for this.count >= this.max {
				k := this.header.link_next.key
				this.remove(k)
				//removeFirst();
			}
			break
		}
	}
	if this.count >= this.threshold {
		this.rehash()
		tab = this.table
		index = this.hash(key) % uint(len(tab))
	}
	e := &LongLongLinkedEntry{key: key, value: value, hash_next: tab[index]}
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

func overflowed(k, v int64) {
	// TODO Auto-generated method stub
}

func (this *LongLongLinkedMap) Remove(key int64) int64 {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.remove(key)
}

func (this *LongLongLinkedMap) remove(key int64) int64 {

	tab := this.table
	index := this.hash(key) % uint(len(tab))
	e := tab[index]
	var prev *LongLongLinkedEntry = nil
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

func (this *LongLongLinkedMap) RemoveFirst() int64 {
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.IsEmpty() {
		return 0
	}
	return this.remove(this.header.link_next.key)
}

func (this *LongLongLinkedMap) RemoveLast() int64 {
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.IsEmpty() {
		return 0
	}

	return this.remove(this.header.link_prev.key)
}

func (this *LongLongLinkedMap) IsEmpty() bool {
	return this.count == 0
}
func (this *LongLongLinkedMap) IsFull() bool {
	return this.max > 0 && this.max <= this.count
}

func (this *LongLongLinkedMap) Clear() {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.clear()
}
func (this *LongLongLinkedMap) clear() {
	tab := this.table
	for index := len(tab) - 1; index >= 0; index-- {
		tab[index] = nil
	}
	this.header.link_next = this.header
	this.header.link_prev = this.header
	this.count = 0
}

func (this *LongLongLinkedMap) chain(link_prev *LongLongLinkedEntry, link_next *LongLongLinkedEntry, e *LongLongLinkedEntry) {
	e.link_prev = link_prev
	e.link_next = link_next
	link_prev.link_next = e
	link_next.link_prev = e
}

func (this *LongLongLinkedMap) unchain(e *LongLongLinkedEntry) {
	e.link_prev.link_next = e.link_next
	e.link_next.link_prev = e.link_prev
	e.link_prev = nil
	e.link_next = nil
}

type longLongSortable struct {
	compare func(a, b int64) bool
	data    []*LongLongLinkedEntry
}

func (this longLongSortable) Len() int {
	return len(this.data)
}
func (this longLongSortable) Less(i, j int) bool {
	return this.compare(this.data[i].GetKey(), this.data[j].GetKey())
}

func (this longLongSortable) Swap(i, j int) {
	this.data[i], this.data[j] = this.data[j], this.data[i]
}

func (this *LongLongLinkedMap) Sort(c func(k1, k2 int64) bool) {
	this.lock.Lock()
	defer this.lock.Unlock()

	sz := this.Size()
	list := make([]*LongLongLinkedEntry, sz)
	en := this.Entries()
	for i := 0; i < sz; i++ {
		list[i] = en.NextElement().(*LongLongLinkedEntry)
	}
	sort.Sort(longLongSortable{compare: c, data: list})

	this.clear()
	for i := 0; i < sz; i++ {
		this.put(list[i].GetKey(), list[i].GetValue(), PUT_LAST)
	}
}

//	public synchronized void sort(Comparator<LongLongLinkedEntry> c) {
//		ArrayList<LongLongLinkedEntry> list = new ArrayList<LongLongLinkedEntry>(this.size());
//		Enumeration<LongLongLinkedEntry> en = this.entries();
//		while (en.hasMoreElements()) {
//			list.add(en.nextElement());
//		}
//		Collections.sort(list, c);
//		this.clear();
//		for (int i = 0; i < list.size(); i++) {
//			LongLongLinkedEntry e = list.get(i);
//			this.put(e.getKey(), e.getValue());
//		}
//	}

func (this *LongLongLinkedMap) ToBytes(dout *io.DataOutputX) {
	dout.WriteDecimal(int64(this.Size()))
	en := this.Entries()
	for en.HasMoreElements() {
		e := en.NextElement().(*LongLongLinkedEntry)
		dout.WriteDecimal(int64(e.GetKey()))
		dout.WriteDecimal(int64(e.GetValue()))
	}
}

func (this *LongLongLinkedMap) ToObject(din *io.DataInputX) *LongLongLinkedMap {
	cnt := int(din.ReadDecimal())
	for i := 0; i < cnt; i++ {
		key := int64(din.ReadDecimal())
		value := int64(din.ReadDecimal())
		this.Put(key, value)
	}
	return this
}

func (this *LongLongLinkedMap) ToString() string {
	this.lock.Lock()
	defer this.lock.Unlock()

	var buffer bytes.Buffer
	x := this.Entries()
	buffer.WriteString("{")
	for i := 0; x.HasMoreElements(); i++ {
		if i > 0 {
			buffer.WriteString(", ")
		}
		e := x.NextElement().(LongLongLinkedEntry)
		buffer.WriteString(e.ToString())
	}
	buffer.WriteString("}")
	return buffer.String()
}

func LongLongLinkedMapMain() {
	m := NewLongLongLinkedMapDefault()
	// m.Put(1, 1)
	// m.Put(2, 1)
	// m.Put(3, 1)
	// m.Put(4, 1)
	//
	e := m.Keys()
	fmt.Println(e.NextLong())
	fmt.Println(e.NextLong())
	fmt.Println(e.NextLong())
	fmt.Println(e.NextLong())
	// fmt.Println(e.NextLong())

}

//	private static void print(Object e) {
//		System.out.println(e);
//	}

type LongLongLinkedEnumer struct {
	isKey  bool
	parent *LongLongLinkedMap
	entry  *LongLongLinkedEntry
	Type   int
}

func NewLongLongLinkedEnumer(parent *LongLongLinkedMap, entry *LongLongLinkedEntry, Type int) *LongLongLinkedEnumer {
	p := new(LongLongLinkedEnumer)
	p.parent = parent
	p.entry = entry
	p.Type = Type

	return p
}
func (this *LongLongLinkedEnumer) HasMoreElements() bool {
	return this.entry != nil && this.parent.header != this.entry
}
func (this *LongLongLinkedEnumer) NextLong() int64 {
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
func (this *LongLongLinkedEnumer) NextElement() interface{} {
	if this.HasMoreElements() {
		e := this.entry
		this.entry = e.link_next
		return e
	}
	return nil
}
