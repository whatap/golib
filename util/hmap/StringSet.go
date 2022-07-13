package hmap

import (
	"fmt"
	"sync"

	"github.com/whatap/golib/util/hash"
)

// TODO SearchPathMap 에서 사용하기 위해 생성 (현재는 생성만)
// 테스트 필요.
type StringSet struct {
	table      []*StringSetry
	count      int
	threshold  int
	loadFactor float32
	lock       sync.Mutex
	max        int
}

func NewStringSet() *StringSet {

	initCapacity := DEFAULT_CAPACITY
	loadFactor := DEFAULT_LOAD_FACTOR

	this := new(StringSet)
	this.loadFactor = float32(loadFactor)
	this.table = make([]*StringSetry, initCapacity)
	this.threshold = int(float64(initCapacity) * loadFactor)
	return this
}

func NewStringSetArray(arr []string) *StringSet {

	initCapacity := DEFAULT_CAPACITY
	loadFactor := DEFAULT_LOAD_FACTOR

	this := new(StringSet)
	this.loadFactor = float32(loadFactor)
	this.table = make([]*StringSetry, initCapacity)
	this.threshold = int(float64(initCapacity) * loadFactor)
	return this
}

func (this *StringSet) Size() int {
	return this.count
}

func (this *StringSet) Keys() StringEnumer {
	return &StringSetEnumer{table: this.table, index: len(this.table), entry: nil}
}

func (this *StringSet) HasKey(key string) bool {
	return this.Contains(key)
}

func (this *StringSet) Contains(key string) bool {
	this.lock.Lock()
	defer this.lock.Unlock()
	if key == "" {
		return false
	}

	tab := this.table
	index := this.hash(key) % uint(len(tab))
	for e := tab[index]; e != nil; e = e.next {
		if e.key == key {
			return true
		}
	}
	return false
}

func (this *StringSet) rehash() {
	oldCapacity := len(this.table)
	oldMap := this.table
	newCapacity := oldCapacity*2 + 1
	newMap := make([]*StringSetry, newCapacity)
	this.threshold = int(float32(newCapacity) * this.loadFactor)
	this.table = newMap
	for i := oldCapacity; i > 0; i-- {
		for old := oldMap[i-1]; old != nil; {
			e := old
			old = old.next
			index := uint(this.hash(e.key) % uint(newCapacity))
			e.next = newMap[index]
			newMap[index] = e
		}
	}
}

func (this *StringSet) Put(key string) string {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.unipoint(key)
}

func (this *StringSet) Unipoint(key string) string {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.unipoint(key)
}

func (this *StringSet) unipoint(key string) string {
	if key == "" {
		return ""
	}
	tab := this.table
	hash := this.hash(key)
	index := hash % uint(len(tab))
	for e := tab[index]; e != nil; e = e.next {
		if e.key == key {
			return e.key
		}
	}
	if this.count >= this.threshold {
		this.rehash()
		tab = this.table
		index = hash % uint(len(tab))
	}
	e := &StringSetry{hash: hash, key: key, next: tab[index]}
	tab[index] = e
	this.count++
	return key
}

func (this *StringSet) Remove(key string) bool {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.remove(key)
}

func (this *StringSet) remove(key string) bool {
	if key == "" {
		return false
	}
	tab := this.table
	index := this.hash(key) % uint(len(tab))
	e := tab[index]
	var prev *StringSetry = nil
	for e != nil {
		if e.key == key {
			if prev != nil {
				prev.next = e.next
			} else {
				tab[index] = e.next
			}
			this.count--
			return true
		}
		prev = e
		e = e.next
	}
	return false
}

func (this *StringSet) Clear() {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.clear()
}
func (this *StringSet) clear() {
	tab := this.table
	for index := len(tab) - 1; index >= 0; index-- {
		tab[index] = nil
	}
	this.count = 0
}

func (this *StringSet) hash(key string) uint {
	return uint(hash.HashStr(key))
}

type StringSetEnumer struct {
	table []*StringSetry
	index int
	entry *StringSetry
}

func NewStringSetEnumer(table []*StringSetry, index int) {
	p := new(StringSetEnumer)
	p.table = table
	p.index = len(table)
	p.entry = nil
}
func (this *StringSetEnumer) HasMoreElements() bool {
	for this.entry == nil && this.index > 0 {
		this.index--
		this.entry = this.table[this.index]
	}
	return this.entry != nil
}

func (this *StringSetEnumer) NextString() string {
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
	return ""
}

func StringSetMain() {
	s := NewStringSet()
	s.Put("aa")
	s.Put("bb")
	s.Put("00")
	//	s.Sort(c func(k1, k2 string) bool)
	fmt.Println(s)
	//	s.Sort(false)
	fmt.Println(s)
}
