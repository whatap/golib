package udp

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/whatap/golib/io"
)

type UdpActiveStackPack struct {
	AbstractPack
	Data  string
	TxId  int64
	Stack string
}

func NewUdpActiveStackPack() *UdpActiveStackPack {
	p := new(UdpActiveStackPack)
	p.Ver = UDP_PACK_VERSION
	p.AbstractPack.Flush = false
	return p
}
func NewUdpActiveStackPackVer(ver int32) *UdpActiveStackPack {
	p := new(UdpActiveStackPack)
	p.Ver = ver
	p.AbstractPack.Flush = false
	return p
}

func (this *UdpActiveStackPack) GetPackType() uint8 {
	return ACTIVE_STACK
}

func (this *UdpActiveStackPack) ToString() string {
	return fmt.Sprint(this.AbstractPack.ToString(), ",data=", this.Data)
}

func (this *UdpActiveStackPack) Clear() {
	this.AbstractPack.Clear()
	this.AbstractPack.Flush = false
	this.Data = ""
	this.TxId = 0
	this.Stack = ""
}

func (this *UdpActiveStackPack) Write(dout *io.DataOutputX) {
	dout.WriteTextShortLength(this.Data)
}

func (this *UdpActiveStackPack) Read(din *io.DataInputX) {
	this.Data = din.ReadTextShortLength()
}

func (this *UdpActiveStackPack) Process() {
	strDatas := strings.Split(this.Data, ", ")
	txId := strDatas[1]
	r, err := strconv.Atoi(txId)
	if err != nil {
		//logutil.Println("WA665", err)
	} else {
		this.TxId = int64(r)
	}
	this.Stack = strDatas[2]

	// txId
	//		dataLen = int(io.ToShort(p.Data, pos))
	//		pos += 2
	//		strDatas := strings.Split(string(p.Data[pos:pos+dataLen]), ", ")
	//		txId := strDatas[1]
	//		r, err := strconv.Atoi(txId)
	//		if err != nil {
	//			logutil.Println("WA665", err)
	//		}
	//		active.SendActiveStack(int64(r), strDatas[2])

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
