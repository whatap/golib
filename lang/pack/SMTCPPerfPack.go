package pack

import (
	"github.com/whatap/golib/io"
)

type TCPPortPerf struct {
	Port    int32
	IsAlive bool
}

func (this *TCPPortPerf) Write(out *io.DataOutputX) {
	dout := io.NewDataOutputX()
	dout.WriteInt(this.Port)
	dout.WriteBool(this.IsAlive)

	out.WriteBlob(dout.ToByteArray())
}

func (this *TCPPortPerf) Read(in *io.DataInputX) {
	din := io.NewDataInputX(in.ReadBlob())
	this.Port = din.ReadInt()
	this.IsAlive = din.ReadBool()
}

type SMTCPPerfPack struct {
	AbstractPack
	TCPPortPerf []TCPPortPerf
}

func NewSMTCPPerfPack() *SMTCPPerfPack {
	p := new(SMTCPPerfPack)
	return p
}

func (this *SMTCPPerfPack) GetPackType() int16 {
	return PACK_SM_PORT_PERF
}

func (this *SMTCPPerfPack) Write(dout *io.DataOutputX) {
	this.AbstractPack.Write(dout)
	dout.WriteDecimal(int64(len(this.TCPPortPerf)))
	for _, tcpPortPerf := range this.TCPPortPerf {
		tcpPortPerf.Write(dout)
	}
}

func (this *SMTCPPerfPack) Read(din *io.DataInputX) {
	this.AbstractPack.Read(din)
	tcpCount := din.ReadDecimal()
	if tcpCount > 0 {
		this.TCPPortPerf = make([]TCPPortPerf, tcpCount)
		for i := int64(0); i < tcpCount; i++ {
			this.TCPPortPerf[i].Read(din)
		}
	}
}
