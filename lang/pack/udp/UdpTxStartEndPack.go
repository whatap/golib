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

	Status int32

	McallerStepId int64
	XTraceId      string

	PeakMem              int64
	ElapsedUserCPUTime   int32
	ElapsedSystemCPUTime int32
	EFuncCount           int32
	ProfEFuncCount       int32
	IFuncCount           int32
	ProfIFuncCount       int32

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

	this.Status = 0
	this.McallerStepId = 0
	this.XTraceId = ""

	this.PeakMem = 0
	this.ElapsedUserCPUTime = 0
	this.ElapsedSystemCPUTime = 0
	this.EFuncCount = 0
	this.ProfEFuncCount = 0
	this.IFuncCount = 0
	this.ProfIFuncCount = 0

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
	if this.Ver > 50000 {
		// Golang
		dout.WriteTextShortLength(this.HttpMethod)
		dout.WriteTextShortLength(stringutil.ParseStringZeroToEmpty(int64(this.Mtid)))
		dout.WriteTextShortLength(stringutil.ParseStringZeroToEmpty(int64(this.Mdepth)))
		dout.WriteTextShortLength(stringutil.ParseStringZeroToEmpty(int64(this.McallerTxid)))
		dout.WriteTextShortLength(stringutil.ParseStringZeroToEmpty(int64(this.McallerPcode)))
		dout.WriteTextShortLength(this.McallerSpec)
		dout.WriteTextShortLength(this.McallerUrl)
		dout.WriteTextShortLength(this.McallerPoidKey)
		if this.Ver >= 50100 {
			dout.WriteTextShortLength(stringutil.ParseStringZeroToEmpty(int64(this.Status)))
		}
		if this.Ver >= 50101 {
			dout.WriteTextShortLength(stringutil.ParseStringZeroToEmpty(this.McallerStepId))
			dout.WriteTextShortLength(this.XTraceId)
		}
	} else if this.Ver > 40000 {
		// Batch
	} else if this.Ver > 30000 {
		// Dotnet
		dout.WriteTextShortLength(this.IsStaticContents)
		dout.WriteTextShortLength(stringutil.ParseStringZeroToEmpty(int64(this.Mtid)))
		dout.WriteTextShortLength(stringutil.ParseStringZeroToEmpty(int64(this.Mdepth)))
		dout.WriteTextShortLength(stringutil.ParseStringZeroToEmpty(int64(this.Mcaller)))
	} else if this.Ver > 20000 {
		// Python
		dout.WriteTextShortLength(this.IsStaticContents)
		dout.WriteTextShortLength(stringutil.ParseStringZeroToEmpty(int64(this.Mtid)))
		dout.WriteTextShortLength(stringutil.ParseStringZeroToEmpty(int64(this.Mdepth)))
		dout.WriteTextShortLength(stringutil.ParseStringZeroToEmpty(int64(this.McallerTxid)))
		dout.WriteTextShortLength(stringutil.ParseStringZeroToEmpty(int64(this.McallerPcode)))
		dout.WriteTextShortLength(this.McallerSpec)
		dout.WriteTextShortLength(this.McallerUrl)
		dout.WriteTextShortLength(this.McallerPoidKey)
	} else {
		// PHP
		dout.WriteTextShortLength(this.HttpMethod)
		dout.WriteTextShortLength(stringutil.ParseStringZeroToEmpty(int64(this.Mtid)))
		dout.WriteTextShortLength(stringutil.ParseStringZeroToEmpty(int64(this.Mdepth)))
		dout.WriteTextShortLength(stringutil.ParseStringZeroToEmpty(int64(this.McallerTxid)))
		dout.WriteTextShortLength(stringutil.ParseStringZeroToEmpty(int64(this.McallerPcode)))
		dout.WriteTextShortLength(this.McallerSpec)
		dout.WriteTextShortLength(this.McallerUrl)
		dout.WriteTextShortLength(this.McallerPoidKey)

		if this.Ver >= 10107 {
			dout.WriteTextShortLength(stringutil.ParseStringZeroToEmpty(int64(this.Status)))
		}

		if this.Ver >= 10108 {
			dout.WriteTextShortLength(stringutil.ParseStringZeroToEmpty(this.McallerStepId))
			dout.WriteTextShortLength(this.XTraceId)
		}

		if this.Ver >= 10110 {
			dout.WriteTextShortLength(stringutil.ParseStringZeroToEmpty(this.PeakMem))
			dout.WriteTextShortLength(stringutil.ParseStringZeroToEmpty(int64(this.ElapsedUserCPUTime)))
			dout.WriteTextShortLength(stringutil.ParseStringZeroToEmpty(int64(this.ElapsedSystemCPUTime)))
			dout.WriteTextShortLength(stringutil.ParseStringZeroToEmpty(int64(this.EFuncCount)))
			dout.WriteTextShortLength(stringutil.ParseStringZeroToEmpty(int64(this.ProfEFuncCount)))
			dout.WriteTextShortLength(stringutil.ParseStringZeroToEmpty(int64(this.IFuncCount)))
			dout.WriteTextShortLength(stringutil.ParseStringZeroToEmpty(int64(this.ProfIFuncCount)))
		}

	}
}

func (this *UdpTxStartEndPack) Read(din *io.DataInputX) {
	this.AbstractPack.Read(din)

	this.Host = din.ReadTextShortLength()
	this.Uri = din.ReadTextShortLength()
	this.Ipaddr = din.ReadTextShortLength()
	this.UAgent = din.ReadTextShortLength()
	this.Ref = din.ReadTextShortLength()
	this.WClientId = din.ReadTextShortLength()

	if this.Ver > 50000 {
		// Golang
		this.HttpMethod = din.ReadTextShortLength()
		this.Mtid = stringutil.ParseInt64(din.ReadTextShortLength())
		this.Mdepth = stringutil.ParseInt32(din.ReadTextShortLength())
		this.McallerTxid = stringutil.ParseInt64(din.ReadTextShortLength())
		this.McallerPcode = stringutil.ParseInt64(din.ReadTextShortLength())
		this.McallerSpec = din.ReadTextShortLength()
		this.McallerUrl = din.ReadTextShortLength()
		this.McallerPoidKey = din.ReadTextShortLength()
		if this.Ver >= 50100 {
			this.Status = stringutil.ParseInt32(din.ReadTextShortLength())
		}
		if this.Ver >= 50101 {
			this.McallerStepId = stringutil.ParseInt64(din.ReadTextShortLength())
			this.XTraceId = din.ReadTextShortLength()
		}

	} else if this.Ver > 40000 {
		// Batch
	} else if this.Ver > 30000 {
		// Dotnet
		this.IsStaticContents = din.ReadTextShortLength()
		this.Mtid = stringutil.ParseInt64(din.ReadTextShortLength())
		this.Mdepth = stringutil.ParseInt32(din.ReadTextShortLength())
		this.Mcaller = stringutil.ParseInt64(din.ReadTextShortLength())
	} else if this.Ver > 20000 {
		// Python
		this.IsStaticContents = din.ReadTextShortLength()
		this.Mtid = stringutil.ParseInt64(din.ReadTextShortLength())
		this.Mdepth = stringutil.ParseInt32(din.ReadTextShortLength())
		this.McallerTxid = stringutil.ParseInt64(din.ReadTextShortLength())
		this.McallerPcode = stringutil.ParseInt64(din.ReadTextShortLength())
		this.McallerSpec = din.ReadTextShortLength()
		this.McallerUrl = din.ReadTextShortLength()
		this.McallerPoidKey = din.ReadTextShortLength()
	} else {
		// PHP
		this.HttpMethod = din.ReadTextShortLength()
		this.Mtid = stringutil.ParseInt64(din.ReadTextShortLength())
		this.Mdepth = stringutil.ParseInt32(din.ReadTextShortLength())
		this.McallerTxid = stringutil.ParseInt64(din.ReadTextShortLength())
		this.McallerPcode = stringutil.ParseInt64(din.ReadTextShortLength())
		this.McallerSpec = din.ReadTextShortLength()
		this.McallerUrl = din.ReadTextShortLength()
		this.McallerPoidKey = din.ReadTextShortLength()

		if this.Ver >= 10107 {
			// reponse code
			this.Status = stringutil.ParseInt32(din.ReadTextShortLength())
		}

		if this.Ver >= 10108 {
			this.McallerStepId = stringutil.ParseInt64(din.ReadTextShortLength())
			this.XTraceId = din.ReadTextShortLength()
		}

		if this.Ver >= 10110 {
			this.PeakMem = stringutil.ParseInt64(din.ReadTextShortLength())
			this.ElapsedUserCPUTime = stringutil.ParseInt32(din.ReadTextShortLength())
			this.ElapsedSystemCPUTime = stringutil.ParseInt32(din.ReadTextShortLength())
			this.EFuncCount = stringutil.ParseInt32(din.ReadTextShortLength())
			this.ProfEFuncCount = stringutil.ParseInt32(din.ReadTextShortLength())
			this.IFuncCount = stringutil.ParseInt32(din.ReadTextShortLength())
			this.ProfIFuncCount = stringutil.ParseInt32(din.ReadTextShortLength())
		}
	}
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

	if this.Ver > 50000 {
		// Golang
		if ret, err := strconv.ParseInt(this.McallerUrl, 10, 32); err == nil {
			this.McallerUrlHash = int32(ret)
		}
	} else if this.Ver > 40000 {
		// Batch
	} else if this.Ver > 30000 {
		// Dotnet
	} else if this.Ver > 20000 {
		// Python
		if ret, err := strconv.ParseInt(this.McallerUrl, 10, 32); err == nil {
			this.McallerUrlHash = int32(ret)
		}
	} else {
		// PHP
		if this.Ver >= 10102 {
			if ret, err := strconv.ParseInt(this.McallerUrl, 10, 32); err == nil {
				this.McallerUrlHash = int32(ret)
			}
		}
	}
}
