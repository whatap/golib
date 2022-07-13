package udp

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/whatap/golib/io"
	"github.com/whatap/golib/util/stringutil"
	"github.com/whatap/golib/util/urlutil"
)

type UdpTxStartEndPack struct {
	AbstractPack
	// Pack
	Host             string
	Uri              string
	Ipaddr           string
	UAgent           string
	Ref              string
	WClientId        string
	HttpMethod       string
	IsStaticContents string

	Mtid    int64
	Mdepth  int32
	Mcaller int64

	McallerTxid    int64
	McallerPcode   int64
	McallerSpec    string
	McallerUrl     string
	McallerPoidKey string

	//Processing data
	ServiceURL     *urlutil.URL
	RefererURL     *urlutil.URL
	IsStatic       bool
	McallerUrlHash int32
}

func NewUdpTxStartEndPack() *UdpTxStartEndPack {
	p := new(UdpTxStartEndPack)
	p.Ver = UDP_PACK_VERSION
	p.AbstractPack.Flush = true
	return p
}
func NewUdpTxStartEndPackVer(ver int32) *UdpTxStartEndPack {
	p := new(UdpTxStartEndPack)
	p.Ver = ver
	p.AbstractPack.Flush = true
	return p
}

func (this *UdpTxStartEndPack) GetPackType() uint8 {
	return TX_START_END
}

func (this *UdpTxStartEndPack) ToString() string {
	return fmt.Sprint(this.AbstractPack.ToString(), ",host=", this.Host, ",uri=", this.Uri)
}

func (this *UdpTxStartEndPack) Clear() {
	this.AbstractPack.Clear()
	this.AbstractPack.Flush = true

	// Pack
	this.Host = ""
	this.Uri = ""
	this.Ipaddr = ""
	this.UAgent = ""
	this.Ref = ""
	this.WClientId = ""
	this.HttpMethod = ""
	this.IsStaticContents = ""

	this.Mtid = 0
	this.Mdepth = 0
	this.Mcaller = 0

	this.McallerTxid = 0
	this.McallerPcode = 0
	this.McallerSpec = ""
	this.McallerUrl = ""
	this.McallerPoidKey = ""

	//Processing data
	this.ServiceURL = nil
	this.RefererURL = nil
	this.IsStatic = false
	this.McallerUrlHash = 0
}
func (this *UdpTxStartEndPack) SetStaticContents(b bool) {
	this.IsStatic = b
	if b {
		this.IsStaticContents = "1"
	} else {
		this.IsStaticContents = "0"
	}
}

func (this *UdpTxStartEndPack) Write(dout *io.DataOutputX) {
	this.AbstractPack.Write(dout)
	dout.WriteTextShortLength(this.Host)
	dout.WriteTextShortLength(this.Uri)
	dout.WriteTextShortLength(this.Ipaddr)
	dout.WriteTextShortLength(this.UAgent)
	dout.WriteTextShortLength(this.Ref)
	dout.WriteTextShortLength(this.WClientId)
	dout.WriteTextShortLength(this.HttpMethod)
	dout.WriteTextShortLength(stringutil.ParseStringZeroToEmpty(int64(this.Mtid)))
	dout.WriteTextShortLength(stringutil.ParseStringZeroToEmpty(int64(this.Mdepth)))
	dout.WriteTextShortLength(stringutil.ParseStringZeroToEmpty(int64(this.McallerTxid)))
	dout.WriteTextShortLength(stringutil.ParseStringZeroToEmpty(int64(this.McallerPcode)))
	dout.WriteTextShortLength(this.McallerSpec)
	dout.WriteTextShortLength(this.McallerUrl)
	dout.WriteTextShortLength(this.McallerPoidKey)
}

func (this *UdpTxStartEndPack) Read(din *io.DataInputX) {
	this.AbstractPack.Read(din)

	this.Host = din.ReadTextShortLength()
	this.Uri = din.ReadTextShortLength()
	this.Ipaddr = din.ReadTextShortLength()
	this.UAgent = din.ReadTextShortLength()
	this.Ref = din.ReadTextShortLength()
	this.WClientId = din.ReadTextShortLength()

	this.HttpMethod = din.ReadTextShortLength()
	this.Mtid = stringutil.ParseInt64(din.ReadTextShortLength())
	this.Mdepth = stringutil.ParseInt32(din.ReadTextShortLength())
	this.McallerTxid = stringutil.ParseInt64(din.ReadTextShortLength())
	this.McallerPcode = stringutil.ParseInt64(din.ReadTextShortLength())
	this.McallerSpec = din.ReadTextShortLength()
	this.McallerUrl = din.ReadTextShortLength()
	this.McallerPoidKey = din.ReadTextShortLength()
}

func (this *UdpTxStartEndPack) Process() {
	if strings.HasPrefix(this.Uri, "/") {
		this.ServiceURL = urlutil.NewURL(this.Host + this.Uri)
	} else {
		this.ServiceURL = urlutil.NewURL(this.Host + "/" + this.Uri)
	}
	this.RefererURL = urlutil.NewURL(this.Ref)
	if this.IsStaticContents == "" {
		this.IsStatic = false
	} else if b, err := strconv.ParseBool(this.IsStaticContents); err == nil {
		this.IsStatic = b
	}

	if ret, err := strconv.ParseInt(this.McallerUrl, 10, 32); err == nil {
		this.McallerUrlHash = int32(ret)
	}
}
