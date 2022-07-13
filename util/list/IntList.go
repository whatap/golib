package list

import (
	//"bytes"
	"fmt"
	"math"
	"sort"
	"strconv"
	"sync"

	"github.com/whatap/golib/io"
	"github.com/whatap/golib/lang/value"
	"github.com/whatap/golib/util/compare"
)

type IntList struct {
	size  int
	table []int
	lock  sync.Mutex
}

func NewIntListDefault() *IntList {
	o := new(IntList)
	o.table = make([]int, 0)
	return o
}

func NewIntList(initialCapa int) *IntList {
	o := new(IntList)
	o.table = make([]int, initialCapa)
	return o
}

func (this *IntList) GetType() byte {
	return ANYLIST_INT
}

func (this *IntList) ensure(minCapacity int) {
	if minCapacity > len(this.table) {
		if this.table == nil {
			minCapacity = int(math.Min(float64(ANYLIST_DEFAULT_CAPACITY), float64(minCapacity)))
		}
		oldSize := len(this.table)
		newSize := oldSize + (oldSize >> 1)
		if newSize < minCapacity {
			newSize = minCapacity
		}
		if newSize > ANYLIST_MAX_SIZE {
			panic("too big size")
		}

		newTable := make([]int, newSize)
		copy(newTable, this.table)
		this.table = newTable
	}
}

func (this *IntList) get(i int) int {
	if i >= this.size {
		panic(fmt.Sprintln("Index out of bounds ", "Index: ", i, ", Size: ", this.size))
	}
	return this.table[i]
}

func (this *IntList) set(i, v int) int {
	if i >= this.size {
		panic(fmt.Sprintln("Index out of bounds Index: ", i, ", Size: ", this.size))
	}

	ov := this.table[i]
	this.table[i] = v
	return ov
}

func (this *IntList) remove(i int) int {
	if i >= this.size {
		panic(fmt.Sprintln("Index out of bounds Index: ", i, ", Size: ", this.size))
	}
	ov := this.table[i]

	movNum := this.size - i - 1
	if movNum > 0 {
		copy(this.table[i:i+movNum], this.table[i+1:i+movNum])
		//System.arraycopy(table, i + 1, table, i, movNum);
	}
	//table[--size] = 0;
	last := this.size - 1
	this.table[last] = 0
	return ov
}

func (this *IntList) add(e int) {
	this.ensure(this.size + 1)
	//table[size++] = e;
	this.table[this.size] = e
	this.size++
}

func (this *IntList) AddAll(other *IntList) {
	this.ensure(this.size + other.size)
	for i := 0; i < other.size; i++ {
		//this.table[this.size++] = other.table[i];
		this.table[this.size] = other.table[i]
		this.size++
	}
}

func (this *IntList) AddAllArray(other []int) {
	this.ensure(this.size + len(other))
	n := len(other)
	for i := 0; i < n; i++ {
		//this.table[this.size++] = other[i];
		this.table[this.size] = other[i]
		this.size++
	}
}

func (this *IntList) ToArray() []int {
	newArray := make([]int, this.size)
	copy(newArray, this.table)
	return newArray
}

func (this *IntList) ToString() string {
	return fmt.Sprintf("%v", this.table)
}

func (this *IntList) SetInt(i, v int) {
	this.set(i, v)
}
func (this *IntList) SetFloat(i int, v float32) {
	this.set(i, int(v))
}
func (this *IntList) SetLong(i int, v int64) {
	this.set(i, int(v))
}
func (this *IntList) SetDouble(i int, v float64) {
	this.set(i, int(v))
}
func (this *IntList) SetString(i int, v string) {
	n, err := strconv.Atoi(v)
	if err == nil {
		this.set(i, n)
	} else {
		panic(" Parse error string(" + v + ") to int")
	}

}

func (this *IntList) AddInt(v int) {
	this.add(v)
}
func (this *IntList) AddFloat(v float32) {
	this.add(int(v))
}
func (this *IntList) AddLong(v int64) {
	this.add(int(v))
}
func (this *IntList) AddDouble(v float64) {
	this.add(int(v))
}
func (this *IntList) AddString(v string) {

	n, err := strconv.Atoi(v)
	if err == nil {
		this.add(n)
	} else {
		panic(" Parse error string(" + v + ") to int")
	}
}

func (this *IntList) GetInt(i int) int {
	return this.get(i)
}
func (this *IntList) GetFloat(i int) float32 {
	return float32(this.get(i))
}
func (this *IntList) GetDouble(i int) float64 {
	return float64(this.get(i))
}
func (this *IntList) GetLong(i int) int64 {
	return int64(this.get(i))
}
func (this *IntList) GetString(i int) string {
	return strconv.Itoa(this.get(i))
}
func (this *IntList) GetObject(i int) interface{} {
	return nil
}
func (this *IntList) GetValue(i int) value.Value {

	return value.NewDecimalValue(int64(this.get(i)))
}

func (this *IntList) Write(out *io.DataOutputX) {
	out.WriteInt3(int32(this.Size()))
	for i := 0; i < this.size; i++ {
		out.WriteDecimal(int64(this.get(i)))
	}
}
func (this *IntList) Read(in *io.DataInputX) {
	count := int(in.ReadInt3())
	for i := 0; i < count; i++ {
		this.add(int(in.ReadDecimal()))
	}
}

func (this *IntList) Size() int {
	return this.size
}

func (this *IntList) Sorting(asc bool) []int {
	this.lock.Lock()
	defer this.lock.Unlock()

	table := make([]*IntListKeyVal, this.size)
	for i := 0; i < this.size; i++ {
		table[i] = &IntListKeyVal{i, this.get(i)}
	}

	c := func(o1, o2 *IntListKeyVal) bool {
		if asc {
			if compare.CompareToInt(o1.value, o2.value) > 0 {
				return false
			} else {
				return true
			}
		} else {
			if compare.CompareToInt(o2.value, o1.value) > 0 {
				return false
			} else {
				return true
			}
		}
	}

	sort.Sort(IntListSortable{compare: c, data: table})

	out := make([]int, this.size)
	for i := 0; i < this.size; i++ {
		out[i] = table[i].key
	}
	return out
}
func (this *IntList) SortingAnyList(asc bool, child AnyList, childAsc bool) []int {
	this.lock.Lock()
	defer this.lock.Unlock()

	table := make([]*IntListKeyVal, this.size)
	for i := 0; i < this.size; i++ {
		table[i] = &IntListKeyVal{i, this.get(i)}
	}

	c := func(o1, o2 *IntListKeyVal) bool {
		var rt int
		if asc {
			rt = compare.CompareToInt(o1.value, o2.value)
		} else {
			rt = compare.CompareToInt(o2.value, o1.value)
		}

		if rt != 0 {
			if rt > 0 {
				return false
			} else {
				return true
			}
		} else {
			rt = CompareChild(child, childAsc, o1.key, o2.key)
			if rt > 0 {
				return false
			} else {
				return true
			}
		}
	}

	sort.Sort(IntListSortable{compare: c, data: table})

	out := make([]int, this.size)
	for i := 0; i < this.size; i++ {
		out[i] = table[i].key
	}
	return out
}
func (this *IntList) Filtering(index []int) AnyList {
	out := NewIntList(this.size)
	sz := len(index)
	for i := 0; i < sz; i++ {
		out.add(this.get(index[i]))
	}
	return out
}

type IntListSortable struct {
	compare func(a, b *IntListKeyVal) bool
	data    []*IntListKeyVal
}

func (this IntListSortable) Len() int {
	return len(this.data)
}

func (this IntListSortable) Less(i, j int) bool {
	return this.compare(this.data[i], this.data[j])
}

func (this IntListSortable) Swap(i, j int) {
	this.data[i], this.data[j] = this.data[j], this.data[i]
}

type IntListKeyVal struct {
	key   int
	value int
}

func IntListMain() {
	list := NewIntList(30)

	for i := 0; i < 30; i++ {
		list.AddInt(i)
	}

	fmt.Println("\nlist")
	for i := 0; i < list.Size(); i++ {
		fmt.Printf("%d ", list.GetInt(i))
	}

	fmt.Println("\nlist Sorting ASC")
	ascList := list.Sorting(true)

	for i := 0; i < len(ascList); i++ {
		fmt.Printf("%d ", ascList[i])
	}

	fmt.Println("\nlist Sorting DESC")
	descList := list.Sorting(false)

	for i := 0; i < len(descList); i++ {
		fmt.Printf("%d ", descList[i])
	}

	fmt.Println("\nlist Filtering")
	var listFilter = make([]int, 30)

	var fList = [30]int{3, 2, 1, 0, 4, 5, 6, 7, 8, 9, 10, 20, 19, 18, 17, 16, 15, 14, 13, 12, 11, 29, 28, 27, 26, 25, 24, 23, 22, 21}

	for i := 0; i < len(fList); i++ {
		listFilter[i] = fList[i]
	}
	aList := list.Filtering(listFilter)

	for i := 0; i < aList.Size(); i++ {
		fmt.Printf("%d ", aList.GetInt(i))
	}

}
