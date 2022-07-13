package udp

import (
	"fmt"

	"github.com/whatap/golib/io"
	"github.com/whatap/golib/util/stringutil"
	"github.com/whatap/golib/util/urlutil"
)

type UdpTxHttpcPack struct {
	AbstractPack
	// Pack
	Url          string
	Mcallee      int64
	ErrorType    string
	ErrorMessage string
	Stack        string

	//Processing data
	HttpcURL *urlutil.URL
}

func NewUdpTxHttpcPack() *UdpTxHttpcPack {
	p := new(UdpTxHttpcPack)
	p.Ver = UDP_PACK_VERSION
	p.AbstractPack.Flush = false
	return p
}
func NewUdpTxHttpcPackVer(ver int32) *UdpTxHttpcPack {
	p := new(UdpTxHttpcPack)
	p.Ver = ver
	p.AbstractPack.Flush = false
	return p
}
func (this *UdpTxHttpcPack) GetPackType() uint8 {
	return TX_HTTPC
}

func (this *UdpTxHttpcPack) ToString() string {
	return fmt.Sprint(this.AbstractPack.ToString(), ",Url=", this.Url, ",Callee=", this.Mcallee)
}

func (this *UdpTxHttpcPack) Clear() {
	this.AbstractPack.Clear()
	this.AbstractPack.Flush = false

	this.Url = ""
	this.Mcallee = 0
	this.ErrorType = ""
	this.ErrorMessage = ""
	this.Stack = ""

	this.HttpcURL = nil

}

func (this *UdpTxHttpcPack) Write(dout *io.DataOutputX) {
	this.AbstractPack.Write(dout)
	dout.WriteTextShortLength(this.Url)
	dout.WriteTextShortLength(stringutil.ParseStringZeroToEmpty(this.Mcallee))
	dout.WriteTextShortLength(this.ErrorType)
	dout.WriteTextShortLength(this.ErrorMessage)
	dout.WriteTextShortLength(this.Stack)
}

func (this *UdpTxHttpcPack) Read(din *io.DataInputX) {
	this.AbstractPack.Read(din)

	this.Url = din.ReadTextShortLength()

	this.Mcallee = stringutil.ParseInt64(din.ReadTextShortLength())
	this.ErrorType = din.ReadTextShortLength()
	this.ErrorMessage = din.ReadTextShortLength()
	this.Stack = din.ReadTextShortLength()
}

func (this *UdpTxHttpcPack) Process() {
	this.HttpcURL = urlutil.NewURL(this.Url)
}
