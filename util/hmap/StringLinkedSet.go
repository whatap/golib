package hmap

import (
	"bytes"
	"fmt"
	"sort"
	"sync"

	"github.com/whatap/golib/util/stringutil"
)

// TODO SearchPathMap 에서 사용하기 위해 생성 (현재는 생성만)
// 테스트 필요.
type StringLinkedSet struct {
	table      []*StringLinkedSetry
	header     *StringLinkedSetry
	count      int
	threshold  int
	loadFactor float32
	lock       sync.Mutex
	max        int
}

func NewStringLinkedSet() *StringLinkedSet {

	initCapacity := DEFAULT_CAPACITY
	loadFactor := DEFAULT_LOAD_FACTOR

	this := new(StringLinkedSet)
	this.loadFactor = float32(loadFactor)
	this.table = make([]*StringLinkedSetry, initCapacity)
	this.header = &StringLinkedSetry{}
	this.header.link_next = this.header
	this.header.link_prev = this.header
	this.threshold = int(float64(initCapacity) * loadFactor)
	return this
}

func (this *StringLinkedSet) Size() int {
	return this.count
}

func (this *StringLinkedSet) GetArray() []string {
	this.lock.Lock()
	defer this.lock.Unlock()

	_keys := make([]string, this.Size())
	en := this.Keys()
	for i := 0; i < len(_keys); i++ {
		_keys[i] = en.NextString()
	}
	return _keys
}

type StringEnumerSetImpl struct {
	parent *StringLinkedSet
	entry  *StringLinkedSetry
	rtype  int
}

func (this *StringEnumerSetImpl) HasMoreElements() bool {
	return this.entry != nil && this.parent.header != this.entry
}

func (this *StringEnumerSetImpl) NextString() string {
	if this.HasMoreElements() {
		e := this.entry
		this.entry = e.link_next
		return e.Get()
	}
	return ""
}
func (this *StringLinkedSet) Keys() StringEnumer {
	return &StringEnumerSetImpl{parent: this, entry: this.header.link_next}
}

func (this *StringLinkedSet) Contains(key string) bool {
	this.lock.Lock()
	defer this.lock.Unlock()
	if key == "" {
		return false
	}

	tab := this.table
	index := this.hash(key) % uint(len(tab))
	for e := tab[index]; e != nil; e = e.hash_next {
		if e.key == key {
			return true
		}
	}
	return false

}

func (this *StringLinkedSet) GetFirst() string {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.header.link_next.key
}

func (this *StringLinkedSet) GetLast() string {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.header.link_prev.key
}
func (this *StringLinkedSet) hash(key string) uint {
	return uint(stringutil.HashCode(key))
}

func (this *StringLinkedSet) rehash() {
	oldCapacity := len(this.table)
	oldMap := this.table
	newCapacity := oldCapacity*2 + 1
	newMap := make([]*StringLinkedSetry, newCapacity)
	this.threshold = int(float32(newCapacity) * this.loadFactor)
	this.table = newMap
	for i := oldCapacity; i > 0; i-- {
		for old := oldMap[i-1]; old != nil; {
			e := old
			old = old.hash_next
			index := uint(this.hash(e.key) % uint(newCapacity))
			e.hash_next = newMap[index]
			newMap[index] = e
		}
	}
}

func (this *StringLinkedSet) SetMax(max int) *StringLinkedSet {
	this.max = max
	return this
}
func (this *StringLinkedSet) Put(key string) interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.put(key, PUT_LAST)
}
func (this *StringLinkedSet) PutLast(key string) interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.put(key, PUT_FORCE_LAST)
}
func (this *StringLinkedSet) PutFirst(key string) interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.put(key, PUT_FORCE_FIRST)
}
func (this *StringLinkedSet) put(key string, m PUT_MODE) interface{} {
	tab := this.table
	keyHash := this.hash(key)
	index := keyHash % uint(len(tab))
	for e := tab[index]; e != nil; e = e.hash_next {
		if e.key == key {
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
			return key
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
		index = keyHash % uint(len(tab))
	}
	e := &StringLinkedSetry{key: key, hash_next: tab[index]}
	tab[index] = e
	switch m {
	case PUT_FORCE_FIRST, PUT_FIRST:
		this.chain(this.header, this.header.link_next, e)
	case PUT_FORCE_LAST, PUT_LAST:
		this.chain(this.header.link_prev, this.header, e)
	}
	this.count++
	return nil
}

func (this *StringLinkedSet) Unipoint(key string) string {
	old := this.put(key, PUT_LAST)
	if old == nil {
		return key
	} else {
		return old.(string)
	}
}

func (this *StringLinkedSet) Remove(key string) interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.remove(key)
}
func (this *StringLinkedSet) RemoveFirst() interface{} {
	if this.IsEmpty() {
		return 0
	}
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.remove(this.header.link_next.key)
}

func (this *StringLinkedSet) RemoveLast() interface{} {
	if this.IsEmpty() {
		return 0
	}
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.remove(this.header.link_prev.key)
}

func (this *StringLinkedSet) remove(key string) interface{} {

	tab := this.table
	index := this.hash(key) % uint(len(tab))
	e := tab[index]
	var prev *StringLinkedSetry = nil
	for e != nil {
		if e.key == key {
			if prev != nil {
				prev.hash_next = e.hash_next
			} else {
				tab[index] = e.hash_next
			}
			this.count--
			//
			this.unchain(e)
			return key
		}
		prev = e
		e = e.hash_next
	}
	return nil
}

func (this *StringLinkedSet) IsEmpty() bool {
	return this.count == 0
}
func (this *StringLinkedSet) IsFull() bool {
	return this.max > 0 && this.max <= this.count
}

func (this *StringLinkedSet) Clear() {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.clear()
}
func (this *StringLinkedSet) clear() {
	tab := this.table
	for index := len(tab) - 1; index >= 0; index-- {
		tab[index] = nil
	}
	this.header.link_next = this.header
	this.header.link_prev = this.header
	this.count = 0
}

func (this *StringLinkedSet) chain(link_prev *StringLinkedSetry, link_next *StringLinkedSetry, e *StringLinkedSetry) {
	e.link_prev = link_prev
	e.link_next = link_next
	link_prev.link_next = e
	link_next.link_prev = e
}

func (this *StringLinkedSet) unchain(e *StringLinkedSetry) {
	e.link_prev.link_next = e.link_next
	e.link_next.link_prev = e.link_prev
	e.link_prev = nil
	e.link_next = nil
}

func (this *StringLinkedSet) ToString() string {
	this.lock.Lock()
	defer this.lock.Unlock()

	var buffer bytes.Buffer
	x := this.Keys()
	buffer.WriteString("{")
	for i := 0; x.HasMoreElements(); i++ {
		if i > 0 {
			buffer.WriteString(", ")
		}
		e := x.NextString()
		buffer.WriteString(e)
	}
	buffer.WriteString("}")
	return buffer.String()
}

type StringLinkedSortable struct {
	compare func(a, b string) bool
	data    []string
}

func (this StringLinkedSortable) Swap(i, j int) {
	this.data[i], this.data[j] = this.data[j], this.data[i]
}
func (this StringLinkedSortable) Len() int {
	return len(this.data)
}
func (this StringLinkedSortable) Less(i, j int) bool {
	return this.compare(this.data[i], this.data[j])
}

func (this *StringLinkedSet) Sort(c func(k1, k2 string) bool) {
	this.lock.Lock()
	defer this.lock.Unlock()

	sz := this.Size()
	list := make([]string, sz)
	en := this.Keys()
	for i := 0; i < sz; i++ {
		list[i] = en.NextString()
	}
	sort.Sort(StringLinkedSortable{compare: c, data: list})

	this.clear()
	for i := 0; i < sz; i++ {
		this.put(list[i], PUT_LAST)
	}
}

func StringLinkedSetMain() {
	s := NewStringLinkedSet()
	s.Put("aa")
	s.Put("bb")
	s.Put("00")
	//	s.Sort(c func(k1, k2 string) bool)
	fmt.Println(s)
	//	s.Sort(false)
	fmt.Println(s)
}
