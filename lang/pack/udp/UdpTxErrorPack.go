package udp

import (
	"fmt"

	"github.com/whatap/golib/io"
)

type UdpTxErrorPack struct {
	AbstractPack
	ErrorType    string
	ErrorMessage string
	Stack        string
}

func NewUdpTxErrorPack() *UdpTxErrorPack {
	p := new(UdpTxErrorPack)
	p.Ver = UDP_PACK_VERSION
	p.AbstractPack.Flush = false
	return p
}
func NewUdpTxErrorPackVer(ver int32) *UdpTxErrorPack {
	p := new(UdpTxErrorPack)
	p.Ver = ver
	p.AbstractPack.Flush = false
	return p
}

func (this *UdpTxErrorPack) GetPackType() uint8 {
	return TX_ERROR
}

func (this *UdpTxErrorPack) ToString() string {
	return fmt.Sprint(this.AbstractPack.ToString(), ",ErrorType=", this.ErrorType, ",ErrorMessage=", this.ErrorMessage, ",Stack=", this.Stack)
}

func (this *UdpTxErrorPack) Clear() {
	this.AbstractPack.Clear()
	this.AbstractPack.Flush = false

	this.ErrorType = ""
	this.ErrorMessage = ""
	this.Stack = ""
}

func (this *UdpTxErrorPack) Write(dout *io.DataOutputX) {
	this.AbstractPack.Write(dout)
	dout.WriteTextShortLength(this.ErrorType)
	dout.WriteTextShortLength(this.ErrorMessage)

	if this.Ver > 50000 {
		// Golang
	} else if this.Ver > 40000 {
		// Batch
	} else if this.Ver > 30000 {
		// Dotnet
		dout.WriteTextShortLength(this.Stack)
	} else if this.Ver > 20000 {
		// Python
		dout.WriteTextShortLength(this.Stack)
	} else {
		// PHP
	}
}

func (this *UdpTxErrorPack) Read(din *io.DataInputX) {
	this.AbstractPack.Read(din)

	this.ErrorType = din.ReadTextShortLength()
	this.ErrorMessage = din.ReadTextShortLength()

	if this.Ver > 50000 {
		// Golang
	} else if this.Ver > 40000 {
		// Batch
	} else if this.Ver > 30000 {
		// Dotnet
		this.Stack = din.ReadTextShortLength()
	} else if this.Ver > 20000 {
		// Python
		this.Stack = din.ReadTextShortLength()
	} else {
		// PHP
	}
}
func (this *UdpTxErrorPack) Process() {
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
