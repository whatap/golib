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

type FloatList struct {
	size  int
	table []float32
	lock  sync.Mutex
}

func NewFloatListDefault() *FloatList {
	o := new(FloatList)
	o.table = make([]float32, 0)
	return o
}

func NewFloatList(initialCapa int) *FloatList {
	o := new(FloatList)
	o.table = make([]float32, initialCapa)
	return o
}

func (this *FloatList) GetType() byte {
	return ANYLIST_FLOAT
}

func (this *FloatList) ensure(minCapacity int) {
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

		newTable := make([]float32, newSize)
		copy(newTable, this.table)
		this.table = newTable
	}
}

func (this *FloatList) get(i int) float32 {
	if i >= this.size {
		panic(fmt.Sprintln("Index out of bounds ", "Index: ", i, ", Size: ", this.size))
	}
	return this.table[i]
}

func (this *FloatList) set(i int, v float32) float32 {
	if i >= this.size {
		panic(fmt.Sprintln("Index out of bounds Index: ", i, ", Size: ", this.size))
	}

	ov := this.table[i]
	this.table[i] = v
	return ov
}

func (this *FloatList) remove(i int) float32 {
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

func (this *FloatList) add(e float32) {
	this.ensure(this.size + 1)
	//table[size++] = e;
	this.table[this.size] = e
	this.size++
}

func (this *FloatList) AddAll(other *FloatList) {
	this.ensure(this.size + other.size)
	for i := 0; i < other.size; i++ {
		//this.table[this.size++] = other.table[i];
		this.table[this.size] = other.table[i]
		this.size++
	}
}

func (this *FloatList) AddAllArray(other []float32) {
	this.ensure(this.size + len(other))
	n := len(other)
	for i := 0; i < n; i++ {
		//this.table[this.size++] = other[i];
		this.table[this.size] = other[i]
		this.size++
	}
}

func (this *FloatList) ToArray() []float32 {
	newArray := make([]float32, this.size)
	copy(newArray, this.table)
	return newArray
}

func (this *FloatList) ToString() string {
	return fmt.Sprintf("%v", this.table)
}

func (this *FloatList) Sort() {
	if this.size > 0 {
		//Arrays.sort(this.table, 0, size-1)
		//sort.Ints(this.table)
	}
}

func (this *FloatList) SetInt(i, v int) {
	this.set(i, float32(v))
}
func (this *FloatList) SetFloat(i int, v float32) {
	this.set(i, v)
}
func (this *FloatList) SetLong(i int, v int64) {
	this.set(i, float32(v))
}
func (this *FloatList) SetDouble(i int, v float64) {
	this.set(i, float32(v))
}
func (this *FloatList) SetString(i int, v string) {
	n, err := strconv.ParseFloat(v, 32)
	if err == nil {
		this.set(i, float32(n))
	} else {
		panic(" Parse error string(" + v + ") to float32")
	}

}

func (this *FloatList) AddInt(v int) {
	this.add(float32(v))
}
func (this *FloatList) AddFloat(v float32) {
	this.add(v)
}
func (this *FloatList) AddLong(v int64) {
	this.add(float32(v))
}
func (this *FloatList) AddDouble(v float64) {
	this.add(float32(v))
}
func (this *FloatList) AddString(v string) {

	n, err := strconv.ParseFloat(v, 32)
	if err == nil {
		this.add(float32(n))
	} else {
		panic(" Parse error string(" + v + ") to float32")
	}
}

func (this *FloatList) GetInt(i int) int {
	return int(this.get(i))
}
func (this *FloatList) GetFloat(i int) float32 {
	return float32(this.get(i))
}
func (this *FloatList) GetDouble(i int) float64 {
	return float64(this.get(i))
}
func (this *FloatList) GetLong(i int) int64 {
	return int64(this.get(i))
}
func (this *FloatList) GetString(i int) string {
	return strconv.FormatFloat(float64(this.get(i)), 'f', 6, 32)
	//return fmt.Sprintf("%.6f", this.get(i))
}
func (this *FloatList) GetObject(i int) interface{} {
	return nil
}
func (this *FloatList) GetValue(i int) value.Value {
	return value.NewFloatValue(this.get(i))
}

func (this *FloatList) Write(out *io.DataOutputX) {
	out.WriteInt3(int32(this.Size()))
	for i := 0; i < this.size; i++ {
		out.WriteFloat(this.get(i))
	}
}
func (this *FloatList) Read(in *io.DataInputX) {
	count := int(in.ReadInt3())
	for i := 0; i < count; i++ {
		this.add(in.ReadFloat())
	}
}

func (this *FloatList) Size() int {
	return this.size
}

func (this *FloatList) Sorting(asc bool) []int {
	this.lock.Lock()
	defer this.lock.Unlock()

	table := make([]*FloatListKeyVal, this.size)
	for i := 0; i < this.size; i++ {
		table[i] = &FloatListKeyVal{i, this.get(i)}
	}

	c := func(o1, o2 *FloatListKeyVal) bool {
		if asc {
			if compare.CompareToFloat(o1.value, o2.value) > 0 {
				return false
			} else {
				return true
			}
		} else {
			if compare.CompareToFloat(o2.value, o1.value) > 0 {
				return false
			} else {
				return true
			}
		}
	}

	sort.Sort(FloatListSortable{compare: c, data: table})

	out := make([]int, this.size)
	for i := 0; i < this.size; i++ {
		out[i] = table[i].key
	}
	return out
}
func (this *FloatList) SortingAnyList(asc bool, child AnyList, childAsc bool) []int {
	this.lock.Lock()
	defer this.lock.Unlock()

	table := make([]*FloatListKeyVal, this.size)
	for i := 0; i < this.size; i++ {
		table[i] = &FloatListKeyVal{i, this.get(i)}
	}

	c := func(o1, o2 *FloatListKeyVal) bool {
		var rt int
		if asc {
			rt = compare.CompareToFloat(o1.value, o2.value)
		} else {
			rt = compare.CompareToFloat(o2.value, o1.value)
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

	sort.Sort(FloatListSortable{compare: c, data: table})

	out := make([]int, this.size)
	for i := 0; i < this.size; i++ {
		out[i] = table[i].key
	}
	return out
}
func (this *FloatList) Filtering(index []int) AnyList {
	out := NewFloatList(this.size)
	sz := len(index)
	for i := 0; i < sz; i++ {
		out.add(this.get(index[i]))
	}
	return out
}

type FloatListSortable struct {
	compare func(a, b *FloatListKeyVal) bool
	data    []*FloatListKeyVal
}

func (this FloatListSortable) Len() int {
	return len(this.data)
}

func (this FloatListSortable) Less(i, j int) bool {
	return this.compare(this.data[i], this.data[j])
}

func (this FloatListSortable) Swap(i, j int) {
	this.data[i], this.data[j] = this.data[j], this.data[i]
}

type FloatListKeyVal struct {
	key   int
	value float32
}

func FloatListMain() {
	list := NewFloatList(30)

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
