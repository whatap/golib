package udp

import (
	"fmt"

	"github.com/whatap/golib/io"
	"github.com/whatap/golib/util/dateutil"
)

const ()

type AbstractPack struct {
	Ver      int32
	Txid     int64
	Time     int64
	Elapsed  int32
	Cpu      int64
	Mem      int64
	Pid      int32
	ThreadId int64

	Flush bool
}

func (this *AbstractPack) Clear() {
	this.Ver = UDP_PACK_VERSION
	this.Txid = 0
	this.Time = 0
	this.Elapsed = 0
	this.Cpu = 0
	this.Mem = 0
	this.Pid = 0
	this.ThreadId = 0

	this.Flush = false
}

func (this *AbstractPack) Write(dout *io.DataOutputX) {
	dout.WriteLong(this.Txid)
	dout.WriteLong(this.Time)
	dout.WriteInt(this.Elapsed)
	dout.WriteLong(this.Cpu)
	dout.WriteLong(this.Mem)
	dout.WriteInt(this.Pid)
	dout.WriteLong(this.ThreadId)
}
func (this *AbstractPack) Read(din *io.DataInputX) {
	this.Txid = din.ReadLong()
	this.Time = din.ReadLong()
	this.Elapsed = din.ReadInt()
	this.Cpu = din.ReadLong()
	this.Mem = din.ReadLong()
	this.Pid = din.ReadInt()
	this.ThreadId = din.ReadLong()
}

// oid 설정   pack interface
func (this *AbstractPack) SetVersion(ver int32) {
	this.Ver = ver
}

// oid 설정   pack interface
func (this *AbstractPack) GetVersion() int32 {
	return this.Ver
}

func (this *AbstractPack) SetFlush(flush bool) {
	this.Flush = flush
}
func (this *AbstractPack) IsFlush() bool {
	return this.Flush
}

func (this *AbstractPack) ToString() string {
	return fmt.Sprint("ver=", this.Ver, " Txid=", this.Txid, ",Time=", dateutil.TimeStamp(this.Time), ",Elapsed=", this.Elapsed,
		",Cpu=", this.Cpu, ",Mem=", this.Mem, ",Pid=", this.Pid, ",Tid=", this.ThreadId)
}
