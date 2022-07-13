package hmap

import (
	"fmt"
	//"log"
	"math"
	"sort"
	"sync"

	//GetKeySet
	"container/list"

	"github.com/whatap/golib/util/stringutil"
)

type IntKeyLinkedMap struct {
	table      []*IntKeyLinkedEntry
	header     *IntKeyLinkedEntry
	count      int
	threshold  int
	loadFactor float32
	max        int
	lock       sync.Mutex
}

func NewIntKeyLinkedMapDefault() *IntKeyLinkedMap {
	p := NewIntKeyLinkedMap(DEFAULT_CAPACITY, DEFAULT_LOAD_FACTOR)
	return p
}

func NewIntKeyLinkedMap(initCapacity int, loadFactor float32) *IntKeyLinkedMap {
	defer func() {
		if r := recover(); r != nil {
			// TODO 추후 hmap 에서 recover 는 없애고 panic 처리. 호출하는 쪽에서 recover 할 것
			//logutil.Println("WA822", r)
			//return NewIntKeyMap(DEFAULT_CAPACITY, DEFAULT_LOAD_FACTOR)
		}
	}()

	p := new(IntKeyLinkedMap)

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
	p.table = make([]*IntKeyLinkedEntry, initCapacity)
	p.header = NewIntKeyLinkedEntry(0, nil, nil)
	p.header.link_prev = p.header
	p.header.link_next = p.header.link_prev

	p.threshold = int(float32(initCapacity) * loadFactor)

	return p
}

func (this *IntKeyLinkedMap) Size() int {
	return this.count
}

func (this *IntKeyLinkedMap) KeyArray() []int32 {
	_keys := make([]int32, this.Size())

	en := this.Keys()
	for i := 0; i < len(_keys); i++ {
		_keys[i] = en.NextInt()
	}
	return _keys
}

func (this *IntKeyLinkedMap) GetKeySet() *IntLinkedSet {
	_keys := NewIntLinkedSet()
	en := this.Keys()
	for en.HasMoreElements() {
		_keys.Put(en.NextInt())
	}
	return _keys
}

func (this *IntKeyLinkedMap) Keys() IntEnumer {
	this.lock.Lock()
	defer this.lock.Unlock()

	//return &LongKeyEnumerImpl{parent: this, entry: this.header.link_next}
	return NewIntKeyLinkedEnumer(this, this.header.link_next, ELEMENT_TYPE_KEYS)
}

func (this *IntKeyLinkedMap) Values() Enumeration {
	this.lock.Lock()
	defer this.lock.Unlock()
	return NewIntKeyLinkedEnumer(this, this.header.link_next, ELEMENT_TYPE_VALUES)
}

func (this *IntKeyLinkedMap) ValueIterator() interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()
	return NewIntKeyLinkedEnumer(this, this.header.link_next, ELEMENT_TYPE_VALUES)
}

func (this *IntKeyLinkedMap) Entries() Enumeration {
	this.lock.Lock()
	defer this.lock.Unlock()
	return NewIntKeyLinkedEnumer(this, this.header.link_next, ELEMENT_TYPE_ENTRIES)
}

func (this *IntKeyLinkedMap) ContainsValue(value interface{}) bool {
	this.lock.Lock()
	defer this.lock.Unlock()
	if value == nil {
		panic("Value is Nil")
		//throw new NullPointerException();
	}
	tab := this.table

	for i := len(tab); i > 0; i-- {
		for e := tab[i]; e != nil; e = e.next {

			// TODO COmpareUtil
			//				if (CompareUtil.equals(e.value, value)) {
			//					return true;
			//				}
			if e.value == value {
				return true
			}
		}
	}
	return false
}

func (this *IntKeyLinkedMap) ContainsKey(key int32) bool {
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

func (this *IntKeyLinkedMap) Get(key int32) interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()
	tab := this.table
	index := this.hash(key) % uint(len(tab))
	for e := tab[index]; e != nil; e = e.next {
		//if (CompareUtil.equals(e.key, key)) {
		if e.key == key {
			return e.value
		}
	}
	return nil
}

func (this *IntKeyLinkedMap) GetLRU(key int32) interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()
	tab := this.table
	index := this.hash(key) % uint(len(tab))

	for e := tab[index]; e != nil; e = e.next {
		if e.key == key {
			old := e.value
			if this.header.link_prev != e {
				this.unchain(e)
				this.chain(this.header.link_prev, this.header, e)
			}
			return old
		}
	}
	return nil
}

func (this *IntKeyLinkedMap) GetFirstKey() int32 {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.header.link_next.key
}

func (this *IntKeyLinkedMap) GetLastKey() int32 {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.header.link_prev.key
}

func (this *IntKeyLinkedMap) GetFirstValue() interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.header.link_next.value
}

func (this *IntKeyLinkedMap) GetLastValue() interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.header.link_prev.value
}

func (this *IntKeyLinkedMap) overflowed(key int32, value interface{}) {
}

func (this *IntKeyLinkedMap) hash(key int32) uint {
	//return uint(key & math.MaxInt32)
	return uint(key & math.MaxInt32)
}

func (this *IntKeyLinkedMap) rehash() {
	oldCapacity := len(this.table)
	oldMap := this.table
	newCapacity := oldCapacity*2 + 1
	newMap := make([]*IntKeyLinkedEntry, newCapacity)
	this.threshold = int(float32(newCapacity) * this.loadFactor)
	this.table = newMap
	for i := oldCapacity; i > 0; i-- {
		old := oldMap[i-1]
		for old != nil {
			e := old
			old = old.next
			index := e.keyHash % uint(newCapacity)
			e.next = newMap[index]
			newMap[index] = e
		}
	}
}

func (this *IntKeyLinkedMap) SetMax(max int) *IntKeyLinkedMap {
	this.max = max
	return this
}

func (this *IntKeyLinkedMap) IsFull() bool {
	return this.max > 0 && this.max <= this.count
}

func (this *IntKeyLinkedMap) Put(key int32, value interface{}) interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.put(key, value, PUT_LAST)
}

func (this *IntKeyLinkedMap) PutLast(key int32, value interface{}) interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.put(key, value, PUT_FORCE_LAST)
}

func (this *IntKeyLinkedMap) PutFirst(key int32, value interface{}) interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.put(key, value, PUT_FORCE_FIRST)
}

func (this *IntKeyLinkedMap) put(key int32, value interface{}, m PUT_MODE) interface{} {

	tab := this.table
	keyHash := this.hash(key)
	index := keyHash % uint(len(tab))

	for e := tab[index]; e != nil; e = e.next {
		//if (CompareUtil.equals(e.key, key)) {
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

			// DEBUG IntKeyLinkedMap
			//logutil.Println("IntKeyLinkedMap Exists return array, key, value", e.key, key, old)
			return old
		}
	}

	if this.max > 0 {
		// DEBUG IntKeyLinkedMap
		//log.Println("IntKeyLinkedMap Exists Max")

		switch m {
		case PUT_FORCE_FIRST, PUT_FIRST:
			for this.count >= this.max {
				// DEBUG IntKeyLinkedMap
				//log.Println("IntKeyLinkedMap PUT_FIRST  Over Max count, max", this.count, this.max)

				// removeLast();
				k := this.header.link_prev.key
				v := this.remove(k)
				this.overflowed(k, v)
			}

		case PUT_FORCE_LAST, PUT_LAST:
			for this.count >= this.max {
				// DEBUG IntKeyLinkedMap
				//log.Println("IntKeyLinkedMap PUT_LAST Over Max count, max", this.count, this.max)

				// removeFirst();
				k := this.header.link_next.key
				v := this.remove(k)
				this.overflowed(k, v)
			}

		}
	}
	if this.count >= this.threshold {
		// DEBUG IntKeyLinkedMap
		//logutil.Println("IntKeyLinkedMap rehash count, threshold", this.count, this.threshold)

		this.rehash()
		tab = this.table
		index = keyHash % uint(len(tab))
	}

	e := &IntKeyLinkedEntry{key: key, keyHash: keyHash, value: value, next: tab[index]}
	tab[index] = e
	switch m {
	case PUT_FORCE_FIRST, PUT_FIRST:
		this.chain(this.header, this.header.link_next, e)
	case PUT_FORCE_LAST, PUT_LAST:
		this.chain(this.header.link_prev, this.header, e)
	}
	this.count++

	// DEBUG IntKeyLinkedMap
	//logutil.Println("IntKeyLinkedMap Add & End", this.count)

	return nil
}
func (this *IntKeyLinkedMap) remove(key int32) interface{} {
	// DEBUG IntKeyLinkedMap
	//log.Println("IntKeyLinkedMap Remove key", key)

	tab := this.table
	index := this.hash(key) % uint(len(tab))
	var prev *IntKeyLinkedEntry
	prev = nil

	for e := tab[index]; e != nil; e = e.next {
		if e.key == key {
			if prev != nil {
				prev.next = e.next
			} else {
				tab[index] = e.next
			}
			this.count--
			oldValue := e.value
			e.value = nil
			this.unchain(e)
			return oldValue
		}

		prev = e
	}

	return nil
}
func (this *IntKeyLinkedMap) Remove(key int32) interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.remove(key)
}

func (this *IntKeyLinkedMap) RemoveFirst() interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()

	if this.IsEmpty() {
		return nil
	}
	return this.remove(this.header.link_next.key)
}

func (this *IntKeyLinkedMap) RemoveLast() interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()

	if this.IsEmpty() {
		return nil
	}

	return this.remove(this.header.link_prev.key)
}

func (this *IntKeyLinkedMap) IsEmpty() bool {
	return this.Size() == 0
}

func (this *IntKeyLinkedMap) Clear() {
	// DEBUG IntKeyLinkedMap
	//logutil.Println("IntKeyLinkedMap Clear")

	this.lock.Lock()
	defer this.lock.Unlock()
	this.clear()
}

func (this *IntKeyLinkedMap) clear() {

	tab := this.table
	//index := this.hash(key) % uint(len(tab))
	for index := len(tab) - 1; index >= 0; index-- {
		tab[index] = nil
	}

	this.header.link_prev = this.header
	this.header.link_next = this.header.link_prev
	this.count = 0
}

func (this *IntKeyLinkedMap) ToString() string {
	buf := stringutil.NewStringBuffer()
	it := this.Entries()
	buf.Append("{")
	for i := 0; it.HasMoreElements(); i++ {
		e := it.NextElement().(*IntKeyLinkedEntry)
		if i > 0 {
			buf.Append(", ")
		}
		buf.Append(fmt.Sprintf("%d=%v", e.GetKey(), e.GetValue()))
	}
	buf.Append("}")
	return buf.ToString()
}

func (this *IntKeyLinkedMap) ToFormatString() string {
	buf := stringutil.NewStringBuffer()
	it := this.Entries()
	buf.Append("{")
	for i := 0; it.HasMoreElements(); i++ {
		e := it.NextElement().(*IntKeyLinkedEntry)
		if i > 0 {
			buf.Append(", ")
		}
		buf.Append(fmt.Sprintf("%d=%v", e.GetKey(), e.GetValue())).Append("\n")

	}
	buf.Append("}")
	return buf.ToString()
}

func (this *IntKeyLinkedMap) chain(link_prev *IntKeyLinkedEntry, link_next *IntKeyLinkedEntry, e *IntKeyLinkedEntry) {
	e.link_prev = link_prev
	e.link_next = link_next
	link_prev.link_next = e
	link_next.link_prev = e
}

func (this *IntKeyLinkedMap) unchain(e *IntKeyLinkedEntry) {
	e.link_prev.link_next = e.link_next
	e.link_next.link_prev = e.link_prev
	e.link_prev = nil
	e.link_next = nil
}

// java hashset 을 list  변환해서 반환
func (this *IntKeyLinkedMap) ToKeySet() *list.List {

	keyList := list.New()
	en := this.Keys()

	for en.HasMoreElements() {

		keyList.PushFront(en.NextInt())

	}
	return keyList
}

func (this *IntKeyLinkedMap) Sort(c func(k1, k2 int32) bool) {
	this.lock.Lock()
	defer this.lock.Unlock()

	sz := this.Size()
	entryList := make([]*IntKeyLinkedEntry, sz)
	en := this.Entries()
	for i := 0; i < sz; i++ {
		entryList[i] = en.NextElement().(*IntKeyLinkedEntry)
	}
	sort.Sort(IntKeySortable{compare: c, data: entryList})
	this.clear()
	for i := 0; i < sz; i++ {
		this.put(entryList[i].GetKey(), entryList[i].GetValue(), PUT_LAST)
	}

}

type IntKeyLinkedEnumer struct {
	parent *IntKeyLinkedMap
	entry  *IntKeyLinkedEntry
	Type   int
}

func NewIntKeyLinkedEnumer(parent *IntKeyLinkedMap, entry *IntKeyLinkedEntry, Type int) *IntKeyLinkedEnumer {
	p := new(IntKeyLinkedEnumer)
	p.parent = parent
	p.entry = entry
	p.Type = Type

	return p
}
func (this *IntKeyLinkedEnumer) HasNext() bool {
	return this.HasMoreElements()
}

func (this *IntKeyLinkedEnumer) HasMoreElements() bool {
	return this.parent.header != this.entry && this.entry != nil
}

func (this *IntKeyLinkedEnumer) Next() interface{} {
	return this.NextElement()
}

func (this *IntKeyLinkedEnumer) NextElement() interface{} {
	if this.HasMoreElements() {
		e := this.entry
		this.entry = e.link_next

		switch this.Type {
		case ELEMENT_TYPE_KEYS:
			//return (V) new Long(e.key);
			return e.key
		case ELEMENT_TYPE_VALUES:
			return e.value
		default:
			return e
		}
	}
	panic("no more next")
	//throw new NoSuchElementException("no more next");
}

func (this *IntKeyLinkedEnumer) NextInt() int32 {
	if this.HasMoreElements() {
		e := this.entry
		this.entry = e.link_next
		return e.key
	}
	panic("no more next")
	//throw new NoSuchElementException("no more next");
}

func (this *IntKeyLinkedEnumer) Remove() {

}

// implements sort.Interface
type IntKeySortable struct {
	// func(a, b, int32) bool
	compare func(a, b int32) bool
	// []*IntKeyLinkedENtry
	data []*IntKeyLinkedEntry
}

func (this IntKeySortable) Len() int {
	return len(this.data)
}
func (this IntKeySortable) Less(i, j int) bool {
	return this.compare(this.data[i].GetKey(), this.data[j].GetKey())
}

func (this IntKeySortable) Swap(i, j int) {
	this.data[i], this.data[j] = this.data[j], this.data[i]
}

//func main() {
//		m = NewIntKeyLinkedMapDefault.setMax(5);
//		// System.out.println(m.getFirstValue());
//		// System.out.println(m.getLastKey());
//		for i := 0; i < 10; i++ {
//			m.putFirst(i, i);
//		}
////		for (int i = 1; i < 10; i+=2) {
////			m.put(i, i);
////		}
////		m.sort(new Comparator<IntKeyLinkedMap.IntKeyLinkedEntry<Integer>>() {
////			public int compare(IntKeyLinkedEntry<Integer> o1, IntKeyLinkedEntry<Integer> o2) {
////				return o1.key - o2.key;
////			}
////		});
//		fmt.Println(m);
//			// System.out.println("==================================");
//		// for(int i=0; i <10; i++){
//		// m.putLast(i, i);
//		// System.out.println(m);
//		// }
//		// System.out.println("==================================");
//		// for(int i=0; i <10; i++){
//		// m.putFirst(i, i);
//		// System.out.println(m);
//		// }
////		IntEnumer e = m.keys();
////		System.out.println("==================================");
////		for (int i = 0; i < 10; i++) {
////			m.removeFirst();
////			System.out.println(m);
////		}
////		System.out.println("==================================");
////		while (e.hasMoreElements()) {
////			System.out.println(e.nextInt());
////		}
//	}

//	private static void print(Object e) {
//		System.out.println(e);
//	}
