package hmap

import (
	"fmt"
	//"log"
	"math"
	"sync"

	"github.com/whatap/golib/util/stringutil"
)

type IntKeyMap struct {
	// []*IntKeyEntry
	table []*IntKeyEntry
	// int
	count int
	// int
	threshold int
	// float32
	loadFactor float32
	// sync.Mutex
	lock sync.Mutex

	// create func(int32) interface{]// create, intern 함수 사용 안함
	//Create func(int32) interface{}
}

func NewIntKeyMapDefault() *IntKeyMap {

	p := NewIntKeyMap(DEFAULT_CAPACITY, DEFAULT_LOAD_FACTOR)

	return p
}

func NewIntKeyMap(initCapacity int, loadFactor float32) *IntKeyMap {
	defer func() {
		if r := recover(); r != nil {
			// TODO 추후 hmap 에서 recover 는 없애고 panic 처리. 호출하는 쪽에서 recover 할 것
			//logutil.Println("WA823", "Recover:",r)
			//return NewIntKeyMap(DEFAULT_CAPACITY, DEFAULT_LOAD_FACTOR)
		}
	}()

	p := new(IntKeyMap)

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
	p.table = make([]*IntKeyEntry, initCapacity)
	p.threshold = int(float32(initCapacity) * loadFactor)

	return p
}

func (this *IntKeyMap) Size() int {
	return this.count
}

func (this *IntKeyMap) Keys() IntEnumer {
	this.lock.Lock()
	defer this.lock.Unlock()

	return NewIntKeyEnumer(ELEMENT_TYPE_KEYS, this.table)
}

func (this *IntKeyMap) Values() Enumeration {
	//public synchronized Enumeration<V> values() {
	this.lock.Lock()
	defer this.lock.Unlock()
	return NewIntKeyEnumer(ELEMENT_TYPE_VALUES, this.table)
}

func (this *IntKeyMap) Entries() Enumeration {
	//public synchronized Enumeration<IntKeyEntry<V>> entries() {
	this.lock.Lock()
	defer this.lock.Unlock()
	return NewIntKeyEnumer(ELEMENT_TYPE_ENTRIES, this.table)
}

func (this *IntKeyMap) ContainsValue(value interface{}) bool {
	//public synchronized boolean containsValue(V value) {
	this.lock.Lock()
	defer this.lock.Unlock()
	if value == nil {
		return false
	}

	// TODO
	//IntKeyEntry<V> tab[] = this.table;
	//		tab := this.table
	//		i := len(tab)
	//		for i> 0 ; i-- {
	//			//for (IntKeyEntry<V> e = tab[i]; e != null; e = e.next) {
	//			for e := tab[i]; e != nil; e = e.Next() {
	//				if CompareUtil.equals(e.value, value) {
	//					return true
	//				}
	//			}
	//		}
	return false
}

func (this *IntKeyMap) ContainsKey(key int32) bool {
	this.lock.Lock()
	defer this.lock.Unlock()

	tab := this.table
	index := int(this.hash(key)) % len(tab)

	//for (IntKeyEntry<V> e = tab[index]; e != null; e = e.next) {
	for e := tab[index]; e != nil; e = e.Next {
		if e.Key == key {
			return true
		}
	}
	return false
}

func (this *IntKeyMap) hash(h int32) uint {
	ret := uint(h)
	// TODO
	ret ^= (uint(h) >> 20) ^ (uint(h) >> 12)
	ret = ret ^ (uint(h) >> 7) ^ (uint(h) >> 4)

	return ret & uint(math.MaxInt32)
}

func (this *IntKeyMap) Get(key int32) interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()

	tab := this.table
	index := int(this.hash(key)) % len(tab)

	//for (IntKeyEntry<V> e = tab[index]; e != null; e = e.next) {
	for e := tab[index]; e != nil; e = e.Next {
		if e.Key == key {
			return e.Value
		}
	}
	return nil
}

func (this *IntKeyMap) rehash() {
	oldCapacity := len(this.table)
	oldMap := this.table
	newCapacity := oldCapacity*2 + 1

	newMap := make([]*IntKeyEntry, newCapacity)
	this.threshold = int(float32(newCapacity) * this.loadFactor)
	this.table = newMap

	for i := oldCapacity; i > 0; i-- {
		for old := oldMap[i-1]; old != nil; {
			e := old
			old = old.Next
			index := int(this.hash(e.Key)) % newCapacity
			e.Next = newMap[index]
			newMap[index] = e
		}
	}
}

func (this *IntKeyMap) KeyArray() []int32 {
	this.lock.Lock()
	defer this.lock.Unlock()

	_keys := make([]int32, this.Size())

	//IntEnumer en = this.keys();
	en := this.Keys()

	for i := 0; i < len(_keys); i++ {
		_keys[i] = en.NextInt()
	}
	return _keys
}

func (this *IntKeyMap) Put(key int32, value interface{}) interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()

	tab := this.table
	_hash := this.hash(key)
	index := _hash % uint(len(tab))

	e := tab[index]
	for e = tab[index]; e != nil; e = e.Next {
		if e.Key == key {
			old := e.Value
			e.Value = value
			return old
		}
	}

	if this.count >= this.threshold {
		this.rehash()
		tab = this.table
		index = _hash % uint(len(tab))
	}
	e = NewIntKeyEntry(key, value, tab[index])
	tab[index] = e
	this.count++

	return nil
}

//  intern 함수 사용 안함
//func (this *IntKeyMap) Intern(key int32) interface{} {
//	this.lock.Lock()
//	defer this.lock.Unlock()
//	tab := this.table
//	_hash := this.hash(key)
//	index := _hash % uint(len(tab))
//
//	for e := tab[index]; e != nil; e = e.Next {
//		if e.Key == key {
//			return e.Value
//		}
//	}
//
//	//var value interface{}
//	value := this.create(key)
//
//	if value == nil {
//		return nil
//	}
//
//	if this.count >= this.threshold {
//		this.rehash()
//		tab = this.table
//		index = _hash % uint(len(tab))
//	}
//	e := NewIntKeyEntry(key, value, tab[index])
//	tab[index] = e
//	this.count++
//	return value
//}

//Override create , intern 함수 사용 안함
//func (this *IntKeyMap) create(key int32) interface{} {
//	if this.Create == nil {
//		panic("Error IntKeyMap.Count is nil")
//	}
//
//	return this.Create(key)
//	//throw new RuntimeException("not implemented create()")
//	//return interface{}
//	//return nil
//}

func (this *IntKeyMap) Remove(key int32) interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()

	tab := this.table
	index := this.hash(key) % uint(len(tab))
	var prev *IntKeyEntry
	prev = nil
	//for e := tab[index], prev := nil; e != nil; prev = e, e = e.Next) {
	for e := tab[index]; e != nil; {
		if e.Key == key {
			if prev != nil {
				prev.Next = e.Next
			} else {
				tab[index] = e.Next
			}
			this.count--
			oldValue := e.Value
			e.Value = nil
			return oldValue
		}

		prev = e
		e = e.Next
	}
	return nil
}

func (this *IntKeyMap) Clear() {
	this.lock.Lock()
	defer this.lock.Unlock()
	tab := this.table

	for index := len(tab) - 1; index >= 0; index-- {
		tab[index] = nil
	}
	this.count = 0

}

func (this *IntKeyMap) ToString() string {
	buf := stringutil.NewStringBuffer()

	it := this.Entries()
	buf.Append("{")

	for i := 0; it.HasMoreElements(); i++ {
		e := it.NextElement().(*IntKeyEntry)

		if i > 0 {
			buf.Append(", ")
		}
		buf.Append(e.ToString())
	}
	buf.Append("}")

	return buf.ToString()
}

func (this *IntKeyMap) ToFormatString() string {
	buf := stringutil.NewStringBuffer()

	it := this.Entries()
	buf.Append("{\n")
	for it.HasMoreElements() {
		e := it.NextElement().(*IntKeyEntry)
		buf.Append("\t").Append(e.ToString()).Append("\n")
	}
	buf.Append("}")
	return buf.ToString()
}

func (this *IntKeyMap) PutAll(other *IntKeyMap) {
	if other == nil {
		return
	}
	it := other.Entries()

	for it.HasMoreElements() {
		e := it.NextElement().(*IntKeyEntry)
		this.Put(e.GetKey(), e.GetValue())
	}
}

type IntKeyEnumer struct {
	table        []*IntKeyEntry
	entry        *IntKeyEntry
	index        int
	lastReturned *IntKeyEntry
	Type         int
}

func NewIntKeyEnumer(Type int, table []*IntKeyEntry) *IntKeyEnumer {
	p := new(IntKeyEnumer)
	p.table = table
	p.index = len(table)
	p.entry = nil
	p.lastReturned = nil
	p.Type = Type
	return p
}

func (this *IntKeyEnumer) HasMoreElements() bool {
	for this.entry == nil && this.index > 0 {
		this.index--
		this.entry = this.table[this.index]
	}
	//fmt.Println("IntKeyMap.IntKeyEnumer HasMoreElements index=", this.index, ",return =", (this.entry != nil))
	return (this.entry != nil)
}

func (this *IntKeyEnumer) NextElement() interface{} {
	//fmt.Println("IntKeyMap.IntKeyEnumer NextElement index=", this.index, ",entry =", this.entry)

	for this.entry == nil && this.index > 0 {
		this.index--
		this.entry = this.table[this.index]
	}

	if this.entry != nil {
		//fmt.Println("IntKeyMap.IntKeyEnumer NextElement this.entry != nil , index=", this.index, ",entry =", this.entry)
		this.lastReturned = this.entry
		e := this.lastReturned

		this.entry = e.Next

		switch this.Type {
		case ELEMENT_TYPE_KEYS:
			return e.Key
		case ELEMENT_TYPE_VALUES:
			return e.Value
		default:
			return e
		}
	}
	//fmt.Println("IntKeyMap.IntKeyEnumer NextElement for , index=", this.index, ",entry =", this.entry)

	//throw new NoSuchElementException("no more next");
	panic("Panic NextElement no more next")

}

func (this *IntKeyEnumer) NextInt() int32 {
	for this.entry == nil && this.index > 0 {
		this.index--
		this.entry = this.table[this.index]
	}

	if this.entry != nil {
		this.lastReturned = this.entry
		e := this.lastReturned
		this.entry = e.Next
		return e.Key
	}
	//throw new NoSuchElementException("no more next");
	panic("Panic NextInt no more next")

}
