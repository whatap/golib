package udp

import (
	"fmt"

	"github.com/whatap/golib/io"
	"github.com/whatap/golib/util/paramtext"
)

type UdpTxDbcPack struct {
	AbstractPack
	Dbc string
	//error
	ErrorType    string
	ErrorMessage string

	Stack string
}

func NewUdpTxDbcPack() *UdpTxDbcPack {
	p := new(UdpTxDbcPack)
	p.Ver = UDP_PACK_VERSION
	p.AbstractPack.Flush = false
	return p
}

func NewUdpTxDbcPackVer(ver int32) *UdpTxDbcPack {
	p := new(UdpTxDbcPack)
	p.Ver = ver
	p.AbstractPack.Flush = false
	return p
}

func (this *UdpTxDbcPack) GetPackType() uint8 {
	return TX_DB_CONN
}

func (this *UdpTxDbcPack) ToString() string {
	return fmt.Sprint(this.AbstractPack.ToString(), ",dbc=", this.Dbc)
}

func (this *UdpTxDbcPack) Clear() {
	this.AbstractPack.Clear()
	this.AbstractPack.Flush = false

	this.Dbc = ""
	//error
	this.ErrorType = ""
	this.ErrorMessage = ""

	this.Stack = ""
}
func (this *UdpTxDbcPack) Write(dout *io.DataOutputX) {
	this.AbstractPack.Write(dout)
	dout.WriteTextShortLength(this.Dbc)

	dout.WriteTextShortLength(this.ErrorType)
	dout.WriteTextShortLength(this.ErrorMessage)
	dout.WriteTextShortLength(this.Stack)
}

func (this *UdpTxDbcPack) Read(din *io.DataInputX) {
	this.AbstractPack.Read(din)

	this.Dbc = din.ReadTextShortLength()
	this.ErrorType = din.ReadTextShortLength()
	this.ErrorMessage = din.ReadTextShortLength()
	this.Stack = din.ReadTextShortLength()
}
func (this *UdpTxDbcPack) Process() {
	if this.Dbc != "" {
		p := paramtext.NewParamKVSeperate(this.Dbc, " ", "=")
		this.Dbc = p.ToStringStr("password", "#")
		p = paramtext.NewParamKVSeperate(this.Dbc, ";", "=")
		this.Dbc = p.ToStringStr("password", "#")
	}
}
