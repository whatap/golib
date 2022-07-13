package pack

import (
	"github.com/whatap/golib/io"
)

type DiskPerf struct {
	DeviceID   int32 /*hash*/
	MountPoint int32
	FileSystem int32

	FreeSpace  int64
	UsedSpace  int64
	TotalSpace int64

	FreePercent float32
	UsedPercent float32

	Blksize     int32
	ReadIops    float64
	WriteIops   float64
	ReadBps     float64
	WriteBps    float64
	IOPercent   float32
	Count       int32
	QueueLength float32

	InodeTotal       int64
	InodeUsed        int64
	InodeUsedPercent float32
	MountOption      int32
}

func (this *DiskPerf) Write(out *io.DataOutputX) {
	dout := io.NewDataOutputX()
	dout.WriteInt(this.DeviceID)
	dout.WriteInt(this.MountPoint)
	dout.WriteInt(this.FileSystem)

	dout.WriteDecimal(this.FreeSpace)
	dout.WriteDecimal(this.UsedSpace)
	dout.WriteDecimal(this.TotalSpace)

	dout.WriteFloat(this.FreePercent)
	dout.WriteFloat(this.UsedPercent)

	dout.WriteInt(this.Blksize)
	dout.WriteDouble(this.ReadIops)
	dout.WriteDouble(this.WriteIops)
	dout.WriteDouble(this.ReadBps)
	dout.WriteDouble(this.WriteBps)
	dout.WriteFloat(this.IOPercent)
	dout.WriteInt(1)
	dout.WriteFloat(this.QueueLength)

	dout.WriteDecimal(this.InodeTotal)
	dout.WriteDecimal(this.InodeUsed)
	dout.WriteFloat(this.InodeUsedPercent)

	dout.WriteInt(this.MountOption)

	out.WriteBlob(dout.ToByteArray())
}
func (this *DiskPerf) Read(in *io.DataInputX) {

	din := io.NewDataInputX(in.ReadBlob())

	this.DeviceID = din.ReadInt() /*hash*/
	this.MountPoint = din.ReadInt()
	this.FileSystem = din.ReadInt()

	this.FreeSpace = din.ReadDecimal()
	this.UsedSpace = din.ReadDecimal()
	this.TotalSpace = din.ReadDecimal()

	this.FreePercent = din.ReadFloat()
	this.UsedPercent = din.ReadFloat()

	this.Blksize = din.ReadInt()
	this.ReadIops = din.ReadDouble()
	this.WriteIops = din.ReadDouble()
	this.ReadBps = din.ReadDouble()
	this.WriteBps = din.ReadDouble()
	this.IOPercent = din.ReadFloat()
	this.Count = din.ReadInt()
	this.QueueLength = din.ReadFloat()

	this.InodeTotal = din.ReadDecimal()
	this.InodeUsed = din.ReadDecimal()
	this.InodeUsedPercent = din.ReadFloat()

	this.MountOption = din.ReadInt()
}

type SMDiskPerfPack struct {
	AbstractPack
	OS   int16
	Disk []DiskPerf
}

func NewSMDiskPerfPack() *SMDiskPerfPack {
	p := new(SMDiskPerfPack)
	return p
}

func (this *SMDiskPerfPack) GetPackType() int16 {
	return PACK_SM_DISK_QUATA
}
func (this *SMDiskPerfPack) Write(dout *io.DataOutputX) {
	this.AbstractPack.Write(dout)
	dout.WriteShort(this.OS)
	cnt := len(this.Disk)
	dout.WriteDecimal(int64(cnt))
	for i := 0; i < cnt; i++ {
		this.Disk[i].Write(dout)
	}
}
func (this *SMDiskPerfPack) Read(din *io.DataInputX) {
	this.AbstractPack.Read(din)
	this.OS = din.ReadShort()
	cnt := int(din.ReadDecimal())
	this.Disk = make([]DiskPerf, cnt)
	for i := 0; i < cnt; i++ {
		this.Disk[i] = DiskPerf{}
		this.Disk[i].Read(din)
	}
}
