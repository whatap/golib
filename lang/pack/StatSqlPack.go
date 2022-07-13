package pack

import (
	"fmt"
	//"math"
	"container/list"

	"github.com/whatap/golib/io"
	"github.com/whatap/golib/util/hmap"
	"github.com/whatap/golib/util/stringutil"
)

type StatSqlPack struct {
	AbstractPack
	// []byte
	Records []byte
	// int32
	RecordCount int32
}

func NewStatSqlPack() *StatSqlPack {
	p := new(StatSqlPack)
	return p
}

// Implements Pack
func (this *StatSqlPack) GetPackType() int16 {
	return PACK_STAT_SQL
}

// String()
//func (this *StatSqlPack) String() string {
//	return this.ToString()
//}

func (this *StatSqlPack) ToString() string {
	sb := stringutil.NewStringBuffer()
	sb.Append("StatSql ")
	sb.Append(this.AbstractPack.ToString())
	sb.Append(fmt.Sprintln(",records=", this.RecordCount))
	sb.Append(fmt.Sprintln(",bytes=", len(this.Records)))
	return sb.ToString()
}

func (this *StatSqlPack) Write(dout *io.DataOutputX) {
	this.AbstractPack.Write(dout)
	dout.WriteBlob(this.Records)
	dout.WriteDecimal(int64(this.RecordCount))
}

func (this *StatSqlPack) Read(din *io.DataInputX) {
	this.AbstractPack.Read(din)
	this.Records = din.ReadBlob()
	this.RecordCount = int32(din.ReadDecimal())
	//return this
}

func (this *StatSqlPack) SetRecords(size int, items hmap.Enumeration) *StatSqlPack {
	o := io.NewDataOutputX()
	o.WriteShort(int16(size))
	for i := 0; i < size; i++ {
		items.NextElement().(*SqlRec).Write(o)
	}
	this.Records = o.ToByteArray()
	this.RecordCount = int32(size)

	return this
}

// TODO 테스트
func (this *StatSqlPack) SetRecordsList(items *list.List) *StatSqlPack {
	o := io.NewDataOutputX()

	o.WriteShort(int16(items.Len()))
	for e := items.Front(); e != nil; e = e.Next() {
		e.Value.(*SqlRec).Write(o)
	}
	this.Records = o.ToByteArray()
	this.RecordCount = int32(items.Len())
	return this
}

func (this *StatSqlPack) GetRecords() *list.List {
	items := list.New()
	if this.Records == nil {
		return items
	}
	in := io.NewDataInputX(this.Records)
	size := int(in.ReadShort()) & 0xffff
	for i := 0; i < size; i++ {
		items.PushBack(NewSqlRec().Read(in))
	}

	return items
}

//func (this *StatSqlPack) WriteSql(o *io.DataOutputX, m *SqlRec) {
//	o.WriteInt(m.Dbc)
//	o.WriteInt(m.Sql)
//	o.WriteByte(m.SqlCrud)
//	o.WriteDecimal(int64(m.CountTotal))
//	o.WriteDecimal(int64(m.CountError))
//	//o.WriteDecimal(int64(m.CountActived))
//	o.WriteDecimal(m.TimeSum)
//	o.WriteDecimal(m.TimeStd)
//	o.WriteDecimal(int64(m.TimeMin))
//	o.WriteDecimal(int64(m.TimeMax))
//	o.WriteDecimal(m.FetchCount)
//	o.WriteDecimal(m.FetchTime)
//	o.WriteDecimal(m.UpdateCount)
//}

// TOOD
// return Java.List to []*SqlRec
//func (this *StatSqlPack) GetRecords() []*SqlRec {
//	in := io.NewDataInputX(this.Records)
//	sz := int(in.ReadShort())
//	items := make([]*SqlRec, sz)
//	for i := 0; i < sz; i++ {
//		items[i] = this.ReadSql(in)
//	}
//	return items
//}
//
//func (this *StatSqlPack) ReadSql(in *io.DataInputX) *SqlRec {
//	m := NewSqlRec()
//	m.Dbc = in.ReadInt()
//	m.Sql = in.ReadInt()
//	m.SqlCrud = in.ReadByte()
//	m.CountTotal = int32(in.ReadDecimal())
//	m.CountError = int32(in.ReadDecimal())
//	//m.CountActived = int32(in.ReadDecimal())
//	m.TimeSum = in.ReadDecimal()
//	m.TimeStd = in.ReadDecimal()
//	m.TimeMin = int32(in.ReadDecimal())
//	m.TimeMax = int32(in.ReadDecimal())
//	m.FetchCount = in.ReadDecimal()
//	m.FetchTime = in.ReadDecimal()
//	m.UpdateCount = in.ReadDecimal()
//	return m
//}

//type SqlRec struct {
//	// int32
//	Dbc int32
//	// int32
//	Sql int32
//	// byte
//	SqlCrud byte
//	// int32
//	CountTotal int32
//	// int32
//	CountError int32
//	// int32
//	CountActived int32
//	// int64
//	TimeSum int64
//	// int64
//	TimeStd int64
//	// int32
//	TimeMin int32
//	// int32
//	TimeMax int32
//	// int64
//	FetchCount int64
//	// int64
//	FetchTime int64
//	// int64
//	UpdateCount int64
//	// []byte
//	Records []byte
//	// int32
//	RecordCount int32
//}
//
//func NewSqlRec() *SqlRec {
//	p := new(SqlRec)
//
//	return p
//}
//
//// func (this *StatSqlPack) Merge(o SqlRec)
//func (this *SqlRec) Merge(o SqlRec) {
//	this.CountTotal += o.CountTotal
//	this.CountError += o.CountError
//	this.CountActived += o.CountActived
//	this.TimeSum += o.TimeSum
//	this.TimeStd = o.TimeStd
//	this.TimeMax = int32(math.Max(float64(this.TimeMax), float64(o.TimeMax)))
//	this.TimeMin = int32(math.Min(float64(this.TimeMin), float64(o.TimeMin)))
//	this.FetchCount += o.FetchCount
//	this.FetchTime += o.FetchTime
//	this.UpdateCount += o.UpdateCount
//}
//
//// func (this *StatSqlPack) Key() int64 {
//func (this *SqlRec) Key() int64 {
//	return bitutil.Composite64(this.Dbc, this.Sql)
//}
//
//// func (this *StatSqlPack) SetDbcSql(dbc, sql int32)  SqlRec {
//func (this *SqlRec) SetDbcSql(dbc, sql int32) *SqlRec {
//	this.Dbc = dbc
//	this.Sql = sql
//	return this
//}
