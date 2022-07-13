package pack

import (
	"fmt"
	"math"

	"github.com/whatap/golib/io"
	"github.com/whatap/golib/util/hmap"
	//	"github.com/whatap/golib/util/stringutil"
	//"github.com/whatap/golib/util/logutil"
)

// StatServicePack -> TransactionRec
type TransactionRec struct {
	Hash int32

	Profiled bool
	Count    int32
	Error    int32

	//Actived         int32
	TimeSum int64
	//TimeStd         int64
	//TimeMin         int32
	TimeMax int32

	SqlCount     int32
	SqlTime      int64
	SqlFetch     int32
	SqlFetchTime int64
	//SqlUpdateRecord int32
	//SqlCommitCount  int32

	//	SqlSelect int32
	//	SqlUpdate int32
	//	SqlDelete int32
	//	SqlInsert int32
	//	SqlOthers int32

	HttpcCount int32
	HttpcTime  int64

	MallocSum int64
	CpuSum    int64

	//	Status200 int32
	//	Status300 int32
	//	Status400 int32
	//	Status500 int32

	SqlMap   *hmap.IntKeyMap
	HttpcMap *hmap.IntKeyMap

	// ver3
	ApdexSatisfied int32
	ApdexTolerated int32

	// ver4
	TimeMin int32
	TimeStd int64
}

func NewTransactionRec() *TransactionRec {
	p := new(TransactionRec)
	return p
}

func (this *TransactionRec) MergeForStat(o *TransactionRec) {
	this.Count += o.Count
	this.Error += o.Error
	this.TimeSum += o.TimeSum
	this.TimeMax = int32(math.Max(float64(this.TimeMax), float64(o.TimeMax)))
	this.SqlCount += o.SqlCount
	this.SqlTime += o.SqlTime
	this.SqlFetch += o.SqlFetch
	this.SqlFetchTime += o.SqlFetchTime

	this.HttpcCount += o.HttpcCount
	this.HttpcTime += o.HttpcTime
	this.MallocSum += o.MallocSum
	this.CpuSum += o.CpuSum
}

func (this *TransactionRec) Merge(o *TransactionRec) {

	this.Profiled = (this.Profiled || o.Profiled)

	this.Count += o.Count
	this.Error += o.Error

	//this.Actived += o.Actived

	this.TimeSum += o.TimeSum
	//this.TimeStd += o.TimeStd

	// TODO 값 확인 필요
	//this.TimeMin = int32(math.Min(float64(this.TimeMin), float64(o.TimeMin)))
	this.TimeMax = int32(math.Max(float64(this.TimeMax), float64(o.TimeMax)))
	this.SqlCount += o.SqlCount
	this.SqlTime += o.SqlTime
	this.SqlFetch += o.SqlFetch
	this.SqlFetchTime += o.SqlFetchTime
	//this.SqlUpdateRecord += o.SqlUpdateRecord
	//this.SqlCommitCount += o.SqlCommitCount

	//	this.SqlSelect += o.SqlSelect
	//	this.SqlUpdate += o.SqlUpdate
	//	this.SqlDelete += o.SqlDelete
	//	this.SqlInsert += o.SqlInsert
	//	this.SqlOthers += o.SqlOthers

	this.HttpcCount += o.HttpcCount
	this.HttpcTime += o.HttpcTime
	this.MallocSum += o.MallocSum
	this.CpuSum += o.CpuSum

	//	this.Status200 += o.Status200
	//	this.Status300 += o.Status300
	//	this.Status400 += o.Status400
	//	this.Status500 += o.Status500

	if o.SqlMap != nil && o.SqlMap.Size() > 0 {
		if this.SqlMap == nil {
			//this.SqlMap = this.CreateMap(o.SqlMap.Size())
			this.SqlMap = hmap.NewIntKeyMap(o.SqlMap.Size(), 1)

			en := o.SqlMap.Entries()
			for en.HasMoreElements() {
				e := en.NextElement().(*hmap.IntKeyEntry)
				this.SqlMap.Put(e.GetKey(), e.GetValue().(*TimeCount).Copy())
			}
		} else {
			en := o.SqlMap.Entries()
			for en.HasMoreElements() {
				e := en.NextElement().(*hmap.IntKeyEntry)
				tc := this.SqlMap.Get(e.GetKey()).(*TimeCount)
				if tc == nil {
					this.SqlMap.Put(e.GetKey(), e.GetValue().(*TimeCount).Copy())
				} else {
					tc.Merge(e.GetValue().(*TimeCount))
				}
			}
		}
	}
	if o.HttpcMap != nil && o.HttpcMap.Size() > 0 {
		if this.HttpcMap == nil {
			//this.HttpcMap = this.CreateMap(o.HttpcMap.Size())
			this.HttpcMap = hmap.NewIntKeyMap(o.HttpcMap.Size(), 1)
			en := o.HttpcMap.Entries()
			for en.HasMoreElements() {
				e := en.NextElement().(*hmap.IntKeyEntry)
				this.HttpcMap.Put(e.GetKey(), e.GetValue().(*TimeCount).Copy())
			}
		} else {
			en := o.HttpcMap.Entries()
			for en.HasMoreElements() {
				e := en.NextElement().(*hmap.IntKeyEntry)
				tc := this.HttpcMap.Get(e.GetKey()).(*TimeCount)
				if tc == nil {
					this.HttpcMap.Put(e.GetKey(), e.GetValue().(*TimeCount).Copy())
				} else {
					tc.Merge(e.GetValue().(*TimeCount))
				}
			}
		}
	}
}

func (this *TransactionRec) SetUrlHash(hash int32) *TransactionRec {
	this.Hash = hash
	return this
}

func (this *TransactionRec) ToString() string {
	return fmt.Sprintln("TransactionRec [hash=", this.Hash, ", count=", this.Count, ", error=", this.Error,
		", time_sum=", this.TimeSum, this.TimeMax,
		", sql_count=", this.SqlCount, ", sql_time=", this.SqlTime, ", sql_fetch=", this.SqlFetch, ", sql_fetch_time=", this.SqlFetchTime,
		", httpc_count=", this.HttpcCount, ", httpc_time=", this.HttpcTime,
		", malloc_sum=", this.MallocSum, ", cpu_sum=", this.CpuSum,
		", sqlMap=", this.SqlMap, ", httpcMap=", this.HttpcMap,
		", satisfied=", this.ApdexSatisfied, ", tolerated=", this.ApdexTolerated, ",timeMin=", this.TimeMin, ", timestd=", this.TimeStd, "]")
}

func ReadTransactionRec(in *io.DataInputX) *TransactionRec {
	urlhash := in.ReadInt()
	ver := in.ReadByte() // old에서는 profield=boolean 값이었다. 그래서 이값을 버전값으로 활용한다.
	if ver <= 1 {
		panic("not supported for a old version")
	}

	m := NewTransactionRec()
	m.Hash = urlhash

	//m.Profiled = in.ReadBool()
	m.Count = int32(in.ReadDecimal())
	m.Error = int32(in.ReadDecimal())
	//m.Actived = int32(in.ReadDecimal())
	m.TimeSum = in.ReadDecimal()
	//m.TimeStd = in.ReadDecimal()
	//m.TimeMin = int32(in.ReadDecimal())
	m.TimeMax = int32(in.ReadDecimal())

	m.SqlCount = int32(in.ReadDecimal())
	m.SqlTime = in.ReadDecimal()
	m.SqlFetch = int32(in.ReadDecimal())
	m.SqlFetchTime = in.ReadDecimal()
	//m.SqlUpdateRecord = int32(in.ReadDecimal())
	//m.SqlCommitCount = int32(in.ReadDecimal())

	//	m.SqlSelect = int32(in.ReadDecimal())
	//	m.SqlUpdate = int32(in.ReadDecimal())
	//	m.SqlDelete = int32(in.ReadDecimal())
	//	m.SqlInsert = int32(in.ReadDecimal())
	//	m.SqlOthers = int32(in.ReadDecimal())

	m.HttpcCount = int32(in.ReadDecimal())
	m.HttpcTime = in.ReadDecimal()
	m.MallocSum = in.ReadDecimal()
	m.CpuSum = in.ReadDecimal()

	//	m.Status200 = int32(in.ReadDecimal())
	//	m.Status300 = int32(in.ReadDecimal())
	//	m.Status400 = int32(in.ReadDecimal())
	//	m.Status500 = int32(in.ReadDecimal())

	sqlcnt := int(in.ReadDecimal())

	if sqlcnt > 0 {
		//m.SqlMap = this.CreateMap(sqlcnt)
		m.SqlMap = hmap.NewIntKeyMap(sqlcnt, 1)
		for i := 0; i < sqlcnt; i++ {
			hash := in.ReadInt()
			m.SqlMap.Put(hash, NewTimeCountDefault().Read(in))
		}
	}
	httpcnt := int(in.ReadDecimal())
	if httpcnt > 0 {
		//m.HttpcMap = this.CreateMap(httpcnt)
		m.HttpcMap = hmap.NewIntKeyMap(httpcnt, 1)
		for i := 0; i < httpcnt; i++ {
			hash := in.ReadInt()
			m.HttpcMap.Put(hash, NewTimeCountDefault().Read(in))
		}
	}

	// 2021.06.28 추가 이하
	if ver <= 2 {
		return m
	}

	m.ApdexSatisfied = int32(in.ReadDecimal())
	m.ApdexTolerated = int32(in.ReadDecimal())

	if ver <= 3 {
		return m
	}

	//Java 2021.05.12추가됨
	m.TimeMin = int32(in.ReadDecimal())
	m.TimeStd = in.ReadDecimal()

	return m
}

// 삭제 - java create 함수를 오버라이드 할 필요 없어서 삭제
//func (this *TransactionRec) CreateMap(cnt int) *hmap.IntKeyMap {
//	p := hmap.NewIntKeyMap(cnt, 1)
//	return p
//}

func WriteTransactionRec(o *io.DataOutputX, m *TransactionRec, version byte) {
	o.WriteInt(m.Hash)
	o.WriteByte(version) // 0,1 == old version
	o.WriteDecimal(int64(m.Count))
	o.WriteDecimal(int64(m.Error))

	o.WriteDecimal(m.TimeSum)
	o.WriteDecimal(int64(m.TimeMax))

	o.WriteDecimal(int64(m.SqlCount))
	o.WriteDecimal(m.SqlTime)
	o.WriteDecimal(int64(m.SqlFetch))
	o.WriteDecimal(m.SqlFetchTime)

	o.WriteDecimal(int64(m.HttpcCount))
	o.WriteDecimal(m.HttpcTime)
	o.WriteDecimal(m.MallocSum)
	o.WriteDecimal(m.CpuSum)

	/////////////////////////////
	if m.SqlMap == nil {
		o.WriteDecimal(0)
	} else {
		o.WriteDecimal(int64(m.SqlMap.Size()))
		en := m.SqlMap.Entries()
		for en.HasMoreElements() {
			ent := en.NextElement().(*hmap.IntKeyEntry)
			o.WriteInt(ent.GetKey())
			ent.GetValue().(*TimeCount).Write(o)
		}

	}
	if m.HttpcMap == nil {
		o.WriteDecimal(0)
	} else {
		o.WriteDecimal(int64(m.HttpcMap.Size()))
		en := m.HttpcMap.Entries()
		for en.HasMoreElements() {
			ent := en.NextElement().(*hmap.IntKeyEntry)
			o.WriteInt(ent.GetKey())
			ent.GetValue().(*TimeCount).Write(o)
		}
	}

	// 2021.06.28
	if version <= 2 {
		return
	}
	//ver 3
	o.WriteDecimal(int64(m.ApdexSatisfied))
	o.WriteDecimal(int64(m.ApdexTolerated))

	if version <= 3 {
		return
	}

	//ver 4
	o.WriteDecimal(int64(m.TimeMin))
	o.WriteDecimal(m.TimeStd)

	//logutil.Infoln(">>>>", m)
}
