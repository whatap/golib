package pack

import (
	"github.com/whatap/golib/io"
	"github.com/whatap/golib/lang/value"
)

// ##################################################
// CPU ELEMENT
// ##################################################
type Cpu interface {
	Write(dout *io.DataOutputX)
	Read(din *io.DataInputX)
}
type CpuLinux struct {
	User    float32
	System  float32
	Idle    float32
	Nice    float32
	Irq     float32
	Softirq float32
	Steal   float32
	Iowait  float32

	Load1  float32
	Load5  float32
	Load15 float32
}

func (this *CpuLinux) Write(doutx *io.DataOutputX) {
	dout := io.NewDataOutputX()
	dout.WriteFloat(this.User)
	dout.WriteFloat(this.System)
	dout.WriteFloat(this.Idle)
	dout.WriteFloat(this.Nice)
	dout.WriteFloat(this.Irq)
	dout.WriteFloat(this.Softirq)
	dout.WriteFloat(this.Steal)
	dout.WriteFloat(this.Iowait)
	dout.WriteFloat(this.Load1)
	dout.WriteFloat(this.Load5)
	dout.WriteFloat(this.Load15)

	doutx.WriteBlob(dout.ToByteArray())
}
func (this *CpuLinux) Read(dinx *io.DataInputX) {
	din := io.NewDataInputX(dinx.ReadBlob())
	this.User = din.ReadFloat()
	this.System = din.ReadFloat()
	this.Idle = din.ReadFloat()
	this.Nice = din.ReadFloat()
	this.Irq = din.ReadFloat()
	this.Softirq = din.ReadFloat()
	this.Steal = din.ReadFloat()
	this.Iowait = din.ReadFloat()
	this.Load1 = din.ReadFloat()
	this.Load5 = din.ReadFloat()
	this.Load15 = din.ReadFloat()
}

type CpuWindow struct {
	User   float32
	System float32
	Idle   float32

	ProcessorQueueLength float32
}

func (this *CpuWindow) Write(doutx *io.DataOutputX) {
	dout := io.NewDataOutputX()
	dout.WriteFloat(this.User)
	dout.WriteFloat(this.System)
	dout.WriteFloat(this.Idle)
	dout.WriteFloat(this.ProcessorQueueLength)

	doutx.WriteBlob(dout.ToByteArray())
}
func (this *CpuWindow) Read(dinx *io.DataInputX) {
	din := io.NewDataInputX(dinx.ReadBlob())
	this.User = din.ReadFloat()
	this.System = din.ReadFloat()
	this.Idle = din.ReadFloat()
	this.ProcessorQueueLength = din.ReadFloat()
}

// CpuOSX CpuOSX
type CpuOSX struct {
	User    float32
	System  float32
	Idle    float32
	Nice    float32
	Irq     float32
	Softirq float32
	Steal   float32
	Iowait  float32

	Load1  float32
	Load5  float32
	Load15 float32
}

func (this *CpuOSX) Write(doutx *io.DataOutputX) {
	dout := io.NewDataOutputX()
	dout.WriteFloat(this.User)
	dout.WriteFloat(this.System)
	dout.WriteFloat(this.Idle)
	dout.WriteFloat(this.Nice)
	dout.WriteFloat(this.Irq)
	dout.WriteFloat(this.Softirq)
	dout.WriteFloat(this.Steal)
	dout.WriteFloat(this.Iowait)
	dout.WriteFloat(this.Load1)
	dout.WriteFloat(this.Load5)
	dout.WriteFloat(this.Load15)

	doutx.WriteBlob(dout.ToByteArray())
}
func (this *CpuOSX) Read(dinx *io.DataInputX) {
	din := io.NewDataInputX(dinx.ReadBlob())
	this.User = din.ReadFloat()
	this.System = din.ReadFloat()
	this.Idle = din.ReadFloat()
	this.Nice = din.ReadFloat()
	this.Irq = din.ReadFloat()
	this.Softirq = din.ReadFloat()
	this.Iowait = din.ReadFloat()
	this.Load1 = din.ReadFloat()
	this.Load5 = din.ReadFloat()
	this.Load15 = din.ReadFloat()
}

// ##################################################
// MEMORY ELEMENT
// ##################################################
type Memory interface {
	Write(dout *io.DataOutputX)
	Read(din *io.DataInputX)
}
type MemoryLinux struct {
	Total      int64
	Free       int64
	Cached     int64
	Used       int64
	Pused      float32
	Available  int64
	Pavailable float32

	Buffers int64
	Shared  int64

	SwapUsed  int64
	SwapPused float32
	SwapTotal int64

	PageFault float32

	Slab         int64
	SReclaimable int64
	SUnreclaim   int64
}

func (this *MemoryLinux) Write(doutx *io.DataOutputX) {
	dout := io.NewDataOutputX()
	dout.WriteDecimal(this.Total)
	dout.WriteDecimal(this.Free)
	dout.WriteDecimal(this.Cached)
	dout.WriteDecimal(this.Used)
	dout.WriteFloat(this.Pused)
	dout.WriteDecimal(this.Available)
	dout.WriteFloat(this.Pavailable)

	dout.WriteDecimal(this.Buffers)
	dout.WriteDecimal(this.Shared)

	dout.WriteDecimal(this.SwapUsed)
	dout.WriteFloat(this.SwapPused)
	dout.WriteDecimal(this.SwapTotal)

	dout.WriteFloat(this.PageFault)

	dout.WriteDecimal(this.Slab)
	dout.WriteDecimal(this.SReclaimable)
	dout.WriteDecimal(this.SUnreclaim)

	doutx.WriteBlob(dout.ToByteArray())
}
func (this *MemoryLinux) Read(dinx *io.DataInputX) {
	din := io.NewDataInputX(dinx.ReadBlob())
	this.Total = din.ReadDecimal()
	this.Free = din.ReadDecimal()
	this.Cached = din.ReadDecimal()
	this.Used = din.ReadDecimal()
	this.Pused = din.ReadFloat()
	this.Available = din.ReadDecimal()
	this.Pavailable = din.ReadFloat()

	this.Buffers = din.ReadDecimal()
	this.Shared = din.ReadDecimal()

	this.SwapUsed = din.ReadDecimal()
	this.SwapPused = din.ReadFloat()
	this.SwapTotal = din.ReadDecimal()

	this.PageFault = din.ReadFloat()

	this.Slab = din.ReadDecimal()
	this.SReclaimable = din.ReadDecimal()
	this.SUnreclaim = din.ReadDecimal()
}

// MemoryWindow MemoryWindow
type MemoryWindow struct {
	Total             int64
	Free              int64
	Cached            int64
	Used              int64
	Pused             float32
	Available         int64
	Pavailable        float32
	PageFault         float32
	SwapUsed          int64
	SwapPused         float32
	SwapTotal         int64
	PoolPagedBytes    int64
	PoolNonpagedBytes int64
}

func (this *MemoryWindow) Write(doutx *io.DataOutputX) {
	dout := io.NewDataOutputX()
	dout.WriteDecimal(this.Total)
	dout.WriteDecimal(this.Free)
	dout.WriteDecimal(this.Cached)
	dout.WriteDecimal(this.Used)
	dout.WriteFloat(this.Pused)
	dout.WriteDecimal(this.Available)
	dout.WriteFloat(this.Pavailable)
	dout.WriteFloat(this.PageFault)
	dout.WriteDecimal(this.SwapUsed)
	dout.WriteFloat(this.SwapPused)
	dout.WriteDecimal(this.SwapTotal)
	dout.WriteDecimal(this.PoolPagedBytes)
	dout.WriteDecimal(this.PoolNonpagedBytes)

	doutx.WriteBlob(dout.ToByteArray())
}
func (this *MemoryWindow) Read(dinx *io.DataInputX) {
	din := io.NewDataInputX(dinx.ReadBlob())
	this.Total = din.ReadDecimal()
	this.Free = din.ReadDecimal()
	this.Cached = din.ReadDecimal()
	this.Used = din.ReadDecimal()
	this.Pused = din.ReadFloat()
	this.Available = din.ReadDecimal()
	this.Pavailable = din.ReadFloat()
	this.PageFault = din.ReadFloat()
	this.SwapUsed = din.ReadDecimal()
	this.SwapPused = din.ReadFloat()
	this.SwapTotal = din.ReadDecimal()
	this.PoolPagedBytes = din.ReadDecimal()
	this.PoolNonpagedBytes = din.ReadDecimal()
}

// ##################################################
// SYSBASE PACK
// ##################################################
type SMBasePack struct {
	AbstractPack
	IP        int32
	OS        int16
	Cpu       Cpu
	CpuCore   []Cpu
	Memory    Memory
	UpTime    int64
	EpochTime int64
	Extra     *value.MapValue
}

func NewSMBasePack() *SMBasePack {
	p := new(SMBasePack)
	return p
}

func (this *SMBasePack) GetPackType() int16 {
	return PACK_SM_BASE
}
func (this *SMBasePack) Write(doutx *io.DataOutputX) {
	dout := io.NewDataOutputX()
	this.AbstractPack.Write(dout)
	dout.WriteInt(this.IP)
	dout.WriteShort(this.OS)
	this.Cpu.Write(dout)
	dout.WriteByte(byte(len(this.CpuCore)))
	for i := 0; i < len(this.CpuCore); i++ {
		this.CpuCore[i].Write(dout)
	}
	this.Memory.Write(dout)
	dout.WriteDecimal(this.UpTime)
	dout.WriteLong(this.EpochTime)
	if this.Extra != nil && this.Extra.Size() > 0 {
		dout.WriteByte(1)
		value.WriteMapValue(dout, this.Extra)
	} else {
		dout.WriteByte(0)
	}

	doutx.WriteBlob(dout.ToByteArray())
}
func (this *SMBasePack) Read(dinx *io.DataInputX) {
	din := io.NewDataInputX(dinx.ReadBlob())
	this.AbstractPack.Read(din)
	this.IP = din.ReadInt()
	this.OS = din.ReadShort()
	switch this.OS {
	case OS_LINUX, OS_OSX, OS_AIX, OS_HPUX:
		this.Cpu = &CpuLinux{}
		this.Cpu.Read(din)
		cnt := int(din.ReadByte())
		this.CpuCore = make([]Cpu, cnt)
		for i := 0; i < cnt; i++ {
			this.CpuCore[i] = &CpuLinux{}
			this.CpuCore[i].Read(din)
		}

		this.Memory = &MemoryLinux{}
		this.Memory.Read(din)
	case OS_WINDOW:
		this.Cpu = &CpuWindow{}
		this.Cpu.Read(din)
		cnt := int(din.ReadByte())
		this.CpuCore = make([]Cpu, cnt)
		for i := 0; i < cnt; i++ {
			this.CpuCore[i] = &CpuWindow{}
			this.CpuCore[i].Read(din)
		}

		this.Memory = &MemoryWindow{}
		this.Memory.Read(din)
	}
	this.UpTime = din.ReadDecimal()
	this.EpochTime = din.ReadLong()

	if din.Available() == 0 {
		return
	}

	if din.ReadByte() > 0 {
		this.Extra = value.ReadMapValue(din)
	}
}
