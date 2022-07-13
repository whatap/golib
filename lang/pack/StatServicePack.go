package pack

import (
	"fmt"
	"math"

	"github.com/whatap/golib/io"
	"github.com/whatap/golib/util/hmap"
)

// Deprecated: Use StatTransactionPack instead
type ServiceRec struct {
	Hash            int32
	Profiled        bool
	Count           int32
	Error           int32
	Actived         int32
	TimeSum         int64
	TimeStd         int64
	TimeMin         int32
	TimeMax         int32
	SqlCount        int32
	SqlTime         int64
	SqlFetch        int32
	SqlFetchTime    int64
	SqlUpdateRecord int32
	SqlCommitCount  int32

	SqlSelect int32
	SqlUpdate int32
	SqlDelete int32
	SqlInsert int32
	SqlOthers int32

	HttpcCount int32
	HttpcTime  int64

	MallocSum int64
	CpuSum    int64

	Status200 int32
	Status300 int32
	Status400 int32
	Status500 int32

	SqlMap   *hmap.IntKeyMap
	HttpcMap *hmap.IntKeyMap
}

func NewServiceRec() *ServiceRec {
	p := new(ServiceRec)
	return p
}

func (this *ServiceRec) Merge(o *ServiceRec) {

	this.Profiled = (this.Profiled || o.Profiled)
	this.Count += o.Count
	this.Error += o.Error
	this.Actived += o.Actived
	this.TimeSum += o.TimeSum
	this.TimeStd += o.TimeStd
	// TODO 값 확인 필요
	this.TimeMin = int32(math.Min(float64(this.TimeMin), float64(o.TimeMin)))
	this.TimeMax = int32(math.Max(float64(this.TimeMax), float64(o.TimeMax)))
	this.SqlCount += o.SqlCount
	this.SqlTime += o.SqlTime
	this.SqlFetch += o.SqlFetch
	this.SqlFetchTime += o.SqlFetchTime
	this.SqlUpdateRecord += o.SqlUpdateRecord
	this.SqlCommitCount += o.SqlCommitCount

	this.SqlSelect += o.SqlSelect
	this.SqlUpdate += o.SqlUpdate
	this.SqlDelete += o.SqlDelete
	this.SqlInsert += o.SqlInsert
	this.SqlOthers += o.SqlOthers

	this.HttpcCount += o.HttpcCount
	this.HttpcTime += o.HttpcTime
	this.MallocSum += o.MallocSum
	this.CpuSum += o.CpuSum

	this.Status200 += o.Status200
	this.Status300 += o.Status300
	this.Status400 += o.Status400
	this.Status500 += o.Status500

	if o.SqlMap != nil && o.SqlMap.Size() > 0 {
		if this.SqlMap == nil {
			this.SqlMap = CreateMap(o.SqlMap.Size())

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
			this.HttpcMap = CreateMap(o.HttpcMap.Size())
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

func (this *ServiceRec) SetUrlHash(hash int32) *ServiceRec {
	this.Hash = hash
	return this
}

func (this *ServiceRec) ToString() string {
	return fmt.Sprintln("ServiceRec [hash=", this.Hash, ", profiled=", this.Profiled, ", count=", this.Count, ", error=", this.Error, ", actived=", this.Actived,
		", time_sum=", this.TimeSum, ", time_std=", this.TimeStd, ", time_min=", this.TimeMin, ", time_max=", this.TimeMax, ", sql_count=", this.SqlCount,
		", sql_time=", this.SqlTime, ", sql_fetch=", this.SqlFetch, ", sql_fetch_time=", this.SqlFetchTime, ", sql_update_record=", this.SqlUpdateRecord,
		", sql_commit_count=", this.SqlCommitCount, ", sql_select=", this.SqlSelect, ", sql_update=", this.SqlUpdate, ", sql_delete=", this.SqlDelete,
		", sql_insert=", this.SqlInsert, ", sql_others=", this.SqlOthers, ", httpc_count=", this.HttpcCount, ", httpc_time=", this.HttpcTime,
		", malloc_sum=", this.MallocSum, ", cpu_sum=", this.CpuSum, ", status200=", this.Status200, ", status300=", this.Status300,
		", status400=", this.Status400, ", status500=", this.Status500, ", sqlMap=", this.SqlMap, ", httpcMap=", this.HttpcMap, "]")
}

type StatServicePack struct {
	AbstractPack
	// [] byte
	Records []byte
	// int
	RecordCount int
}

func NewStatServicePack() *StatServicePack {
	p := new(StatServicePack)
	return p
}

func (this *StatServicePack) GetPackType() int16 {
	return PACK_STAT_SERVICE
}

func (this *StatServicePack) ToString() string {
	//	sb.Append(",bytes=" + ArrayUtil.len(records));

	return fmt.Sprintln("StatService ", this.Oid, ",", this.Pcode, ",", this.Time, ",records=", this.RecordCount, ",bytes=", len(this.Records))
}

func (this *StatServicePack) Write(dout *io.DataOutputX) {
	this.AbstractPack.Write(dout)
	dout.WriteBlob(this.Records)
	dout.WriteDecimal(int64(this.RecordCount))
}

func (this *StatServicePack) Read(din *io.DataInputX) {
	this.AbstractPack.Read(din)
	this.Records = din.ReadBlob()
	this.RecordCount = int(din.ReadDecimal())
	//return this
}

func (this *StatServicePack) SetRecords(size int, items hmap.Enumeration) *StatServicePack {
	//public StatServicePack setRecords(int size, Enumeration<ServiceRec> items) {
	//fmt.Println("StatServicePack:SetRecords")
	o := io.NewDataOutputX()
	o.WriteShort(int16(size))
	for i := 0; i < size; i++ {
		//fmt.Println("StatServicePack:SetRecords i=", i)
		this.WriteRec(o, items.NextElement().(*ServiceRec))
	}
	this.Records = o.ToByteArray()
	this.RecordCount = size
	return this
}

func (this *StatServicePack) WriteRec(o *io.DataOutputX, m *ServiceRec) {
	o.WriteInt(m.Hash)
	o.WriteBool(m.Profiled)
	o.WriteDecimal(int64(m.Count))
	o.WriteDecimal(int64(m.Error))
	o.WriteDecimal(int64(m.Actived))
	o.WriteDecimal(m.TimeSum)
	o.WriteDecimal(m.TimeStd)
	o.WriteDecimal(int64(m.TimeMin))
	o.WriteDecimal(int64(m.TimeMax))
	o.WriteDecimal(int64(m.SqlCount))
	o.WriteDecimal(m.SqlTime)
	o.WriteDecimal(int64(m.SqlFetch))
	o.WriteDecimal(m.SqlFetchTime)
	o.WriteDecimal(int64(m.SqlUpdateRecord))
	o.WriteDecimal(int64(m.SqlCommitCount))
	o.WriteDecimal(int64(m.SqlSelect))
	o.WriteDecimal(int64(m.SqlUpdate))
	o.WriteDecimal(int64(m.SqlDelete))
	o.WriteDecimal(int64(m.SqlInsert))
	o.WriteDecimal(int64(m.SqlOthers))

	o.WriteDecimal(int64(m.HttpcCount))
	o.WriteDecimal(m.HttpcTime)
	o.WriteDecimal(m.MallocSum)
	o.WriteDecimal(m.CpuSum)

	o.WriteDecimal(int64(m.Status200))
	o.WriteDecimal(int64(m.Status300))
	o.WriteDecimal(int64(m.Status400))
	o.WriteDecimal(int64(m.Status500))
	/////////////////////////////
	if m.SqlMap == nil {
		o.WriteDecimal(0)
	} else {
		o.WriteDecimal(int64(m.SqlMap.Size()))
		//Enumeration<IntKeyEntry<TimeCount>> en = m.sqlMap.entries();
		en := m.SqlMap.Entries()
		for en.HasMoreElements() {
			ent := en.NextElement().(*hmap.IntKeyEntry)
			o.WriteInt(ent.GetKey())
			o.WriteDecimal(int64(ent.GetValue().(*TimeCount).Count))
			o.WriteDecimal(int64(ent.GetValue().(*TimeCount).Error))
			o.WriteDecimal(ent.GetValue().(*TimeCount).Time)
		}

	}
	if m.HttpcMap == nil {
		o.WriteDecimal(0)
	} else {
		o.WriteDecimal(int64(m.HttpcMap.Size()))
		//Enumeration<IntKeyEntry<TimeCount>> en = m.httpcMap.entries();
		en := m.HttpcMap.Entries()
		for en.HasMoreElements() {
			ent := en.NextElement().(*hmap.IntKeyEntry)
			o.WriteInt(ent.GetKey())
			o.WriteDecimal(int64(ent.GetValue().(*TimeCount).Count))
			o.WriteDecimal(int64(ent.GetValue().(*TimeCount).Error))
			o.WriteDecimal(ent.GetValue().(*TimeCount).Time)
		}
	}
}

// TODO
//func (this * StatServicePack) GetRecords() List<ServiceRec> {
//        // TODO
//        List<ServiceRec> items = new ArrayList<StatServicePack.ServiceRec>()
//        if (this.Records == nil)
//            return nil
//        in = io.NewDataInputX(this.Records)
//        size := in.ReadShort() & 0xffff;
//        for i := 0; i < size; i++ {
//            items.add(this.ReadRec(in))
//        }
//        return items
//    }

func ReadRec(in *io.DataInputX) *ServiceRec {

	m := NewServiceRec()
	m.Hash = in.ReadInt()
	m.Profiled = in.ReadBool()
	m.Count = int32(in.ReadDecimal())
	m.Error = int32(in.ReadDecimal())
	m.Actived = int32(in.ReadDecimal())
	m.TimeSum = in.ReadDecimal()
	m.TimeStd = in.ReadDecimal()
	m.TimeMin = int32(in.ReadDecimal())
	m.TimeMax = int32(in.ReadDecimal())

	m.SqlCount = int32(in.ReadDecimal())
	m.SqlTime = in.ReadDecimal()
	m.SqlFetch = int32(in.ReadDecimal())
	m.SqlFetchTime = in.ReadDecimal()
	m.SqlUpdateRecord = int32(in.ReadDecimal())
	m.SqlCommitCount = int32(in.ReadDecimal())

	m.SqlSelect = int32(in.ReadDecimal())
	m.SqlUpdate = int32(in.ReadDecimal())
	m.SqlDelete = int32(in.ReadDecimal())
	m.SqlInsert = int32(in.ReadDecimal())
	m.SqlOthers = int32(in.ReadDecimal())

	m.HttpcCount = int32(in.ReadDecimal())
	m.HttpcTime = in.ReadDecimal()
	m.MallocSum = in.ReadDecimal()
	m.CpuSum = in.ReadDecimal()

	m.Status200 = int32(in.ReadDecimal())
	m.Status300 = int32(in.ReadDecimal())
	m.Status400 = int32(in.ReadDecimal())
	m.Status500 = int32(in.ReadDecimal())

	sqlcnt := int(in.ReadDecimal())

	if sqlcnt > 0 {
		m.SqlMap = CreateMap(sqlcnt)
		for i := 0; i < sqlcnt; i++ {
			hash := in.ReadInt()
			count := int32(in.ReadDecimal())
			err := int32(in.ReadDecimal())
			time := in.ReadDecimal()
			m.SqlMap.Put(hash, NewTimeCount(count, err, time))
		}
	}
	httpcnt := int(in.ReadDecimal())
	if httpcnt > 0 {
		m.HttpcMap = CreateMap(httpcnt)
		for i := 0; i < httpcnt; i++ {
			hash := in.ReadInt()
			count := int32(in.ReadDecimal())
			err := int32(in.ReadDecimal())
			time := in.ReadDecimal()
			m.HttpcMap.Put(hash, NewTimeCount(count, err, time))
		}
	}
	return m
}

// TODO 추후 이름이 혼동되지 않도록 CreateIntKeyMap 또는 StatServicePack 내부로 넣을 것 ?
// StatServicePack. func CreateMap(cnt int) *hmap.IntKeyMap {
func CreateMap(cnt int) *hmap.IntKeyMap {
	p := hmap.NewIntKeyMap(cnt, 1)

	// TODO create override 사용 안함
	//	        return new IntKeyMap<StatServicePack.TimeCount>(cnt, 1f) {
	//            @Override
	//            protected TimeCount create(int key) {
	//                return this.size() >= 1000 ? null : new TimeCount();
	//            }

	return p
}
