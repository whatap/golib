package list

import (
	"math"

	"github.com/whatap/golib/io"
	"github.com/whatap/golib/lang/value"
	"github.com/whatap/golib/util/compare"
)

const (
	ANYLIST_DEFAULT_CAPACITY = 10
	ANYLIST_MAX_SIZE         = math.MaxInt32 - 8

	ANYLIST_INT    = 1
	ANYLIST_LONG   = 2
	ANYLIST_FLOAT  = 3
	ANYLIST_DOUBLE = 4
	ANYLIST_STRING = 5
)

type AnyList interface {
	GetType() byte

	SetInt(i, v int)
	SetFloat(i int, v float32)
	SetLong(i int, v int64)
	SetDouble(i int, v float64)
	SetString(i int, v string)

	AddInt(v int)
	AddFloat(v float32)
	AddLong(v int64)
	AddDouble(v float64)
	AddString(v string)

	GetInt(i int) int
	GetFloat(i int) float32
	GetDouble(i int) float64
	GetLong(i int) int64
	GetString(i int) string
	GetObject(i int) interface{}
	GetValue(i int) value.Value

	Write(out *io.DataOutputX)
	Read(in *io.DataInputX)

	Size() int

	Sorting(asc bool) []int
	SortingAnyList(asc bool, child AnyList, childAsc bool) []int
	Filtering(index []int) AnyList
}

func CompareChild(child AnyList, ord bool, i1 int, i2 int) int {
	if child.GetType() == ANYLIST_STRING {
		if ord {
			return compare.CompareToString(child.GetString(i1), child.GetString(i2))
		} else {
			return compare.CompareToString(child.GetString(i2), child.GetString(i1))
		}
	} else {
		if ord {
			return compare.CompareToDouble(child.GetDouble(i1), child.GetDouble(i2))
		} else {
			return compare.CompareToDouble(child.GetDouble(i2), child.GetDouble(i1))
		}
	}
}
