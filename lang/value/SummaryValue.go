package value

import ()

const (
	BYTE_LEN = 29
)

type SummaryValue interface {
	Value
	AddCount()
	// java.lang.Number 안함.
	
	//Add(Number value) *SummaryValue
	Add(num SummaryValue) SummaryValue
	LongSum() int64
	LongMin() int64
	LongMax() int64
	LongAvg() int64
	DoubleSum() float64
	DoubleMin() float64
	DoubleMax() float64
	DoubleAvg() float64
	GetCount() int32
}
