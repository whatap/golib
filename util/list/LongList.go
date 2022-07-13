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

type LongList struct {
	size  int
	table []int64
	lock  sync.Mutex
}

func NewLongListDefault() *LongList {
	o := new(LongList)
	o.table = make([]int64, 0)
	return o
}

func NewLongList(initialCapa int) *LongList {
	o := new(LongList)
	o.table = make([]int64, initialCapa)
	return o
}

func (this *LongList) GetType() byte {
	return ANYLIST_LONG
}

func (this *LongList) ensure(minCapacity int) {
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

		newTable := make([]int64, newSize)
		copy(newTable, this.table)
		this.table = newTable
	}
}

func (this *LongList) get(i int) int64 {
	if i >= this.size {
		panic(fmt.Sprintln("Index out of bounds ", "Index: ", i, ", Size: ", this.size))
	}
	return this.table[i]
}

func (this *LongList) set(i int, v int64) int64 {
	if i >= this.size {
		panic(fmt.Sprintln("Index out of bounds Index: ", i, ", Size: ", this.size))
	}

	ov := this.table[i]
	this.table[i] = v
	return ov
}

func (this *LongList) remove(i int) int64 {
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

func (this *LongList) add(e int64) {
	this.ensure(this.size + 1)
	//table[size++] = e;
	this.table[this.size] = e
	this.size++
}

func (this *LongList) AddAll(other *LongList) {
	this.ensure(this.size + other.size)
	for i := 0; i < other.size; i++ {
		//this.table[this.size++] = other.table[i];
		this.table[this.size] = other.table[i]
		this.size++
	}
}

func (this *LongList) AddAllArray(other []int64) {
	this.ensure(this.size + len(other))
	n := len(other)
	for i := 0; i < n; i++ {
		//this.table[this.size++] = other[i];
		this.table[this.size] = other[i]
		this.size++
	}
}

func (this *LongList) ToArray() []int64 {
	newArray := make([]int64, this.size)
	copy(newArray, this.table)
	return newArray
}

func (this *LongList) ToString() string {
	return fmt.Sprintf("%v", this.table)
}

func (this *LongList) Sort() {
	if this.size > 0 {
		//Arrays.sort(this.table, 0, size-1)
		//sort.Ints(this.table)
	}
}

func (this *LongList) SetInt(i, v int) {
	this.set(i, int64(v))
}
func (this *LongList) SetFloat(i int, v float32) {
	this.set(i, int64(v))
}
func (this *LongList) SetLong(i int, v int64) {
	this.set(i, v)
}
func (this *LongList) SetDouble(i int, v float64) {
	this.set(i, int64(v))
}
func (this *LongList) SetString(i int, v string) {
	n, err := strconv.ParseInt(v, 10, 64)
	if err == nil {
		this.set(i, n)
	} else {
		panic(" Parse error string(" + v + ") to int64")
	}

}

func (this *LongList) AddInt(v int) {
	this.add(int64(v))
}
func (this *LongList) AddFloat(v float32) {
	this.add(int64(v))
}
func (this *LongList) AddLong(v int64) {
	this.add(v)
}
func (this *LongList) AddDouble(v float64) {
	this.add(int64(v))
}
func (this *LongList) AddString(v string) {

	n, err := strconv.ParseInt(v, 10, 64)
	if err == nil {
		this.add(n)
	} else {
		panic(" Parse error string(" + v + ") to int64")
	}
}

func (this *LongList) GetInt(i int) int {
	return int(this.get(i))
}
func (this *LongList) GetFloat(i int) float32 {
	return float32(this.get(i))
}
func (this *LongList) GetDouble(i int) float64 {
	return float64(this.get(i))
}
func (this *LongList) GetLong(i int) int64 {
	return this.get(i)
}
func (this *LongList) GetString(i int) string {
	return strconv.FormatInt(this.get(i), 10)
	//return fmt.Sprintf("%d", this.get(i))
}
func (this *LongList) GetObject(i int) interface{} {
	return nil
}
func (this *LongList) GetValue(i int) value.Value {

	return value.NewDecimalValue(this.get(i))
}

func (this *LongList) Write(out *io.DataOutputX) {
	out.WriteInt3(int32(this.Size()))
	for i := 0; i < this.size; i++ {
		out.WriteDecimal(this.get(i))
	}
}
func (this *LongList) Read(in *io.DataInputX) {
	count := int(in.ReadInt3())
	for i := 0; i < count; i++ {
		this.add(in.ReadDecimal())
	}
}

func (this *LongList) Size() int {
	return this.size
}

func (this *LongList) Sorting(asc bool) []int {
	this.lock.Lock()
	defer this.lock.Unlock()

	table := make([]*LongListKeyVal, this.size)
	for i := 0; i < this.size; i++ {
		table[i] = &LongListKeyVal{i, this.get(i)}
	}

	c := func(o1, o2 *LongListKeyVal) bool {
		if asc {
			if compare.CompareToLong(o1.value, o2.value) > 0 {
				return false
			} else {
				return true
			}
		} else {
			if compare.CompareToLong(o2.value, o1.value) > 0 {
				return false
			} else {
				return true
			}
		}
	}

	sort.Sort(LongListSortable{compare: c, data: table})

	out := make([]int, this.size)
	for i := 0; i < this.size; i++ {
		out[i] = table[i].key
	}
	return out
}
func (this *LongList) SortingAnyList(asc bool, child AnyList, childAsc bool) []int {
	this.lock.Lock()
	defer this.lock.Unlock()

	table := make([]*LongListKeyVal, this.size)
	for i := 0; i < this.size; i++ {
		table[i] = &LongListKeyVal{i, this.get(i)}
	}

	c := func(o1, o2 *LongListKeyVal) bool {
		var rt int
		if asc {
			rt = compare.CompareToLong(o1.value, o2.value)
		} else {
			rt = compare.CompareToLong(o2.value, o1.value)
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

	sort.Sort(LongListSortable{compare: c, data: table})

	out := make([]int, this.size)
	for i := 0; i < this.size; i++ {
		out[i] = table[i].key
	}
	return out
}
func (this *LongList) Filtering(index []int) AnyList {
	out := NewLongList(this.size)
	sz := len(index)
	for i := 0; i < sz; i++ {
		out.add(this.get(index[i]))
	}
	return out
}

type LongListSortable struct {
	compare func(a, b *LongListKeyVal) bool
	data    []*LongListKeyVal
}

func (this LongListSortable) Len() int {
	return len(this.data)
}

func (this LongListSortable) Less(i, j int) bool {
	return this.compare(this.data[i], this.data[j])
}

func (this LongListSortable) Swap(i, j int) {
	this.data[i], this.data[j] = this.data[j], this.data[i]
}

type LongListKeyVal struct {
	key   int
	value int64
}

func LongListMain() {
	list := NewLongList(30)

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
