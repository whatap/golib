package udp

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/whatap/golib/io"
	"github.com/whatap/golib/util/stringutil"
	"github.com/whatap/golib/util/urlutil"
)

type UdpTxEndPack struct {
	AbstractPack

	// Pack
	Host    string
	Uri     string
	Mtid    int64
	Mdepth  int32
	Mcaller int64

	McallerTxid    int64
	McallerPcode   int64
	McallerSpec    string
	McallerUrl     string
	McallerPoidKey string

	Status int32

	// Processing data
	ServiceURL     *urlutil.URL
	McallerUrlHash int32
}

func NewUdpTxEndPack() *UdpTxEndPack {
	p := new(UdpTxEndPack)
	p.Ver = UDP_PACK_VERSION
	p.AbstractPack.Flush = true
	return p
}

func NewUdpTxEndPackVer(ver int32) *UdpTxEndPack {
	p := new(UdpTxEndPack)
	p.Ver = ver
	p.AbstractPack.Flush = true
	return p
}

func (this *UdpTxEndPack) GetPackType() uint8 {
	return TX_END
}

func (this *UdpTxEndPack) ToString() string {
	return fmt.Sprint(this.AbstractPack.ToString(), ",host=", this.Host, ",uri=", this.Uri, ",elapsed=", this.Elapsed)
}

func (this *UdpTxEndPack) Clear() {
	this.AbstractPack.Clear()
	this.AbstractPack.Flush = true

	this.Host = ""
	this.Uri = ""
	this.Mtid = 0
	this.Mdepth = 0
	this.Mcaller = 0

	this.McallerTxid = 0
	this.McallerPcode = 0
	this.McallerSpec = ""
	this.McallerUrl = ""
	this.McallerPoidKey = ""

	this.Status = 0

	// Processing data
	this.ServiceURL = nil
	this.McallerUrlHash = 0
}

func (this *UdpTxEndPack) SetMcallerUrlHash(v int32) {
	this.McallerUrlHash = v
	this.McallerUrl = fmt.Sprintf("%d", v)
}
func (this *UdpTxEndPack) Write(dout *io.DataOutputX) {
	this.AbstractPack.Write(dout)

	dout.WriteTextShortLength(this.Host)
	dout.WriteTextShortLength(this.Uri)
	dout.WriteTextShortLength(stringutil.ParseStringZeroToEmpty(int64(this.Mtid)))
	dout.WriteTextShortLength(stringutil.ParseStringZeroToEmpty(int64(this.Mdepth)))
	dout.WriteTextShortLength(stringutil.ParseStringZeroToEmpty(int64(this.McallerTxid)))
	dout.WriteTextShortLength(stringutil.ParseStringZeroToEmpty(int64(this.McallerPcode)))
	dout.WriteTextShortLength(this.McallerSpec)
	dout.WriteTextShortLength(this.McallerUrl)
	dout.WriteTextShortLength(this.McallerPoidKey)
	dout.WriteTextShortLength(stringutil.ParseStringZeroToEmpty(int64(this.Status)))

}

func (this *UdpTxEndPack) Read(din *io.DataInputX) {
	this.AbstractPack.Read(din)
	this.Host = din.ReadTextShortLength()
	this.Uri = din.ReadTextShortLength()
	this.Mtid = stringutil.ParseInt64(din.ReadTextShortLength())
	this.Mdepth = stringutil.ParseInt32(din.ReadTextShortLength())
	this.McallerTxid = stringutil.ParseInt64(din.ReadTextShortLength())
	this.McallerPcode = stringutil.ParseInt64(din.ReadTextShortLength())
	this.McallerSpec = din.ReadTextShortLength()
	this.McallerUrl = din.ReadTextShortLength()
	this.McallerPoidKey = din.ReadTextShortLength()
	this.Status = stringutil.ParseInt32(din.ReadTextShortLength())
}
func (this *UdpTxEndPack) Process() {
	if this.Host != "" && this.Uri != "" {
		if strings.HasPrefix(this.Uri, "/") {
			this.ServiceURL = urlutil.NewURL(this.Host + this.Uri)
		} else {
			this.ServiceURL = urlutil.NewURL(this.Host + "/" + this.Uri)
		}
	}
	if ret, err := strconv.ParseInt(this.McallerUrl, 10, 32); err == nil {
		this.McallerUrlHash = int32(ret)
	}
}
