package udp

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/whatap/golib/io"
)

type UdpDBConPoolPack struct {
	AbstractPack
	Data     string
	Pid      int32
	Url      string
	ActCnt   int32
	InactCnt int32
}

func NewUdpDBConPoolPack() *UdpDBConPoolPack {
	p := new(UdpDBConPoolPack)
	p.Ver = UDP_PACK_VERSION
	p.AbstractPack.Flush = false
	return p
}
func NewUdpDBConPoolPackVer(ver int32) *UdpDBConPoolPack {
	p := new(UdpDBConPoolPack)
	p.Ver = ver
	p.AbstractPack.Flush = false
	return p
}

func (this *UdpDBConPoolPack) GetPackType() uint8 {
	return DBCONN_POOL
}

func (this *UdpDBConPoolPack) ToString() string {
	return fmt.Sprint(this.AbstractPack.ToString(), ", data=", this.Data, ",Pid=", this.Pid, ", Url=", this.Url, ", ActCnt=", this.ActCnt, ", inactCnt=", this.InactCnt)
}

func (this *UdpDBConPoolPack) Clear() {
	this.AbstractPack.Clear()
	this.AbstractPack.Flush = false
	this.Data = ""
	this.Pid = 0
	this.Url = ""
	this.ActCnt = 0
	this.InactCnt = 0
}

func (this *UdpDBConPoolPack) Write(dout *io.DataOutputX) {
	dout.WriteTextShortLength(this.Data)
}

func (this *UdpDBConPoolPack) Read(din *io.DataInputX) {
	this.Data = din.ReadTextShortLength()
}

func (this *UdpDBConPoolPack) Process() {
	strDbConns := strings.Split(this.Data, ",")
	for _, strDbConn := range strDbConns {
		words := strings.Split(strDbConn, "|")
		if len(words) == 4 {
			pid, _ := strconv.ParseInt(words[0], 10, 32)
			this.Pid = int32(pid)
			this.Url = words[1]
			actCnt, _ := strconv.ParseInt(words[2], 10, 32)
			this.ActCnt = int32(actCnt)
			inactCnt, _ := strconv.ParseInt(words[3], 10, 32)
			this.InactCnt = int32(inactCnt)
		}
	}
	//
	//		dataLen := int(io.ToShort(p.Data, pos))
	//		pos += 2
	//		strDbConns := strings.Split(string(p.Data[pos:pos+dataLen]), ",")
	//		for _, strDbConn := range strDbConns {
	//			words := strings.Split(strDbConn, "|")
	//			if len(words) == 4 {
	//				pid, _ := strconv.ParseInt(words[0], 10, 32)
	//				url := words[1]
	//				actCnt, _ := strconv.ParseInt(words[2], 10, 32)
	//				inactCnt, _ := strconv.ParseInt(words[3], 10, 32)
	//				trace.AddDBConnPool(int32(pid), url, int32(actCnt), int32(inactCnt))
	//			}
	//
	//		}
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
