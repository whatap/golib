package pack

import (
	"github.com/whatap/golib/io"
)

type TimeCount struct {
	Count int32
	Error int32
	Time  int64
}

func NewTimeCountDefault() *TimeCount {
	p := new(TimeCount)
	return p
}

func NewTimeCount(count int32, err int32, time int64) *TimeCount {
	p := new(TimeCount)
	p.Count = count
	p.Error = err
	p.Time = time

	return p
}

func (this *TimeCount) Add(time int32, err bool) {
	this.Count++
	this.Time += int64(time)
	if err {
		this.Error++
	}
}

func (this *TimeCount) Merge(o *TimeCount) {
	this.Count += o.Count
	this.Time += o.Time
	this.Error += o.Error
}

func (this *TimeCount) Copy() *TimeCount {
	return NewTimeCount(this.Count, this.Error, this.Time)
}

func (this *TimeCount) Read(in *io.DataInputX) *TimeCount {

	this.Count = int32(in.ReadDecimal())
	this.Error = int32(in.ReadDecimal())
	this.Time = in.ReadDecimal()
	return this
}

func (this *TimeCount) Write(o *io.DataOutputX) {
	o.WriteDecimal(int64(this.Count))
	o.WriteDecimal(int64(this.Error))
	o.WriteDecimal(this.Time)
}
