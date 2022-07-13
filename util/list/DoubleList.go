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

type DoubleList struct {
	size  int
	table []float64
	lock  sync.Mutex
}

func NewDoubleListDefault() *DoubleList {
	o := new(DoubleList)
	o.table = make([]float64, 0)
	return o
}

func NewDoubleList(initialCapa int) *DoubleList {
	o := new(DoubleList)
	o.table = make([]float64, initialCapa)
	return o
}
func (this *DoubleList) GetType() byte {
	return ANYLIST_DOUBLE
}

func (this *DoubleList) ensure(minCapacity int) {
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

		newTable := make([]float64, newSize)
		copy(newTable, this.table)
		this.table = newTable
	}
}

func (this *DoubleList) get(i int) float64 {
	if i >= this.size {
		panic(fmt.Sprintln("Index out of bounds ", "Index: ", i, ", Size: ", this.size))
	}
	return this.table[i]
}

func (this *DoubleList) set(i int, v float64) float64 {
	if i >= this.size {
		panic(fmt.Sprintln("Index out of bounds Index: ", i, ", Size: ", this.size))
	}

	ov := this.table[i]
	this.table[i] = v
	return ov
}

func (this *DoubleList) remove(i int) float64 {
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

func (this *DoubleList) add(e float64) {
	this.ensure(this.size + 1)
	//table[size++] = e;
	this.table[this.size] = e
	this.size++
}

func (this *DoubleList) AddAll(other *DoubleList) {
	this.ensure(this.size + other.size)
	for i := 0; i < other.size; i++ {
		//this.table[this.size++] = other.table[i];
		this.table[this.size] = other.table[i]
		this.size++
	}
}

func (this *DoubleList) AddAllArray(other []float64) {
	this.ensure(this.size + len(other))
	n := len(other)
	for i := 0; i < n; i++ {
		//this.table[this.size++] = other[i];
		this.table[this.size] = other[i]
		this.size++
	}
}

func (this *DoubleList) ToArray() []float64 {
	newArray := make([]float64, this.size)
	copy(newArray, this.table)
	return newArray
}

func (this *DoubleList) ToString() string {
	return fmt.Sprintf("%v", this.table)
}

func (this *DoubleList) Sort() {
	if this.size > 0 {
		//Arrays.sort(this.table, 0, size-1)
		//sort.Ints(this.table)
	}
}

func (this *DoubleList) SetInt(i, v int) {
	this.set(i, float64(v))
}
func (this *DoubleList) SetFloat(i int, v float32) {
	this.set(i, float64(v))
}
func (this *DoubleList) SetLong(i int, v int64) {
	this.set(i, float64(v))
}
func (this *DoubleList) SetDouble(i int, v float64) {
	this.set(i, v)
}
func (this *DoubleList) SetString(i int, v string) {
	n, err := strconv.ParseFloat(v, 32)
	if err == nil {
		this.set(i, float64(n))
	} else {
		panic(" Parse error string(" + v + ") to float64")
	}

}

func (this *DoubleList) AddInt(v int) {
	this.add(float64(v))
}
func (this *DoubleList) AddFloat(v float32) {
	this.add(float64(v))
}
func (this *DoubleList) AddLong(v int64) {
	this.add(float64(v))
}
func (this *DoubleList) AddDouble(v float64) {
	this.add(v)
}
func (this *DoubleList) AddString(v string) {

	n, err := strconv.ParseFloat(v, 64)
	if err == nil {
		this.add(n)
	} else {
		panic(" Parse error string(" + v + ") to float64")
	}
}

func (this *DoubleList) GetInt(i int) int {
	return int(this.get(i))
}
func (this *DoubleList) GetFloat(i int) float32 {
	return float32(this.get(i))
}
func (this *DoubleList) GetDouble(i int) float64 {
	return this.get(i)
}
func (this *DoubleList) GetLong(i int) int64 {
	return int64(this.get(i))
}
func (this *DoubleList) GetString(i int) string {
	return strconv.FormatFloat(this.get(i), 'f', 6, 64)
	//return fmt.Sprintf("%.6f", this.get(i))
}
func (this *DoubleList) GetObject(i int) interface{} {
	return nil
}
func (this *DoubleList) GetValue(i int) value.Value {
	return value.NewDoubleValue(this.get(i))
}

func (this *DoubleList) Write(out *io.DataOutputX) {
	out.WriteInt3(int32(this.Size()))
	for i := 0; i < this.size; i++ {
		out.WriteDouble(this.get(i))
	}
}
func (this *DoubleList) Read(in *io.DataInputX) {
	count := int(in.ReadInt3())
	for i := 0; i < count; i++ {
		this.add(in.ReadDouble())
	}
}

func (this *DoubleList) Size() int {
	return this.size
}

func (this *DoubleList) Sorting(asc bool) []int {
	this.lock.Lock()
	defer this.lock.Unlock()

	table := make([]*DoubleListKeyVal, this.size)
	for i := 0; i < this.size; i++ {
		table[i] = &DoubleListKeyVal{i, this.get(i)}
	}

	c := func(o1, o2 *DoubleListKeyVal) bool {
		if asc {
			if compare.CompareToDouble(o1.value, o2.value) > 0 {
				return false
			} else {
				return true
			}
		} else {
			if compare.CompareToDouble(o2.value, o1.value) > 0 {
				return false
			} else {
				return true
			}
		}
	}

	sort.Sort(DoubleListSortable{compare: c, data: table})

	out := make([]int, this.size)
	for i := 0; i < this.size; i++ {
		out[i] = table[i].key
	}
	return out
}
func (this *DoubleList) SortingAnyList(asc bool, child AnyList, childAsc bool) []int {
	this.lock.Lock()
	defer this.lock.Unlock()

	table := make([]*DoubleListKeyVal, this.size)
	for i := 0; i < this.size; i++ {
		table[i] = &DoubleListKeyVal{i, this.get(i)}
	}

	c := func(o1, o2 *DoubleListKeyVal) bool {
		var rt int
		if asc {
			rt = compare.CompareToDouble(o1.value, o2.value)
		} else {
			rt = compare.CompareToDouble(o2.value, o1.value)
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

	sort.Sort(DoubleListSortable{compare: c, data: table})

	out := make([]int, this.size)
	for i := 0; i < this.size; i++ {
		out[i] = table[i].key
	}
	return out
}
func (this *DoubleList) Filtering(index []int) AnyList {
	out := NewDoubleList(this.size)
	sz := len(index)
	for i := 0; i < sz; i++ {
		out.add(this.get(index[i]))
	}
	return out
}

type DoubleListSortable struct {
	compare func(a, b *DoubleListKeyVal) bool
	data    []*DoubleListKeyVal
}

func (this DoubleListSortable) Len() int {
	return len(this.data)
}

func (this DoubleListSortable) Less(i, j int) bool {
	return this.compare(this.data[i], this.data[j])
}

func (this DoubleListSortable) Swap(i, j int) {
	this.data[i], this.data[j] = this.data[j], this.data[i]
}

type DoubleListKeyVal struct {
	key   int
	value float64
}

func DoubleListMain() {
	list := NewDoubleList(30)

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
