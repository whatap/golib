package value

import (
	"fmt"
	"math"

	"github.com/whatap/golib/io"
)

type LongSummary struct {
	Sum   int64
	Count int32
	Min   int64
	Max   int64
}

func NewLongSummary() *LongSummary {
	p := new(LongSummary)
	return p
}

func (this *LongSummary) CompareTo(o Value) int {
	if o != nil && o.GetValueType() == this.GetValueType() {
		if this.Sum == o.(*LongSummary).Sum && this.Count == o.(*LongSummary).Count {
			return 0
		}
		if this.Sum < o.(*LongSummary).Sum {
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

func (this *LongSummary) Equals(o Value) bool {
	if o != nil && o.GetValueType() == this.GetValueType() {
		return this.Sum == o.(*LongSummary).Sum && this.Count == o.(*LongSummary).Count
	}
	return false
}

func (this *LongSummary) GetValueType() byte {
	return VALUE_LONG_SUMMARY
}

func (this *LongSummary) Write(out *io.DataOutputX) {
	// import cycle not allowd (logutil 에서 value 사용) 여기서 recover 안함
	//	defer func() {
	//		if r := recover(); r != nil {
	//			//logutil.Println("WA130001", "Recover ", r)
	//		}
	//	}()
	out.WriteLong(this.Sum)
	out.WriteInt(this.Count)
	out.WriteLong(this.Min)
	out.WriteLong(this.Max)
}

func (this *LongSummary) Read(in *io.DataInputX) {
	// import cycle not allowd (logutil 에서 value 사용) 여기서 recover 안함
	//	defer func() {
	//		if r := recover(); r != nil {
	//			//logutil.Println("WA130002", "Recover ", r)
	//		}
	//	}()

	this.Sum = in.ReadLong()
	this.Count = in.ReadInt()
	this.Min = in.ReadLong()
	this.Max = in.ReadLong()
}

func (this *LongSummary) ToString() string {
	return fmt.Sprintln("[sum=", this.Sum, ",count=", this.Count, ",min=", this.Min, ",max=", this.Max, "]")
}

//	public Object toJavaObject() {
//		return this;
//	}
//

func (this *LongSummary) AddCount() {
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

func (this *LongSummary) Add(other SummaryValue) SummaryValue {
	if other == nil || other.GetCount() == 0 {
		return this
	}

	this.Count += other.GetCount()
	this.Sum += other.LongSum()
	this.Min = int64(math.Min(float64(this.Min), float64(other.LongMin())))
	this.Max = int64(math.Max(float64(this.Max), float64(other.LongMax())))

	return this
}

func (this *LongSummary) LongSum() int64 {
	return this.Sum
}
func (this *LongSummary) LongMin() int64 {
	return this.Min
}
func (this *LongSummary) LongMax() int64 {
	return this.Max
}
func (this *LongSummary) LongAvg() int64 {
	if this.Count == 0 {
		return 0
	} else {
		return int64(float64(this.Sum) / float64(this.Count))
	}
}
func (this *LongSummary) DoubleSum() float64 {
	return float64(this.Sum)
}
func (this *LongSummary) DoubleMin() float64 {
	return float64(this.Min)
}
func (this *LongSummary) DoubleMax() float64 {
	return float64(this.Max)
}
func (this *LongSummary) DoubleAvg() float64 {
	if this.Count == 0 {
		return 0
	} else {
		return float64(this.Sum) / float64(this.Count)
	}
}
func (this *LongSummary) GetCount() int32 {
	return this.Count
}
