package pack

import (
	"github.com/whatap/golib/io"
)

type ProcNetPerf struct {
	IP    int32
	Port  int16
	Count int32
}

func (this *ProcNetPerf) Write(out *io.DataOutputX) {
	dout := io.NewDataOutputX()
	dout.WriteInt(this.IP)
	dout.WriteShort(this.Port)
	dout.WriteInt(this.Count)

	out.WriteBlob(dout.ToByteArray())
}

func (this *ProcNetPerf) Read(in *io.DataInputX) {
	din := io.NewDataInputX(in.ReadBlob())
	this.IP = din.ReadInt()
	this.Port = din.ReadShort()
	this.Count = din.ReadInt()
}

type ProcFilePerf struct {
	FilePath int32
	Size     int64
}

func (this *ProcFilePerf) Write(out *io.DataOutputX) {
	dout := io.NewDataOutputX()
	dout.WriteInt(this.FilePath)
	dout.WriteLong(this.Size)

	out.WriteBlob(dout.ToByteArray())
}

func (this *ProcFilePerf) Read(in *io.DataInputX) {
	din := io.NewDataInputX(in.ReadBlob())
	this.FilePath = din.ReadInt()
	this.Size = din.ReadLong()
}

type ProcPerf struct {
	Ppid          int32
	Pid           int32
	Cpu           float32
	MemoryBytes   int64
	MemoryPercent float32
	ReadBps       float32
	WriteBps      float32
	Cmd1          int32
	Cmd2          int32
	ReadIops      float32
	WriteIops     float32

	User       int32
	State      int32
	CreateTime int64

	Group int64

	Net                 []ProcNetPerf
	File                []ProcFilePerf
	MemoryShared        int64
	OpenFileDescriptors int64
}

func (this *ProcPerf) Write(out *io.DataOutputX) {
	dout := io.NewDataOutputX()
	dout.WriteInt(this.Ppid)
	dout.WriteInt(this.Pid)
	dout.WriteFloat(this.Cpu)
	dout.WriteDecimal(this.MemoryBytes)
	dout.WriteFloat(this.MemoryPercent)
	dout.WriteFloat(this.ReadBps)
	dout.WriteFloat(this.WriteBps)

	dout.WriteInt(this.Cmd1)
	dout.WriteInt(this.Cmd2)

	dout.WriteFloat(this.ReadIops)
	dout.WriteFloat(this.WriteIops)

	dout.WriteInt(this.User)
	dout.WriteInt(this.State)
	dout.WriteLong(this.CreateTime)

	dout.WriteDecimal(this.Group)

	dout.WriteDecimal(int64(len(this.Net)))
	for _, procNetPerf := range this.Net {
		procNetPerf.Write(dout)
	}

	dout.WriteDecimal(int64(len(this.File)))
	for _, procFile := range this.File {
		procFile.Write(dout)
	}
	dout.WriteDecimal(this.MemoryShared)
	dout.WriteDecimal(this.OpenFileDescriptors)

	out.WriteBlob(dout.ToByteArray())
}
func (this *ProcPerf) Read(in *io.DataInputX) {
	din := io.NewDataInputX(in.ReadBlob())
	this.Ppid = din.ReadInt()
	this.Pid = din.ReadInt()
	this.Cpu = din.ReadFloat()
	this.MemoryBytes = din.ReadDecimal()
	this.MemoryPercent = din.ReadFloat()
	this.ReadBps = din.ReadFloat()
	this.WriteBps = din.ReadFloat()
	this.Cmd1 = din.ReadInt()
	this.Cmd2 = din.ReadInt()
	this.ReadIops = din.ReadFloat()
	this.WriteIops = din.ReadFloat()

	this.User = din.ReadInt()
	this.State = din.ReadInt()
	this.CreateTime = din.ReadLong()

	this.Group = din.ReadDecimal()

	netCount := din.ReadDecimal()
	if netCount > 0 {
		this.Net = make([]ProcNetPerf, netCount)
		for i := int64(0); i < netCount; i++ {
			this.Net[i].Read(din)
		}
	}

	fileCount := din.ReadDecimal()
	if fileCount > 0 {
		this.File = make([]ProcFilePerf, fileCount)
		for i := int64(0); i < fileCount; i++ {
			this.File[i].Read(din)
		}
	}
	this.MemoryShared = din.ReadDecimal()
	this.OpenFileDescriptors = din.ReadDecimal()
}

type SMProcPerfPack struct {
	AbstractPack
	OS   int16
	Proc []ProcPerf
}

func NewSMProcPerfPack() *SMProcPerfPack {
	p := new(SMProcPerfPack)
	return p
}

func (this *SMProcPerfPack) GetPackType() int16 {
	return PACK_SM_PROC_PERF
}
func (this *SMProcPerfPack) Write(dout *io.DataOutputX) {
	this.AbstractPack.Write(dout)
	dout.WriteShort(this.OS)
	cnt := len(this.Proc)
	dout.WriteDecimal(int64(cnt))
	for i := 0; i < cnt; i++ {
		this.Proc[i].Write(dout)
	}
}
func (this *SMProcPerfPack) Read(din *io.DataInputX) {
	this.AbstractPack.Read(din)
	this.OS = din.ReadShort()
	cnt := int(din.ReadDecimal())
	this.Proc = make([]ProcPerf, cnt)
	for i := 0; i < cnt; i++ {
		this.Proc[i] = ProcPerf{}
		this.Proc[i].Read(din)
	}
}
