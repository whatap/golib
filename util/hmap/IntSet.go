package hmap

import (
	"bytes"
	"fmt"
	"sync"
	"strconv"
)

// TODO SearchPathMap 에서 사용하기 위해 생성 (현재는 생성만)
// 테스트 필요.
type IntSet struct {
	table      []*IntSetry
	count      int
	threshold  int
	loadFactor float32
	lock       sync.Mutex
	max        int
}

func NewIntSet() *IntSet {

	initCapacity := DEFAULT_CAPACITY
	loadFactor := DEFAULT_LOAD_FACTOR

	this := new(IntSet)
	this.loadFactor = float32(loadFactor)
	this.table = make([]*IntSetry, initCapacity)
	this.threshold = int(float64(initCapacity) * loadFactor)
	return this
}

func NewIntSetArray(arr []string) *IntSet {

	initCapacity := DEFAULT_CAPACITY
	loadFactor := DEFAULT_LOAD_FACTOR

	this := new(IntSet)
	this.loadFactor = float32(loadFactor)
	this.table = make([]*IntSetry, initCapacity)
	this.threshold = int(float64(initCapacity) * loadFactor)
	return this
}

func (this *IntSet) Size() int {
	return this.count
}

func (this *IntSet) Values() *IntSetEnumer {
	return &IntSetEnumer{table: this.table, index: len(this.table), entry: nil}
}

func (this *IntSet) Contains(key int32) bool {
	this.lock.Lock()
	defer this.lock.Unlock()
	tab := this.table
	index := uint(key) % uint(len(tab))
	for e := tab[index]; e != nil; e = e.next {
		if e.key == key {
			return true
		}
	}
	return false
}

func (this *IntSet) rehash() {
	oldCapacity := len(this.table)
	oldMap := this.table
	newCapacity := oldCapacity*2 + 1
	newMap := make([]*IntSetry, newCapacity)
	this.threshold = int(float32(newCapacity) * this.loadFactor)
	this.table = newMap
	for i := oldCapacity; i > 0; i-- {
		for old := oldMap[i-1]; old != nil; {
			e := old
			old = old.next
			index := uint(e.key) % uint(newCapacity)
			e.next = newMap[index]
			newMap[index] = e
		}
	}
}

func (this *IntSet) PutAll(values []int32) {
	if values == nil {
		return
	}
	ln := len(values)
	for i := 0; i < ln; i++ {
		this.Put(values[i])
	}
}

func (this *IntSet) Put(key int32) bool {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.put(key)
}

func (this *IntSet) put(value int32) bool {
	tab := this.table
	index := uint(value) % uint(len(tab))
	for e := tab[index]; e != nil; e = e.next {
		if e.key == value {
			return false
		}
	}
	if this.count >= this.threshold {
		this.rehash()
		tab = this.table
		index = uint(value) % uint(len(tab))
	}
	e := &IntSetry{key: value, next: tab[index]}
	tab[index] = e
	this.count++
	return true
}

func (this *IntSet) Remove(key int32) int32 {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.remove(key)
}

func (this *IntSet) remove(key int32) int32 {
	tab := this.table
	index := uint(key) % uint(len(tab))
	e := tab[index]
	var prev *IntSetry = nil
	for e != nil {
		if e.key == key {
			if prev != nil {
				prev.next = e.next
			} else {
				tab[index] = e.next
			}
			this.count--
			return key
		}
		prev = e
		e = e.next
	}
	return 0
}

func (this *IntSet) Clear() {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.clear()
}
func (this *IntSet) clear() {
	tab := this.table
	for index := len(tab) - 1; index >= 0; index-- {
		tab[index] = nil
	}
	this.count = 0
}

func (this *IntSet) ToString() string {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.toString()
}

func (this *IntSet) toString() string {
	max := this.Size() - 1

	var buf bytes.Buffer

	it := this.Values()
	buf.WriteString("{")

	for i := 0; i <= max; i++ {
		key := it.NextInt()
		buf.WriteString(strconv.Itoa(int(key)))
		
		if i < max {
			buf.WriteString(", ")
		}
	}

	buf.WriteString("}")
	return buf.String()
}

//public static IntEnumer emptyEnumer = new IntEnumer() {
//		public int nextInt() {
//			return 0;
//		}
//
//		public boolean hasMoreElements() {
//			return false;
//		}
//	};

type IntSetEnumer struct {
	table []*IntSetry
	index int
	entry *IntSetry
}

func NewIntSetEnumer(table []*IntSetry, index int) {
	p := new(IntSetEnumer)
	p.table = table
	p.index = len(table)
	p.entry = nil
}
func (this *IntSetEnumer) HasMoreElements() bool {
	for this.entry == nil && this.index > 0 {
		this.index--
		this.entry = this.table[this.index]
	}
	return this.entry != nil
}

func (this *IntSetEnumer) NextInt() int32 {
	for this.entry == nil && this.index > 0 {
		this.index--
		this.entry = this.table[this.index]
	}
	if this.entry != nil {
		e := this.entry
		this.entry = e.next
		return e.key
	}
	//panic("no more next")
	return 0
}

func IntSetMain() {
	s := NewIntSet()
	s.Put(111)
	s.Put(2222)
	s.Put(3333)
	//	s.Sort(c func(k1, k2 string) bool)
	fmt.Println(s)
	//	s.Sort(false)
	fmt.Println(s)
}
