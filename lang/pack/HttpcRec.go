package pack

import (
	"fmt"
	"math"

	"github.com/whatap/golib/io"
	"github.com/whatap/golib/lang/variable"
)

type HttpcRec struct {
	Url        int32
	Host       int32
	Port       int32
	CountTotal int32
	CountError int32
	//CountActived int32
	TimeSum int64
	TimeStd int64
	TimeMin int32
	TimeMax int32

	Service int32
}

func NewHttpcRec() *HttpcRec {
	p := new(HttpcRec)
	return p
}

func (this *HttpcRec) Merge(o *HttpcRec) {
	this.CountTotal += o.CountTotal
	this.CountError += o.CountError
	//this.CountActived += o.CountActived
	this.TimeSum += o.TimeSum
	this.TimeStd = o.TimeStd
	if o.TimeMax > this.TimeMax {
		this.TimeMax = o.TimeMax
	}
	if o.TimeMin < this.TimeMin {
		this.TimeMin = o.TimeMin
	}

	if this.Service == 0 {
		this.Service = o.Service
	}
}

func (this *HttpcRec) KeyHostPort() *variable.I2 {
	return variable.NewI2(this.Host, this.Port)
}

func (this *HttpcRec) KeyFull() *variable.I3 {
	return variable.NewI3(this.Url, this.Host, this.Port)
}

func (this *HttpcRec) SetUrlHostPort(url, host, port int32) {
	this.Url = url
	this.Host = host
	this.Port = port
}

func (this *HttpcRec) GetDeviation() float64 {
	if this.CountTotal == 0 {
		return 0
	}
	avg := this.TimeSum / int64(this.CountTotal)
	variation := (this.TimeStd - (2 * avg * this.TimeSum) + (int64(this.CountTotal) * avg * avg)) / int64(this.CountTotal)
	//double ret = Math.sqrt(variation);
	//return ret == Double.NaN ? 0 : ret;
	ret := math.Sqrt(float64(variation))
	if ret == math.NaN() {
		return 0
	} else {
		return ret
	}
}

func (this *HttpcRec) ToString() string {

	return fmt.Sprintln("HttpRec ", "url=", this.Url, ",host=", this.Host, ",port=", this.Port, ",count_total=", this.CountTotal, ",count_error=", this.CountError,
		//",count_actived=" ,this.count_actived,
		",time_sum=", this.TimeSum, ",time_std=", this.TimeStd, ",time_min=", this.TimeMin, ",time_max=", this.TimeMax)
}

func (this *HttpcRec) Write(o *io.DataOutputX) {
	o.WriteInt(this.Url)
	o.WriteInt(this.Host)
	o.WriteInt(this.Port)
	o.WriteDecimal(int64(this.CountTotal))
	o.WriteDecimal(int64(this.CountError))
	o.WriteDecimal(-1) //version은 음수(-)로 증가시켜야 한다.
	o.WriteDecimal(this.TimeSum)
	o.WriteDecimal(this.TimeStd)
	o.WriteDecimal(int64(this.TimeMin))
	o.WriteDecimal(int64(this.TimeMax))

	o.WriteDecimal(int64(this.Service))
}
func (this *HttpcRec) Read(in *io.DataInputX) *HttpcRec {
	this.Url = in.ReadInt()
	this.Host = in.ReadInt()
	this.Port = in.ReadInt()
	this.CountTotal = int32(in.ReadDecimal())
	this.CountError = int32(in.ReadDecimal())
	ver := int32(in.ReadDecimal())
	this.TimeSum = in.ReadDecimal()
	this.TimeStd = in.ReadDecimal()
	this.TimeMin = int32(in.ReadDecimal())
	this.TimeMax = int32(in.ReadDecimal())
	if ver >= 0 {
		return this
	}
	this.Service = int32(in.ReadDecimal())
	return this
}
