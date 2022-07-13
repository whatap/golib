package udp

import (
	"fmt"

	"github.com/whatap/golib/io"
	"github.com/whatap/golib/util/stringutil"
)

type UdpTxMessagePack struct {
	AbstractPack
	Hash  string
	Value string
	Desc  string
}

func NewUdpTxMessagePack() *UdpTxMessagePack {
	p := new(UdpTxMessagePack)
	p.Ver = UDP_PACK_VERSION
	p.AbstractPack.Flush = false
	return p
}

func NewUdpTxMessagePackVer(ver int32) *UdpTxMessagePack {
	p := new(UdpTxMessagePack)
	p.Ver = ver
	p.AbstractPack.Flush = false
	return p
}

func (this *UdpTxMessagePack) GetPackType() uint8 {
	return TX_MSG
}

func (this *UdpTxMessagePack) ToString() string {
	return fmt.Sprint(this.AbstractPack.ToString(), ",hash=", this.Hash, ",value=", this.Value, ",desc=", this.Desc)
}

func (this *UdpTxMessagePack) Clear() {
	this.AbstractPack.Clear()
	this.AbstractPack.Flush = false

	this.Hash = ""
	this.Value = ""
	this.Desc = ""
}
func (this *UdpTxMessagePack) SetHeader(m map[string][]string) {
	this.Desc = ParseHeader(m, HTTP_HEADER_MAX_COUNT, HTTP_HEADER_KEY_MAX_SIZE, HTTP_HEADER_VALUE_MAX_SIZE)
}

func (this *UdpTxMessagePack) Write(dout *io.DataOutputX) {
	this.AbstractPack.Write(dout)
	dout.WriteTextShortLength(stringutil.Truncate(this.Hash, HTTP_URI_MAX_SIZE))
	dout.WriteTextShortLength(this.Value)
	dout.WriteTextShortLength(stringutil.Truncate(this.Desc, PACKET_MESSAGE_MAX_SIZE))
}

func (this *UdpTxMessagePack) Read(din *io.DataInputX) {
	this.AbstractPack.Read(din)

	this.Hash = din.ReadTextShortLength()
	this.Value = din.ReadTextShortLength()
	this.Desc = din.ReadTextShortLength()
}
func (this *UdpTxMessagePack) Process() {
}

func ParseHeader(m map[string][]string, maxCount, keyMaxSize, valueMaxSize int) string {
	rt := ""
	if m != nil && len(m) > 0 {
		sb := stringutil.NewStringBuffer()
		idx := 0
		for k, v := range m {
			if idx > maxCount {
				break
			}
			sb.Append(stringutil.Truncate(k, keyMaxSize)).Append("=")
			if len(v) > 0 {
				sb.AppendLine(stringutil.Truncate(v[0], valueMaxSize))
			}
		}
		rt = sb.ToString()
		sb.Clear()
	}
	return rt
}
