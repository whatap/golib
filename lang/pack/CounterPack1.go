package pack

import (
	"fmt"
	"os"
	"strconv"

	"github.com/whatap/go-api/config"
	"github.com/whatap/golib/io"
	"github.com/whatap/golib/lang"
	"github.com/whatap/golib/lang/value"
	"github.com/whatap/golib/util/hmap"
)

type CounterPack1 struct {
	AbstractPack

	Version byte

	Duration int32
	Cputime  int64

	HeapTot                 int64
	HeapUse                 int64
	HeapPerm                int64
	HeapPendingFinalization int32
	HeapMax                 int64

	GcCount       int32
	GcTime        int64
	GcOldgenCount int32

	ServiceCount int32
	ServiceError int32
	ServiceTime  int64

	TxDbcTime   float32
	TxSqlTime   float32
	TxHttpcTime float32

	SqlCount int32
	SqlError int32
	SqlTime  int64

	SqlFetchCount int64
	SqlFetchTime  int64

	HttpcCount int32
	HttpcError int32
	HttpcTime  int64

	//active_tx_cont
	ActSvcCount int32
	//active_tx_slice
	ActSvcSlice []int16

	Cpu      float32
	CpuSys   float32
	CpuUsr   float32
	CpuWait  float32
	CpuSteal float32
	CpuIrq   float32

	CpuProc  float32
	CpuCores int32

	Mem  float32
	Swap float32
	Disk float32

	ThreadTotalStarted int64
	ThreadCount        int32
	ThreadDaemon       int32
	ThreadPeakCount    int32

	Starttime   int64
	PackDropped int64

	DbNumActive *hmap.IntIntMap
	DbNumIdle   *hmap.IntIntMap

	Netstat *NETSTAT

	HostIp      int32
	ProcFd      int32
	Tps         float32
	RespTime    int32
	ArrivalRate float32

	// Java = 1, Node = 2, Python = 3, PHP = 4,  URL = 5
	ApType int16

	Websocket *WEBSOCKET

	MacHash int32

	//public IntMapValue extra;
	Extra *value.IntMapValue

	//public int pid;
	Pid int32

	//public final static String[] active_stat_keys = { "method", "sql", "httpc", "dbc", "socket" };
	ActiveStatKeys []string
	//public []int16 active_stat;
	ActiveStat []int16

	//public int32 threadpool_activeCount;
	ThreadPoolActiveCount int32
	//public int32 threadpool_queueSize;
	ThreadPoolQueueSize int32

	//public IntKeyLinkedMap<TxMeter> txcaller_oid_meter;
	TxcallerOidMeter  *hmap.IntKeyLinkedMap
	TxcallerPOidMeter *hmap.LinkedMap
	//public IntKeyLinkedMap<SqlMeter> sql_meter;
	SqlMeter *hmap.IntKeyLinkedMap
	//public IntKeyLinkedMap<HttpcMeter> httpc_meter;
	HttpcMeter *hmap.IntKeyLinkedMap
	//public LinkedMap<PKIND, TxMeter> txcaller_group_meter;
	TxcallerGroupMeter *hmap.LinkedMap
	//public TxMeter txcaller_unknown;
	TxcallerUnknown *TxMeter

	ContainerKey int32

	// AppDex
	ApdexSatisfied int32
	ApdexTolerated int32
	ApdexTotal     int32

	// TODO
	ProcFdMax int32

	// -1 초기화 안됨 , 0 old버전, 0.1이상 정상 수집
	Metering float32

	Resp90     int32
	Resp95     int32
	TimeSqrSum int64

	// java transient , read, write 없음.
	CollectIntervalMs int
}

func NewCounterPack1() *CounterPack1 {
	p := new(CounterPack1)
	p.Starttime, _ = strconv.ParseInt(os.Getenv("WHATAP.starttime"), 10, 64)
	p.ActSvcSlice = make([]int16, 3)
	p.ApType = config.GetConfig().AppType
	p.ActiveStatKeys = []string{"method", "sql", "httpc", "dbc", "socket"}
	p.ActiveStat = make([]int16, 0)
	return p
}

func (this *CounterPack1) GetPackType() int16 {
	return PACK_COUNTER_1
}
func (this *CounterPack1) Write(out *io.DataOutputX) {
	this.AbstractPack.Write(out)

	dout := io.NewDataOutputX()
	dout.WriteDecimal(int64(this.Duration))
	dout.WriteDecimal(int64(this.Cputime))

	dout.WriteDecimal(this.HeapTot)
	dout.WriteDecimal(this.HeapUse)
	dout.WriteDecimal(this.HeapPerm)
	dout.WriteDecimal(int64(this.HeapPendingFinalization))

	dout.WriteDecimal(int64(this.GcCount))
	dout.WriteDecimal(int64(this.GcTime))

	dout.WriteDecimal(int64(this.ServiceCount))
	dout.WriteDecimal(int64(this.ServiceError))
	dout.WriteDecimal(int64(this.ServiceTime))

	dout.WriteDecimal(int64(this.SqlCount))
	dout.WriteDecimal(int64(this.SqlError))
	dout.WriteDecimal(int64(this.SqlTime))
	dout.WriteDecimal(int64(this.SqlFetchCount))
	dout.WriteDecimal(int64(this.SqlFetchTime))

	dout.WriteDecimal(int64(this.HttpcCount))
	dout.WriteDecimal(int64(this.HttpcError))
	dout.WriteDecimal(int64(this.HttpcTime))

	dout.WriteDecimal(int64(this.ActSvcCount))
	this.writeShortArray(dout, this.ActSvcSlice)

	dout.WriteFloat(this.Cpu)
	dout.WriteFloat(this.CpuSys)
	dout.WriteFloat(this.CpuUsr)
	dout.WriteFloat(this.CpuWait)
	dout.WriteFloat(this.CpuSteal)
	dout.WriteFloat(this.CpuIrq)

	dout.WriteFloat(this.CpuProc)
	dout.WriteDecimal(int64(this.CpuCores))
	dout.WriteFloat(this.Mem)
	dout.WriteFloat(this.Swap)
	dout.WriteFloat(this.Disk)

	dout.WriteDecimal(int64(this.ThreadTotalStarted))
	dout.WriteDecimal(int64(this.ThreadCount))
	dout.WriteDecimal(int64(this.ThreadDaemon))
	dout.WriteDecimal(int64(this.ThreadPeakCount))

	if this.DbNumActive == nil || this.DbNumIdle == nil {
		dout.WriteByte(0)
	} else {
		dout.WriteByte(1)
		this.DbNumActive.ToBytes(dout)
		this.DbNumIdle.ToBytes(dout)
	}
	//	dout.WriteByte(0)

	if this.Netstat == nil {
		dout.WriteByte(0)
	} else {
		dout.WriteByte(1)
		dout.WriteDecimal(int64(this.Netstat.Est))
		dout.WriteDecimal(int64(this.Netstat.FinW))
		dout.WriteDecimal(int64(this.Netstat.CloW))
		dout.WriteDecimal(int64(this.Netstat.TimW))
	}
	//	dout.WriteByte(0)

	dout.WriteDecimal(int64(this.ProcFd))
	dout.WriteFloat(this.Tps)
	dout.WriteDecimal(int64(this.RespTime))

	dout.WriteShort(this.ApType)

	//Websocket
	if this.Websocket == nil {
		dout.WriteByte(0)
	} else {
		dout.WriteByte(1)
		dout.WriteDecimal(int64(this.Websocket.Count))
		dout.WriteDecimal(this.Websocket.In)
		dout.WriteDecimal(this.Websocket.Out)
	}
	//dout.WriteByte(0)

	dout.WriteDecimal(int64(this.Starttime))
	dout.WriteDecimal(int64(this.PackDropped))
	dout.WriteDecimal(int64(this.HostIp))
	dout.WriteDecimal(int64(this.MacHash))

	if this.Extra == nil {
		dout.WriteByte(0)
	} else {
		dout.WriteByte(1)
		value.WriteValue(dout, this.Extra)
	}
	dout.WriteInt(this.Pid)
	sz := 0
	if this.ActiveStat != nil {
		sz = len(this.ActiveStat)
	}
	dout.WriteByte(byte(sz))

	for i := 0; i < sz; i++ {
		dout.WriteShort(this.ActiveStat[i])
	}

	dout.WriteDecimal(int64(this.ThreadPoolActiveCount))
	dout.WriteDecimal(int64(this.ThreadPoolQueueSize))

	this.writeTxcallerOidMeter(dout)
	this.writeSqlMeter(dout)
	this.writeHttpcMeter(dout)
	this.writeTxcallerGroupMeter(dout)
	dout.WriteDecimal(0)
	this.writeTxcallerOther(dout)

	// containerKey
	dout.WriteDecimal(int64(this.ContainerKey))

	// tx_dbc_tim
	dout.WriteFloat(this.TxDbcTime)
	// tx_sql_time
	dout.WriteFloat(this.TxSqlTime)
	// tx_httpc_time
	dout.WriteFloat(this.TxHttpcTime)

	// AppDex
	// service_statisfied
	dout.WriteDecimal(int64(this.ApdexSatisfied))
	// service_tolerated
	dout.WriteDecimal(int64(this.ApdexTolerated))

	// arrival_rate
	dout.WriteFloat(this.ArrivalRate)

	// gc_oldgen_count
	dout.WriteDecimal(int64(this.GcOldgenCount))
	// Version
	dout.WriteByte(this.Version)
	// heap_max
	dout.WriteDecimal(this.HeapMax)

	// proc_fd_max
	dout.WriteDecimal(int64(this.ProcFdMax))
	//
	dout.WriteFloat(this.Metering)

	dout.WriteDecimal(int64(this.ApdexTotal))

	this.writeTxcallerPOidMeter(dout)

	//v2.2.1 2022.11.25
	dout.WriteDecimal(int64(this.Resp90))
	dout.WriteDecimal(int64(this.Resp95))
	//
	dout.WriteDecimal(this.TimeSqrSum)

	out.WriteBlob(dout.ToByteArray())
}

func (this *CounterPack1) Read(in *io.DataInputX) {
	this.AbstractPack.Read(in)
	din := io.NewDataInputX(in.ReadBlob())

	this.Duration = int32(din.ReadDecimal())
	this.Cputime = din.ReadDecimal()
	this.HeapTot = din.ReadDecimal()
	this.HeapUse = din.ReadDecimal()
	this.HeapPerm = din.ReadDecimal()
	this.HeapPendingFinalization = int32(din.ReadDecimal())

	this.GcCount = int32(din.ReadDecimal())
	this.GcTime = din.ReadDecimal()

	this.ServiceCount = int32(din.ReadDecimal())
	this.ServiceError = int32(din.ReadDecimal())
	this.ServiceTime = din.ReadDecimal()

	this.SqlCount = int32(din.ReadDecimal())
	this.SqlError = int32(din.ReadDecimal())
	this.SqlTime = din.ReadDecimal()
	this.SqlFetchCount = din.ReadDecimal()
	this.SqlFetchTime = din.ReadDecimal()

	this.HttpcCount = int32(din.ReadDecimal())
	this.HttpcError = int32(din.ReadDecimal())
	this.HttpcTime = din.ReadDecimal()

	this.ActSvcCount = int32(din.ReadDecimal())
	this.ActSvcSlice = this.readShortArray(din)

	this.Cpu = din.ReadFloat()
	this.CpuSys = din.ReadFloat()
	this.CpuUsr = din.ReadFloat()
	this.CpuWait = din.ReadFloat()
	this.CpuSteal = din.ReadFloat()
	this.CpuIrq = din.ReadFloat()

	this.CpuProc = din.ReadFloat()
	this.CpuCores = int32(din.ReadDecimal())

	this.Mem = din.ReadFloat()
	this.Swap = din.ReadFloat()
	this.Disk = din.ReadFloat()

	this.ThreadTotalStarted = din.ReadDecimal()
	this.ThreadCount = int32(din.ReadDecimal())
	this.ThreadDaemon = int32(din.ReadDecimal())
	this.ThreadPeakCount = int32(din.ReadDecimal())

	// db_num_active ... db_num_idle
	if din.ReadByte() != 0 {
		this.ReadDropMap(din)
		//		this.db_num_active = new IntIntMap(7, 1f);
		//		this.db_num_idle = new IntIntMap(7, 1f);
		//		this.db_num_active.toObject(din);
		//		this.db_num_idle.toObject(din);
	}

	// Netstat
	if din.ReadByte() != 0 {
		din.ReadDecimal()
		din.ReadDecimal()
		din.ReadDecimal()
		din.ReadDecimal()

		//		this.netstat = new NETSTAT();
		//		this.netstat.est = (int) din.readDecimal();
		//		this.netstat.fin_w = (int) din.readDecimal();
		//		this.netstat.clo_w = (int) din.readDecimal();
		//		this.netstat.tim_w = (int) din.readDecimal();
	}

	this.ProcFd = int32(din.ReadDecimal())
	this.Tps = din.ReadFloat()
	this.RespTime = int32(din.ReadDecimal())

	this.ApType = din.ReadShort()

	// websocket
	if din.ReadByte() != 0 {
		//this.websocket = new WEBSOCKET();
		// this.websocket.count
		din.ReadDecimal()
		// this.websocket.in
		din.ReadDecimal()
		// this.websocket.out
		din.ReadDecimal()
	}

	this.Starttime = din.ReadDecimal()
	this.PackDropped = din.ReadDecimal()
	this.HostIp = int32(din.ReadDecimal())

	//if (din.available() > 0) {
	this.MacHash = int32(din.ReadDecimal())

	//if (din.available() > 0) {
	if din.ReadByte() == 1 {
		//this.Extra = (IntMapValue) din.readValue();
		this.Extra = value.NewIntMapValue()
		this.Extra.Read(din)
	}

	this.ActiveStat = make([]int16, 0)

	//if (din.available() > 0) {
	this.Pid = din.ReadInt()

	sz := int(din.ReadByte())
	for i := 0; i < sz; i++ {
		this.ActiveStat = append(this.ActiveStat, din.ReadShort())
	}

	//if (din.available() > 0) {
	this.ThreadPoolActiveCount = int32(din.ReadDecimal())
	this.ThreadPoolQueueSize = int32(din.ReadDecimal())

	//if (din.available() > 0) {
	this.readTxcallerOidMeter(din)
	//if (din.available() > 0) {
	this.readSqlMeter(din)
	//if (din.available() > 0) {
	this.readHttpcMeter(din)
	//if (din.available() > 0) {
	this.readTxcallerGroupMeter(din)
	//if (din.available() > 0) {
	this.readTxcallerOkindMeterDeprecated(din)
	//if (din.available() > 0) {
	this.readTxcallerUnknown(din)

	//if (din.available() > 0) {
	this.ContainerKey = int32(din.ReadDecimal())

	//if (din.available() > 0) {
	this.TxDbcTime = din.ReadFloat()
	this.TxSqlTime = din.ReadFloat()
	this.TxHttpcTime = din.ReadFloat()

	//	if (din.available() > 0) {
	this.ApdexSatisfied = int32(din.ReadDecimal())
	this.ApdexTolerated = int32(din.ReadDecimal())

	// if (din.available() > 0) {
	this.ArrivalRate = din.ReadFloat()

	//		if (din.available() > 0) {
	this.GcOldgenCount = int32(din.ReadDecimal())
	this.Version = din.ReadByte()
	this.HeapMax = din.ReadDecimal()
	//		}
	//		if (din.available() > 0) {
	this.ProcFdMax = int32(din.ReadDecimal())
	//		}
	//		if (din.available() > 0) {
	this.Metering = din.ReadFloat()
	//		}
	//		if (din.available() > 0) {
	this.ApdexTotal = int32(din.ReadDecimal())
	//		}

	this.readTxcallerPOidMeter(din)

	this.Resp90 = int32(din.ReadDecimal())
	this.Resp95 = int32(din.ReadDecimal())

	this.TimeSqrSum = din.ReadDecimal()

}

func (this *CounterPack1) writeShortArray(out *io.DataOutputX, v []int16) {
	if v == nil {
		out.WriteByte(0)
	} else {
		out.WriteByte(byte(len(v)))
		for i := 0; i < len(v); i++ {
			out.WriteShort(v[i])
		}
	}
}
func (this *CounterPack1) readShortArray(in *io.DataInputX) []int16 {
	sz := int(in.ReadByte())
	out := make([]int16, sz)
	for i := 0; i < sz; i++ {
		out[i] = in.ReadShort()
	}
	return out
}
func (this *CounterPack1) ReadDropMap(in *io.DataInputX) {
	cnt := int(in.ReadDecimal())
	for i := 0; i < cnt; i++ {
		in.ReadDecimal()
		in.ReadDecimal()
	}
}

func (this *CounterPack1) readTxcallerUnknown(din *io.DataInputX) {
	ver := din.ReadByte()
	if ver > 0 {
		this.TxcallerUnknown = new(TxMeter)
		this.TxcallerUnknown.Time = din.ReadDecimal()
		this.TxcallerUnknown.Count = int32(din.ReadDecimal())
		this.TxcallerUnknown.Error = int32(din.ReadDecimal())
		if ver >= 2 {
			this.TxcallerUnknown.Actx = int32(din.ReadDecimal())
		}
	}
}

func (this *CounterPack1) readTxcallerGroupMeter(din *io.DataInputX) {
	ver := int(din.ReadByte())
	if ver == 0 {
		return
	}
	count := 0
	if ver <= 8 {
		count = int(din.ReadDecimalLen(ver))
	} else {
		count = int(din.ReadDecimal())
	}
	this.TxcallerGroupMeter = hmap.NewLinkedMapDefault()

	for i := 0; i < count; i++ {
		m := new(TxMeter)
		pcode := din.ReadDecimal()
		if ver <= 8 {
			m.Time = din.ReadDecimal()
			m.Count = int32(din.ReadDecimal())
			m.Error = int32(din.ReadDecimal())
			this.TxcallerGroupMeter.Put(lang.NewPKIND(pcode, int32(0)), m)
		} else {
			okind := int32(din.ReadDecimal())
			m.Time = din.ReadDecimal()
			m.Count = int32(din.ReadDecimal())
			m.Error = int32(din.ReadDecimal())
			m.Actx = int32(din.ReadDecimal())
			this.TxcallerGroupMeter.Put(lang.NewPKIND(pcode, okind), m)
		}
	}
}

func (this *CounterPack1) readTxcallerPOidMeter(din *io.DataInputX) {
	ver := int(din.ReadByte())
	if ver == 0 {
		return
	}
	count := 0
	if ver <= 8 {
		count = int(din.ReadDecimalLen(ver))
	} else {
		count = int(din.ReadDecimal())
	}
	this.TxcallerPOidMeter = hmap.NewLinkedMapDefault()

	for i := 0; i < count; i++ {
		m := new(TxMeter)
		pcode := din.ReadDecimal()
		oid := int32(din.ReadDecimal())
		m.Time = din.ReadDecimal()
		m.Count = int32(din.ReadDecimal())
		m.Error = int32(din.ReadDecimal())
		m.Acts = ReadShortArray(din, count)
		m.Actx = int32(din.ReadDecimal())
		this.TxcallerPOidMeter.Put(lang.NewPOID(pcode, oid), m)

	}
}
func (this *CounterPack1) readTxcallerOkindMeterDeprecated(din *io.DataInputX) {
	count := int(din.ReadDecimal())
	for i := 0; i < count; i++ {
		m := new(TxMeter)

		//okind := din.ReadInt()
		din.ReadInt()
		m.Time = din.ReadDecimal()
		m.Count = int32(din.ReadDecimal())
		m.Error = int32(din.ReadDecimal())
	}

}

func (this *CounterPack1) readHttpcMeter(din *io.DataInputX) {
	ver := int(din.ReadByte())
	if ver == 0 {
		return
	}

	count := 0
	if ver <= 8 {
		count = int(din.ReadDecimalLen(ver))
	} else {
		count = int(din.ReadDecimal())
	}
	//this.httpc_meter = new IntKeyLinkedMap<CounterPack1.HttpcMeter>();
	this.HttpcMeter = hmap.NewIntKeyLinkedMapDefault()

	for i := 0; i < count; i++ {
		m := new(HttpcMeter)
		host := din.ReadInt()
		m.Time = din.ReadDecimal()
		m.Count = int32(din.ReadDecimal())
		m.Error = int32(din.ReadDecimal())
		if ver >= 9 {
			m.Actx = int32(din.ReadDecimal())
		}
		this.HttpcMeter.Put(host, m)
	}

}

func (this *CounterPack1) readSqlMeter(din *io.DataInputX) {
	ver := int(din.ReadByte())
	if ver == 0 {
		return
	}
	count := 0
	if ver <= 8 {
		count = int(din.ReadDecimalLen(ver))
	} else {
		count = int(din.ReadDecimal())
	}
	//this.sql_meter = new IntKeyLinkedMap<CounterPack1.SqlMeter>();
	this.SqlMeter = hmap.NewIntKeyLinkedMapDefault()

	for i := 0; i < count; i++ {
		m := new(SqlMeter)
		dbc := din.ReadInt()
		m.Time = din.ReadDecimal()
		m.Count = int32(din.ReadDecimal())
		m.Error = int32(din.ReadDecimal())
		if ver >= 9 {
			m.Actx = int32(din.ReadDecimal())
		}
		m.FetchCount = din.ReadDecimal()
		m.FetchTime = din.ReadDecimal()
		this.SqlMeter.Put(dbc, m)
	}

}

func (this *CounterPack1) readTxcallerOidMeter(din *io.DataInputX) {

	ver := int(din.ReadByte())
	if ver == 0 {
		return
	}
	count := 0
	if ver <= 8 {
		count = int(din.ReadDecimalLen(ver))
	} else {
		count = int(din.ReadDecimal())
	}

	//this.txcaller_oid_meter = new IntKeyLinkedMap<CounterPack1.TxMeter>();
	this.TxcallerOidMeter = hmap.NewIntKeyLinkedMapDefault()

	for i := 0; i < count; i++ {
		m := new(TxMeter)
		key := din.ReadInt()
		m.Time = din.ReadDecimal()
		m.Count = int32(din.ReadDecimal())
		m.Error = int32(din.ReadDecimal())
		if ver >= 9 {
			m.Actx = int32(din.ReadDecimal())
		}
		this.TxcallerOidMeter.Put(key, m)
	}

}

func (this *CounterPack1) dbcCount() int {
	dbc := 0
	if this.DbNumActive != nil {
		en := this.DbNumActive.Values()
		for en.HasMoreElements() {
			dbc += int(en.NextInt())
		}
	}
	if this.DbNumIdle != nil {
		en := this.DbNumIdle.Values()
		for en.HasMoreElements() {
			dbc += int(en.NextInt())
		}
	}
	return dbc
}

func (this *CounterPack1) writeTxcallerOther(dout *io.DataOutputX) {
	if this.TxcallerUnknown != nil {
		dout.WriteByte(2)
		dout.WriteDecimal(this.TxcallerUnknown.Time)
		dout.WriteDecimal(int64(this.TxcallerUnknown.Count))
		dout.WriteDecimal(int64(this.TxcallerUnknown.Error))
		dout.WriteDecimal(int64(this.TxcallerUnknown.Actx))
	} else {
		dout.WriteByte(0)
	}
}

func (this *CounterPack1) writeTxcallerOidMeter(dout *io.DataOutputX) {
	if this.TxcallerOidMeter == nil {
		dout.WriteDecimal(0)
	} else {
		dout.WriteByte(9)
		dout.WriteDecimal(int64(this.TxcallerOidMeter.Size()))
		en := this.TxcallerOidMeter.Entries()
		for en.HasMoreElements() {
			ent := en.NextElement().(*hmap.IntKeyLinkedEntry)
			m := ent.GetValue().(*TxMeter)

			dout.WriteInt(ent.GetKey())
			dout.WriteDecimal(m.Time)
			dout.WriteDecimal(int64(m.Count))
			dout.WriteDecimal(int64(m.Error))
			dout.WriteDecimal(int64(m.Actx))
		}
	}
}

func (this *CounterPack1) writeSqlMeter(dout *io.DataOutputX) {
	if this.SqlMeter == nil {
		dout.WriteDecimal(0)
	} else {
		dout.WriteByte(9)
		dout.WriteDecimal(int64(this.SqlMeter.Size()))
		idx := 0
		en := this.SqlMeter.Entries()
		for en.HasMoreElements() {
			ent := en.NextElement().(*hmap.IntKeyLinkedEntry)
			m := ent.GetValue().(*SqlMeter)
			idx++

			dout.WriteInt(ent.GetKey())
			dout.WriteDecimal(m.Time)
			dout.WriteDecimal(int64(m.Count))
			dout.WriteDecimal(int64(m.Error))
			dout.WriteDecimal(int64(m.Actx))
			dout.WriteDecimal(int64(m.FetchCount))
			dout.WriteDecimal(int64(m.FetchTime))
		}
	}
}

func (this *CounterPack1) writeHttpcMeter(dout *io.DataOutputX) {
	if this.HttpcMeter == nil {
		dout.WriteDecimal(0)
	} else {
		dout.WriteByte(9)
		dout.WriteDecimal(int64(this.HttpcMeter.Size()))
		en := this.HttpcMeter.Entries()
		for en.HasMoreElements() {
			ent := en.NextElement().(*hmap.IntKeyLinkedEntry)
			m := ent.GetValue().(*HttpcMeter)

			dout.WriteInt(ent.GetKey())
			dout.WriteDecimal(m.Time)
			dout.WriteDecimal(int64(m.Count))
			dout.WriteDecimal(int64(m.Error))
			dout.WriteDecimal(int64(m.Actx))
		}
	}
}

func (this *CounterPack1) writeTxcallerGroupMeter(dout *io.DataOutputX) {
	if this.TxcallerGroupMeter == nil {
		dout.WriteDecimal(0)
	} else {
		dout.WriteByte(9)
		dout.WriteDecimal(int64(this.TxcallerGroupMeter.Size()))
		//Enumeration<LinkedEntry<PKIND, TxMeter>> en = this.txcaller_group_meter.entries();
		en := this.TxcallerGroupMeter.Entries()
		for en.HasMoreElements() {
			ent := en.NextElement().(*hmap.LinkedEntry)
			m := ent.GetValue().(*TxMeter)
			dout.WriteDecimal(ent.GetKey().(*lang.PKIND).PCode)
			dout.WriteDecimal(int64(ent.GetKey().(*lang.PKIND).OKind))
			dout.WriteDecimal(m.Time)
			dout.WriteDecimal(int64(m.Count))
			dout.WriteDecimal(int64(m.Error))
			dout.WriteDecimal(int64(m.Actx))
		}
	}
}

func (this *CounterPack1) writeTxcallerPOidMeter(dout *io.DataOutputX) {
	if this.TxcallerPOidMeter == nil {
		dout.WriteDecimal(0)
	} else {
		dout.WriteDecimal(int64(this.TxcallerPOidMeter.Size()))
		en := this.TxcallerPOidMeter.Entries()
		for en.HasMoreElements() {
			ent := en.NextElement().(*hmap.LinkedEntry)
			m := ent.GetValue().(*TxMeter)
			dout.WriteDecimal(ent.GetKey().(*lang.POID).PCode)
			dout.WriteDecimal(int64(ent.GetKey().(*lang.POID).Oid))
			dout.WriteDecimal(m.Time)
			dout.WriteDecimal(int64(m.Count))
			dout.WriteDecimal(int64(m.Error))
			dout.WriteDecimal(int64(m.Actx))
		}
	}
}

type TxMeter struct {
	Time  int64
	Count int32
	Error int32
	Actx  int32
	Acts  []int16
}

func NewTxMeter() *TxMeter {
	return new(TxMeter)
}
func (this *TxMeter) ToString() string {
	return fmt.Sprintln("[time=", this.Time, ", count=", this.Count, ", error=", this.Error, ", actx=", this.Actx, "]")
}

type HttpcMeter struct {
	TxMeter
}

func NewHttpcMeter() *HttpcMeter {
	return new(HttpcMeter)
}
func (this *HttpcMeter) ToString() string {
	return fmt.Sprintln("HttpcMeter ", this.TxMeter.ToString())
}

type SqlMeter struct {
	TxMeter
	FetchCount int64
	FetchTime  int64
}

func NewSqlMeter() *SqlMeter {
	return new(SqlMeter)
}
func (this *SqlMeter) ToString() string {
	return fmt.Sprintln("SqlMeter [fetch_count=", this.FetchCount, ", fetch_time=", this.FetchTime, "]", this.TxMeter.ToString())
}

func ReadShortArray(din *io.DataInputX, sz int) []int16 {
	len := int(din.ReadByte())
	arr := make([]int16, sz) // sz는 최대 사이즈, 0인경우가 있음
	for i := 0; i < len; i++ {
		arr[i] = din.ReadShort()
	}
	return arr
}
