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

type StringList struct {
	size  int
	table []string
	lock  sync.Mutex
}

func NewStringListDefault() *StringList {
	o := new(StringList)
	o.table = make([]string, 0)
	return o
}

func NewStringList(initialCapa int) *StringList {
	o := new(StringList)
	o.table = make([]string, initialCapa)
	return o
}

func (this *StringList) GetType() byte {
	return ANYLIST_STRING
}

func (this *StringList) ensure(minCapacity int) {
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

		newTable := make([]string, newSize)
		copy(newTable, this.table)
		this.table = newTable
	}
}

func (this *StringList) get(i int) string {
	if i >= this.size {
		panic(fmt.Sprintln("Index out of bounds ", "Index: ", i, ", Size: ", this.size))
	}
	return this.table[i]
}

func (this *StringList) set(i int, v string) string {
	if i >= this.size {
		panic(fmt.Sprintln("Index out of bounds Index: ", i, ", Size: ", this.size))
	}

	ov := this.table[i]
	this.table[i] = v
	return ov
}

func (this *StringList) remove(i int) string {
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
	this.table[last] = ""
	return ov
}

func (this *StringList) add(e string) {
	this.ensure(this.size + 1)
	//table[size++] = e;
	this.table[this.size] = e
	this.size++
}

func (this *StringList) AddAll(other *StringList) {
	this.ensure(this.size + other.size)
	for i := 0; i < other.size; i++ {
		//this.table[this.size++] = other.table[i];
		this.table[this.size] = other.table[i]
		this.size++
	}
}

func (this *StringList) AddAllArray(other []string) {
	this.ensure(this.size + len(other))
	n := len(other)
	for i := 0; i < n; i++ {
		//this.table[this.size++] = other[i];
		this.table[this.size] = other[i]
		this.size++
	}
}

func (this *StringList) ToArray() []string {
	newArray := make([]string, this.size)
	copy(newArray, this.table)
	return newArray
}

func (this *StringList) ToString() string {
	return fmt.Sprintf("%v", this.table)
}

func (this *StringList) Sort() {
	if this.size > 0 {
		//Arrays.sort(this.table, 0, size-1)
		//sort.Ints(this.table)
	}
}

func (this *StringList) SetInt(i, v int) {
	this.set(i, strconv.Itoa(v))
}
func (this *StringList) SetFloat(i int, v float32) {
	this.set(i, strconv.FormatFloat(float64(v), 'f', 6, 32))
	//this.set(i, fmt.Sprintf("%.6f", this.get(i)))
}
func (this *StringList) SetLong(i int, v int64) {
	this.set(i, strconv.FormatInt(v, 10))
}
func (this *StringList) SetDouble(i int, v float64) {
	this.set(i, strconv.FormatFloat(v, 'f', 6, 64))
	//this.set(i, fmt.Sprintf("%.6f", this.get(i)))
}
func (this *StringList) SetString(i int, v string) {
	this.set(i, v)
}

func (this *StringList) AddInt(v int) {
	this.add(strconv.Itoa(v))
}
func (this *StringList) AddFloat(v float32) {
	this.add(strconv.FormatFloat(float64(v), 'f', 6, 32))
	//this.add(fmt.Sprintf("%.6f", this.get(i)))
}
func (this *StringList) AddLong(v int64) {
	this.add(strconv.FormatInt(v, 10))
}
func (this *StringList) AddDouble(v float64) {
	this.add(strconv.FormatFloat(float64(v), 'f', 6, 64))
	//this.add(i, fmt.Sprintf("%.6f", this.get(i)))
}
func (this *StringList) AddString(v string) {
	this.add(v)
}

func (this *StringList) GetInt(i int) int {
	n, err := strconv.Atoi(this.get(i))
	if err == nil {
		return n
	} else {
		panic(" Parse error string(" + this.get(i) + ") to int")
	}
}
func (this *StringList) GetFloat(i int) float32 {
	n, err := strconv.ParseFloat(this.get(i), 32)
	if err == nil {
		return float32(n)
	} else {
		panic(" Parse error string(" + this.get(i) + ") to float32")
	}
}
func (this *StringList) GetDouble(i int) float64 {
	n, err := strconv.ParseFloat(this.get(i), 64)
	if err == nil {
		return n
	} else {
		panic(" Parse error string(" + this.get(i) + ") to float64")
	}
}
func (this *StringList) GetLong(i int) int64 {
	n, err := strconv.ParseInt(this.get(i), 10, 64)
	if err == nil {
		return n
	} else {
		panic(" Parse error string(" + this.get(i) + ") to int64")
	}
}
func (this *StringList) GetString(i int) string {
	return this.get(i)
}
func (this *StringList) GetObject(i int) interface{} {
	return nil
}
func (this *StringList) GetValue(i int) value.Value {
	return value.NewTextValue(this.get(i))
}

func (this *StringList) Write(out *io.DataOutputX) {
	out.WriteInt3(int32(this.Size()))
	for i := 0; i < this.size; i++ {
		out.WriteText(this.get(i))
	}
}
func (this *StringList) Read(in *io.DataInputX) {
	count := int(in.ReadInt3())
	for i := 0; i < count; i++ {
		this.add(in.ReadText())
	}
}

func (this *StringList) Size() int {
	return this.size
}

func (this *StringList) Sorting(asc bool) []int {
	this.lock.Lock()
	defer this.lock.Unlock()

	table := make([]*StringListKeyVal, this.size)
	for i := 0; i < this.size; i++ {
		table[i] = &StringListKeyVal{i, this.get(i)}
	}

	c := func(o1, o2 *StringListKeyVal) bool {
		if asc {
			if compare.CompareToString(o1.value, o2.value) > 0 {
				return false
			} else {
				return true
			}
		} else {
			if compare.CompareToString(o2.value, o1.value) > 0 {
				return false
			} else {
				return true
			}
		}
	}

	sort.Sort(StringListSortable{compare: c, data: table})

	out := make([]int, this.size)
	for i := 0; i < this.size; i++ {
		out[i] = table[i].key
	}
	return out
}
func (this *StringList) SortingAnyList(asc bool, child AnyList, childAsc bool) []int {
	this.lock.Lock()
	defer this.lock.Unlock()

	table := make([]*StringListKeyVal, this.size)
	for i := 0; i < this.size; i++ {
		table[i] = &StringListKeyVal{i, this.get(i)}
	}

	c := func(o1, o2 *StringListKeyVal) bool {
		var rt int
		if asc {
			rt = compare.CompareToString(o1.value, o2.value)
		} else {
			rt = compare.CompareToString(o2.value, o1.value)
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

	sort.Sort(StringListSortable{compare: c, data: table})

	out := make([]int, this.size)
	for i := 0; i < this.size; i++ {
		out[i] = table[i].key
	}
	return out
}
func (this *StringList) Filtering(index []int) AnyList {
	out := NewStringList(this.size)
	sz := len(index)
	for i := 0; i < sz; i++ {
		out.add(this.get(index[i]))
	}
	return out
}

type StringListSortable struct {
	compare func(a, b *StringListKeyVal) bool
	data    []*StringListKeyVal
}

func (this StringListSortable) Len() int {
	return len(this.data)
}

func (this StringListSortable) Less(i, j int) bool {
	return this.compare(this.data[i], this.data[j])
}

func (this StringListSortable) Swap(i, j int) {
	this.data[i], this.data[j] = this.data[j], this.data[i]
}

type StringListKeyVal struct {
	key   int
	value string
}

func StringListMain() {
	list := NewStringList(30)

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
