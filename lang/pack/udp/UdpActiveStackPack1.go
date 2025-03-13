package udp

import (
	"fmt"

	"github.com/whatap/golib/io"
)

type UdpActiveStackPack1 struct {
	AbstractPack
	Stack string
}

func NewUdpActiveStackPack1() *UdpActiveStackPack1 {
	p := new(UdpActiveStackPack1)
	p.Ver = UDP_PACK_VERSION
	p.AbstractPack.Flush = false
	return p
}
func NewUdpActiveStackPack1Ver(ver int32) *UdpActiveStackPack1 {
	p := new(UdpActiveStackPack1)
	p.Ver = ver
	p.AbstractPack.Flush = false
	return p
}

func (this *UdpActiveStackPack1) GetPackType() uint8 {
	return ACTIVE_STACK_1
}

func (this *UdpActiveStackPack1) ToString() string {
	return fmt.Sprint(this.AbstractPack.ToString(), ",data=", this.Stack)
}

func (this *UdpActiveStackPack1) Clear() {
	this.AbstractPack.Clear()
	this.AbstractPack.Flush = false
	this.Stack = ""
}

func (this *UdpActiveStackPack1) Write(dout *io.DataOutputX) {
	this.AbstractPack.Write(dout)
	dout.WriteTextShortLength(this.Stack)
}

func (this *UdpActiveStackPack1) Read(din *io.DataInputX) {
	this.AbstractPack.Read(din)
	this.Stack = din.ReadTextShortLength()
}

func (this *UdpActiveStackPack1) Process() {
	if this.Ver > 50000 {
		// Golang
	} else if this.Ver > 40000 {
		// Batch
	} else if this.Ver > 30000 {
		// Dotnet
	} else if this.Ver > 20000 {
		// Python
	} else {
		// PHP
	}
}
