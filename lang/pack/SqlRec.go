package pack

import (
	"math"

	"github.com/whatap/golib/io"
	"github.com/whatap/golib/util/bitutil"
)

type SqlRec struct {
	// int32
	Dbc int32
	// int32
	Sql int32
	// byte
	SqlCrud byte
	// int32
	CountTotal int32
	// int32
	CountError int32
	// int32
	//CountActived int32
	// int64
	TimeSum int64
	// int64
	TimeStd int64
	// int32
	TimeMin int32
	// int32
	TimeMax int32
	// int64
	FetchCount int64
	// int64
	FetchTime int64
	// int64
	UpdateCount int64

	Service int32
}

func NewSqlRec() *SqlRec {
	p := new(SqlRec)

	return p
}

// func (this *StatSqlPack) Merge(o SqlRec)
func (this *SqlRec) Merge(o SqlRec) {
	this.CountTotal += o.CountTotal
	this.CountError += o.CountError
	//this.CountActived += o.CountActived
	this.TimeSum += o.TimeSum
	this.TimeStd = o.TimeStd
	this.TimeMax = int32(math.Max(float64(this.TimeMax), float64(o.TimeMax)))
	this.TimeMin = int32(math.Min(float64(this.TimeMin), float64(o.TimeMin)))
	this.FetchCount += o.FetchCount
	this.FetchTime += o.FetchTime
	this.UpdateCount += o.UpdateCount
	if this.Service == 0 {
		this.Service = o.Service
	}
}

// func (this *StatSqlPack) Key() int64 {
func (this *SqlRec) Key() int64 {
	return bitutil.Composite64(this.Dbc, this.Sql)
}

// func (this *StatSqlPack) SetDbcSql(dbc, sql int32)  SqlRec {
func (this *SqlRec) SetDbcSql(dbc, sql int32) *SqlRec {
	this.Dbc = dbc
	this.Sql = sql
	return this
}

func (this *SqlRec) GetDeviation() float64 {
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

func (this *SqlRec) Write(o *io.DataOutputX) {
	o.WriteInt(this.Dbc)
	o.WriteInt(this.Sql)
	o.WriteByte(this.SqlCrud)
	o.WriteDecimal(int64(this.CountTotal))
	o.WriteDecimal(int64(this.CountError))
	o.WriteDecimal(-1)
	o.WriteDecimal(this.TimeSum)
	o.WriteDecimal(this.TimeStd)
	o.WriteDecimal(int64(this.TimeMin))
	o.WriteDecimal(int64(this.TimeMax))
	o.WriteDecimal(this.FetchCount)
	o.WriteDecimal(this.FetchTime)
	o.WriteDecimal(this.UpdateCount)
	o.WriteDecimal(int64(this.Service))
}

func (this *SqlRec) Read(in *io.DataInputX) *SqlRec {
	this.Dbc = in.ReadInt()
	this.Sql = in.ReadInt()
	this.SqlCrud = in.ReadByte()
	this.CountTotal = int32(in.ReadDecimal())
	this.CountError = int32(in.ReadDecimal())

	ver := int32(in.ReadDecimal())
	this.TimeSum = in.ReadDecimal()
	this.TimeStd = in.ReadDecimal()
	this.TimeMin = int32(in.ReadDecimal())
	this.TimeMax = int32(in.ReadDecimal())
	this.FetchCount = in.ReadDecimal()
	this.FetchTime = in.ReadDecimal()
	this.UpdateCount = in.ReadDecimal()
	if ver >= 0 {
		return this
	}
	this.Service = int32(in.ReadDecimal())
	return this
}
