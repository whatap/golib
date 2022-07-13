package value

import (
	"fmt"
	"math"

	"github.com/whatap/golib/io"
)

type DoubleSummary struct {
	Sum   float64
	Count int32
	Min   float64
	Max   float64
}

func NewDoubleSummary() *DoubleSummary {
	p := new(DoubleSummary)
	return p
}

func (this *DoubleSummary) CompareTo(o Value) int {
	if o != nil && o.GetValueType() == this.GetValueType() {
		if this.Sum == o.(*DoubleSummary).Sum && this.Count == o.(*DoubleSummary).Count {
			return 0
		}
		if this.Sum < o.(*DoubleSummary).Sum {
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

func (this *DoubleSummary) Equals(o Value) bool {
	if o != nil && o.GetValueType() == this.GetValueType() {
		return this.Sum == o.(*DoubleSummary).Sum && this.Count == o.(*DoubleSummary).Count
	}
	return false
}

func (this *DoubleSummary) GetValueType() byte {
	return VALUE_DOUBLE_SUMMARY
}

func (this *DoubleSummary) Write(out *io.DataOutputX) {
	// import cycle not allowd (logutil 에서 value 사용) 여기서 recover 안함
	//	defer func() {
	//		if r := recover(); r != nil {
	//			//logutil.Println("WA130001", "Recover ", r)
	//		}
	//	}()
	out.WriteDouble(this.Sum)
	out.WriteInt(this.Count)
	out.WriteDouble(this.Min)
	out.WriteDouble(this.Max)
}

func (this *DoubleSummary) Read(in *io.DataInputX) {
	// import cycle not allowd (logutil 에서 value 사용) 여기서 recover 안함
	//	defer func() {
	//		if r := recover(); r != nil {
	//			//logutil.Println("WA130002", "Recover ", r)
	//		}
	//	}()

	this.Sum = in.ReadDouble()
	this.Count = in.ReadInt()
	this.Min = in.ReadDouble()
	this.Max = in.ReadDouble()
}

func (this *DoubleSummary) ToString() string {
	return fmt.Sprintln("[sum=", this.Sum, ",count=", this.Count, ",min=", this.Min, ",max=", this.Max, "]")
}

//	public Object toJavaObject() {
//		return this;
//	}
//

func (this *DoubleSummary) AddCount() {
	this.Count++
}

// java.lang.Number 없음
//	public SummaryValue add(Number value) {
//		if(value==null)
//			return this;
//		if (this.count == 0) {
//			this.sum = value.longValue();
//			this.count = 1;
//			this.max = value.longValue();
//			this.min = value.longValue();
//		} else {
//			this.sum += value.doubleValue();
//			this.count++;
//			this.max = Math.max(this.max, value.longValue());
//			this.min = Math.min(this.min, value.longValue());
//		}
//		return this;
//	}

func (this *DoubleSummary) Add(other SummaryValue) SummaryValue {
	if other == nil || other.GetCount() == 0 {
		return this
	}

	this.Count += other.GetCount()
	this.Sum += other.DoubleSum()
	this.Min = math.Min(this.Min, other.DoubleMin())
	this.Max = math.Max(this.Max, other.DoubleMax())

	return this
}

func (this *DoubleSummary) LongSum() int64 {
	return int64(this.Sum)
}
func (this *DoubleSummary) LongMin() int64 {
	return int64(this.Min)
}
func (this *DoubleSummary) LongMax() int64 {
	return int64(this.Max)
}
func (this *DoubleSummary) LongAvg() int64 {
	if this.Count == 0 {
		return 0
	} else {
		return int64(this.Sum / float64(this.Count))
	}
}
func (this *DoubleSummary) DoubleSum() float64 {
	return this.Sum
}
func (this *DoubleSummary) DoubleMin() float64 {
	return this.Min
}
func (this *DoubleSummary) DoubleMax() float64 {
	return this.Max
}
func (this *DoubleSummary) DoubleAvg() float64 {
	if this.Count == 0 {
		return 0
	} else {
		return this.Sum / float64(this.Count)
	}
}
func (this *DoubleSummary) GetCount() int32 {
	return this.Count
}
