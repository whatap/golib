package value

import (
	"fmt"
	"math"

	"github.com/whatap/golib/io"
)

type FloatSummary struct {
	Sum   float32
	Count int32
	Min   float32
	Max   float32
}

func NewFloatSummary() *FloatSummary {
	return new(FloatSummary)
}

func (this *FloatSummary) CompareTo(o Value) int {
	if o != nil && o.GetValueType() == this.GetValueType() {
		if this.Sum == o.(*FloatSummary).Sum && this.Count == o.(*FloatSummary).Count {
			return 0
		}
		if this.Sum < o.(*FloatSummary).Sum {
			return 1
		} else {
			return -1
		}
	}
	if o == nil {
		return 1
	} else {
		return int(this.GetValueType() - o.GetValueType())
	}
}

func (this *FloatSummary) Equals(o Value) bool {
	if o != nil && o.GetValueType() == this.GetValueType() {
		return this.Sum == o.(*FloatSummary).Sum && this.Count == o.(*FloatSummary).Count
	}
	return false
}

func (this *FloatSummary) GetValueType() byte {
	return FLOAT_SUMMARY
}

func (this *FloatSummary) Write(out *io.DataOutputX) {
	out.WriteInt(this.Count)
	out.WriteFloat(this.Sum)
	out.WriteFloat(this.Min)
	out.WriteFloat(this.Max)
}

func (this *FloatSummary) Read(in *io.DataInputX) {
	this.Count = in.ReadInt()
	this.Sum = in.ReadFloat()
	this.Min = in.ReadFloat()
	this.Max = in.ReadFloat()
}

func (this *FloatSummary) String() string {
	return fmt.Sprintf("[sum=%f,count=%d,min=%f,max=%f]", this.Sum, this.Count, this.Min, this.Max)
}

func (this *FloatSummary) AddCount() {
	this.Count++
}

func (this *FloatSummary) AddFloat(v float32) *FloatSummary {
	if this.Count == 0 {
		this.Sum = v
		this.Count = 1
		this.Max = v
		this.Min = v
	} else {
		this.Sum += v
		this.Count++
		this.Max = float32(math.Max(float64(this.Max), float64(v)))
		this.Min = float32(math.Min(float64(this.Min), float64(v)))
	}
	return this
}

func (this *FloatSummary) Add(other SummaryValue) SummaryValue {
	if other == nil || other.GetCount() == 0 {
		return this
	}
	this.Count += other.GetCount()
	this.Sum += float32(other.DoubleSum())
	this.Min = float32(math.Min(float64(this.Min), other.DoubleMin()))
	this.Max = float32(math.Max(float64(this.Max), other.DoubleMax()))
	return this
}

func (this *FloatSummary) LongSum() int64     { return int64(this.Sum) }
func (this *FloatSummary) LongMin() int64     { return int64(this.Min) }
func (this *FloatSummary) LongMax() int64     { return int64(this.Max) }
func (this *FloatSummary) DoubleSum() float64 { return float64(this.Sum) }
func (this *FloatSummary) DoubleMin() float64 { return float64(this.Min) }
func (this *FloatSummary) DoubleMax() float64 { return float64(this.Max) }
func (this *FloatSummary) GetCount() int32    { return this.Count }

func (this *FloatSummary) LongAvg() int64 {
	if this.Count == 0 {
		return 0
	}
	return int64(this.Sum / float32(this.Count))
}

func (this *FloatSummary) DoubleAvg() float64 {
	if this.Count == 0 {
		return 0
	}
	return float64(this.Sum) / float64(this.Count)
}
