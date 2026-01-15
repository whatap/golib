package udp

import (
	"fmt"

	"github.com/whatap/golib/io"
	"github.com/whatap/golib/util/stringutil"
)

type UdpTxResultSetPack struct {
	AbstractPack
	Dbc   string
	Sql   string
	Fetch int32
}

func NewUdpTxResultSetPack() *UdpTxResultSetPack {
	p := new(UdpTxResultSetPack)
	p.Ver = UDP_PACK_VERSION
	p.AbstractPack.Flush = false
	return p
}

func NewUdpTxResultSetPackVer(ver int32) *UdpTxResultSetPack {
	p := new(UdpTxResultSetPack)
	p.Ver = ver
	p.AbstractPack.Flush = false
	return p
}

func (this *UdpTxResultSetPack) GetPackType() uint8 {
	return TX_RESULT_SET
}

func (this *UdpTxResultSetPack) ToString() string {
	return fmt.Sprint(this.AbstractPack.ToString(), ",dbc=", this.Dbc, ",sql=", this.Sql, ",fetch=", this.Fetch)
}

func (this *UdpTxResultSetPack) Clear() {
	this.AbstractPack.Clear()
	this.AbstractPack.Flush = false

	this.Dbc = ""
	this.Sql = ""
	this.Fetch = 0
}

func (this *UdpTxResultSetPack) Write(dout *io.DataOutputX) {
	this.AbstractPack.Write(dout)
	dout.WriteTextShortLength(this.Dbc)
	dout.WriteTextShortLength(this.Sql)
	dout.WriteTextShortLength(stringutil.ParseStringZeroToEmpty(int64(this.Fetch)))

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
		// if this.Ver >= 10105 {
		// }
	}
}

func (this *UdpTxResultSetPack) Read(din *io.DataInputX) {
	this.AbstractPack.Read(din)

	this.Dbc = din.ReadTextShortLength()
	this.Sql = din.ReadTextShortLength()
	this.Fetch = stringutil.ParseInt32(din.ReadTextShortLength())

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

func (this *UdpTxResultSetPack) Process() {

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
