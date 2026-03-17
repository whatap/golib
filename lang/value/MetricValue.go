package value

import (
	"fmt"
	"math"

	"github.com/whatap/golib/io"
)

const VALUE_METRIC = 82

type MetricValue struct {
	Count int32
	Sum   float64
	Min   float64
	Max   float64
	Last  float64
}

func NewMetricValue() *MetricValue {
	return new(MetricValue)
}

func (this *MetricValue) GetValueType() byte {
	return VALUE_METRIC
}

func (this *MetricValue) Write(out *io.DataOutputX) {
	if this.Count == 0 {
		out.WriteByte(0)
	} else {
		out.WriteByte(1)
		out.WriteDecimal(int64(this.Count))
		out.WriteDouble(this.Sum)
		out.WriteFloat(float32(this.Min))
		out.WriteFloat(float32(this.Max))
		out.WriteFloat(float32(this.Last))
	}
}

func (this *MetricValue) Read(in *io.DataInputX) {
	mode := in.ReadByte()
	switch mode {
	case 0:
		return
	case 1:
		this.Count = int32(in.ReadDecimal())
		this.Sum = in.ReadDouble()
		this.Min = float64(in.ReadFloat())
		this.Max = float64(in.ReadFloat())
		this.Last = float64(in.ReadFloat())
	}
}

func (this *MetricValue) CompareTo(o Value) int {
	if o != nil && o.GetValueType() == this.GetValueType() {
		other := o.(*MetricValue)
		if this.Sum == other.Sum {
			return 0
		}
		if this.Sum < other.Sum {
			return 1
		}
		return -1
	}
	if o == nil {
		return 1
	}
	return int(this.GetValueType() - o.GetValueType())
}

func (this *MetricValue) Equals(o Value) bool {
	if o != nil && o.GetValueType() == this.GetValueType() {
		other := o.(*MetricValue)
		return this.Count == other.Count && this.Sum == other.Sum
	}
	return false
}

func (this *MetricValue) Add(v float64) *MetricValue {
	if this.Count == 0 {
		this.Count = 1
		this.Sum = v
		this.Max = v
		this.Min = v
	} else {
		this.Count++
		this.Sum += v
		this.Max = math.Max(this.Max, v)
		this.Min = math.Min(this.Min, v)
	}
	this.Last = v
	return this
}

func (this *MetricValue) AddLong(v int64) *MetricValue {
	return this.Add(float64(v))
}

func (this *MetricValue) Merge(other *MetricValue) *MetricValue {
	if other == nil {
		return this
	}
	if this.Count == 0 {
		this.Count = other.Count
		this.Sum = other.Sum
		this.Max = other.Max
		this.Min = other.Min
	} else {
		this.Count += other.Count
		this.Sum += other.Sum
		this.Max = math.Max(this.Max, other.Max)
		this.Min = math.Min(this.Min, other.Min)
	}
	this.Last = other.Last
	return this
}

func (this *MetricValue) Avg() float64 {
	if this.Count == 0 {
		return 0
	}
	return this.Sum / float64(this.Count)
}

func (this *MetricValue) GetSum() float64 { return this.Sum }

func (this *MetricValue) GetMin() float64 {
	if this.Count == 1 {
		return this.Sum
	}
	return this.Min
}

func (this *MetricValue) GetMax() float64 {
	if this.Count == 1 {
		return this.Sum
	}
	return this.Max
}

func (this *MetricValue) GetLast() float64 {
	if this.Count == 1 {
		return this.Sum
	}
	return this.Last
}

func (this *MetricValue) GetCount() int32 { return this.Count }

func (this *MetricValue) String() string {
	return fmt.Sprintf("[sum=%.2f,count=%d,min=%.2f,max=%.2f,last=%.2f]",
		this.Sum, this.Count, this.Min, this.Max, this.Last)
}
