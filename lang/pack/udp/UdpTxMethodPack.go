package udp

import (
	"fmt"

	"github.com/whatap/golib/io"
)

type UdpTxMethodPack struct {
	AbstractPack
	Method string
	Stack  string
}

func NewUdpTxMethodPack() *UdpTxMethodPack {
	p := new(UdpTxMethodPack)
	p.Ver = UDP_PACK_VERSION
	p.AbstractPack.Flush = false
	return p
}

func NewUdpTxMethodPackVer(ver int32) *UdpTxMethodPack {
	p := new(UdpTxMethodPack)
	p.Ver = UDP_PACK_VERSION
	p.AbstractPack.Flush = false
	return p
}

func (this *UdpTxMethodPack) GetPackType() uint8 {
	return TX_METHOD
}

func (this *UdpTxMethodPack) ToString() string {
	return fmt.Sprint(this.AbstractPack.ToString(), ",method=", this.Method, ",stack=", this.Stack)
}

func (this *UdpTxMethodPack) Clear() {
	this.AbstractPack.Clear()
	this.AbstractPack.Flush = false

	this.Method = ""
	this.Stack = ""
}

func (this *UdpTxMethodPack) Write(dout *io.DataOutputX) {
	this.AbstractPack.Write(dout)
	dout.WriteTextShortLength(this.Method)
	dout.WriteTextShortLength(this.Stack)

	if this.Ver > 60000 {
		// Node.js
	} else if this.Ver > 50000 {
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

func (this *UdpTxMethodPack) Read(din *io.DataInputX) {
	this.AbstractPack.Read(din)

	this.Method = din.ReadTextShortLength()
	this.Stack = din.ReadTextShortLength()
	if this.Ver > 60000 {
		// Node.js
	} else if this.Ver > 50000 {
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
func (this *UdpTxMethodPack) Process() {
	if this.Ver > 60000 {
		// Node.js
	} else if this.Ver > 50000 {
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
