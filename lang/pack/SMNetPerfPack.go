package pack

import (
	"github.com/whatap/golib/io"
)

type NetPerf struct {
	Desc       int32
	IP         []byte
	HwAddr     string
	TrafficIn  float64
	TrafficOut float64
	PacketIn   float64
	PacketOut  float64
	ErrorOut   float64
	ErrorIn    float64
	DroppedOut float64
	DroppedIn  float64
	Count      int32
}

func (this *NetPerf) Write(out *io.DataOutputX) {
	dout := io.NewDataOutputX()
	dout.WriteInt(this.Desc)
	dout.WriteBlob(this.IP)
	dout.WriteText(this.HwAddr)

	dout.WriteDouble(this.TrafficIn)
	dout.WriteDouble(this.TrafficOut)
	dout.WriteDouble(this.PacketIn)
	dout.WriteDouble(this.PacketOut)

	dout.WriteDouble(this.ErrorOut)
	dout.WriteDouble(this.ErrorIn)
	dout.WriteDouble(this.DroppedOut)
	dout.WriteDouble(this.DroppedIn)
	dout.WriteInt(1)

	out.WriteBlob(dout.ToByteArray())

}
func (this *NetPerf) Read(in *io.DataInputX) {
	din := io.NewDataInputX(in.ReadBlob())
	this.Desc = din.ReadInt()
	this.IP = din.ReadBlob()
	this.HwAddr = din.ReadText()
	this.TrafficIn = din.ReadDouble()
	this.TrafficOut = din.ReadDouble()
	this.PacketIn = din.ReadDouble()
	this.PacketOut = din.ReadDouble()
	this.ErrorOut = din.ReadDouble()
	this.ErrorIn = din.ReadDouble()
	this.DroppedOut = din.ReadDouble()
	this.DroppedIn = din.ReadDouble()
	this.Count = din.ReadInt()
}

type SMNetPerfPack struct {
	AbstractPack
	OS  int16
	Net []NetPerf
}

func NewSMNetPerfPack() *SMNetPerfPack {
	p := new(SMNetPerfPack)
	return p
}

func (this *SMNetPerfPack) GetPackType() int16 {
	return PACK_SM_NET_PERF
}
func (this *SMNetPerfPack) Write(dout *io.DataOutputX) {
	this.AbstractPack.Write(dout)
	dout.WriteShort(this.OS)
	cnt := len(this.Net)
	dout.WriteDecimal(int64(cnt))
	for i := 0; i < cnt; i++ {
		this.Net[i].Write(dout)
	}
}
func (this *SMNetPerfPack) Read(din *io.DataInputX) {
	this.AbstractPack.Read(din)
	this.OS = din.ReadShort()
	cnt := int(din.ReadDecimal())
	this.Net = make([]NetPerf, cnt)
	for i := 0; i < cnt; i++ {
		this.Net[i] = NetPerf{}
		this.Net[i].Read(din)
	}
}
