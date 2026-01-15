package udp

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/whatap/golib/io"
	"github.com/whatap/golib/lang/pack"
	"github.com/whatap/golib/util/dateutil"
)

type UdpTxParamPack struct {
	AbstractPack
	ParamId       string
	ParamResponse string
	Data          string
	StrDatas      []string
	ParamPack     *pack.ParamPack
}

func NewUdpTxParamPack() *UdpTxParamPack {
	p := new(UdpTxParamPack)
	p.Ver = UDP_PACK_VERSION
	p.AbstractPack.Flush = false
	return p
}
func NewUdpTxParamPackVer(ver int32) *UdpTxParamPack {
	p := new(UdpTxParamPack)
	p.Ver = ver
	p.AbstractPack.Flush = false
	return p
}

func (this *UdpTxParamPack) GetPackType() uint8 {
	return TX_PARAM
}

func (this *UdpTxParamPack) ToString() string {
	return fmt.Sprint(this.AbstractPack.ToString(), ",id=", this.ParamId, ",resp=", this.ParamResponse, ",data=", this.Data)
}

func (this *UdpTxParamPack) Clear() {
	this.AbstractPack.Clear()
	this.AbstractPack.Flush = false
	this.ParamId = ""
	this.ParamResponse = ""
	this.Data = ""
	this.StrDatas = nil
	this.ParamPack = nil
}

func (this *UdpTxParamPack) Write(dout *io.DataOutputX) {
	dout.WriteTextShortLength(this.ParamId)
	dout.WriteTextShortLength(this.ParamResponse)
	dout.WriteTextShortLength(this.Data)
}

func (this *UdpTxParamPack) Read(din *io.DataInputX) {
	this.ParamId = din.ReadTextShortLength()
	this.ParamResponse = din.ReadTextShortLength()
	this.Data = din.ReadTextShortLength()
}

func (this *UdpTxParamPack) Process() {
	this.ParamPack = pack.NewParamPack()
	this.ParamPack.Time = dateutil.Now()

	// ParamPack Id
	r, err := strconv.Atoi(this.ParamId)
	if err != nil {
		//logutil.Println("WA617", err)
	}
	this.ParamPack.Id = int32(r)

	// ParamPack Response
	r, err = strconv.Atoi(this.ParamResponse)
	if err != nil {
		//logutil.Println("WA618", err)
	}
	this.ParamPack.Response = int64(r)

	// ParamPack Data
	this.StrDatas = strings.Split(this.Data, ", ")

	//		var dataLen int
	//
	//		param := pack.NewParamPack()
	//		param.Time = dateutil.Now()
	//
	//		// ParamPack Id
	//		dataLen = int(io.ToShort(p.Data, pos))
	//		pos += 2
	//		paramId := string(p.Data[pos : pos+dataLen])
	//		pos += dataLen
	//		r, err := strconv.Atoi(paramId)
	//		if err != nil {
	//			logutil.Println("WA617", err)
	//		}
	//		param.Id = int32(r)
	//
	//		// ParamPack Response
	//		dataLen = int(io.ToShort(p.Data, pos))
	//		pos += 2
	//		response := string(p.Data[pos : pos+dataLen])
	//		pos += dataLen
	//		r, err = strconv.Atoi(response)
	//		if err != nil {
	//			logutil.Println("WA618", err)
	//		}
	//		param.Response = int64(r)
	//
	//		// ParamPack Data
	//		dataLen = int(io.ToShort(p.Data, pos))
	//		pos += 2
	//		str := string(p.Data[pos : pos+dataLen])
	//
	//		strDatas := strings.Split(str, ", ")
	//
	//		// Send Data
	//		data.Send(ParseMapValue(param, strDatas))
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
