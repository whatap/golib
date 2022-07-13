package udp

import (
	"fmt"

	"github.com/whatap/golib/io"
	"github.com/whatap/golib/util/stringutil"
)

type UdpTxSecureMessagePack struct {
	AbstractPack
	Hash  string
	Value string
	Desc  string
}

func NewUdpTxSecureMessagePack() *UdpTxSecureMessagePack {
	p := new(UdpTxSecureMessagePack)
	p.Ver = UDP_PACK_VERSION
	p.AbstractPack.Flush = false
	return p
}

func NewUdpTxSecureMessagePackVer(ver int32) *UdpTxSecureMessagePack {
	p := new(UdpTxSecureMessagePack)
	p.Ver = ver
	p.AbstractPack.Flush = false
	return p
}

func (this *UdpTxSecureMessagePack) GetPackType() uint8 {
	return TX_SECURE_MSG
}

func (this *UdpTxSecureMessagePack) ToString() string {
	return fmt.Sprint(this.AbstractPack.ToString(), ",hash=", this.Hash, ",value=", this.Value, ",desc=", this.Desc)
}

func (this *UdpTxSecureMessagePack) Clear() {
	this.AbstractPack.Clear()
	this.AbstractPack.Flush = false

	this.Hash = ""
	this.Value = ""
	this.Desc = ""
}

func (this *UdpTxSecureMessagePack) SetParameter(m map[string][]string) {
	this.Desc = ParseParameter(m, HTTP_PARAM_MAX_COUNT, HTTP_PARAM_KEY_MAX_SIZE, HTTP_PARAM_VALUE_MAX_SIZE)
}

func (this *UdpTxSecureMessagePack) Write(dout *io.DataOutputX) {
	this.AbstractPack.Write(dout)
	dout.WriteTextShortLength(this.Hash)
	dout.WriteTextShortLength(this.Value)
	dout.WriteTextShortLength(this.Desc)
}

func (this *UdpTxSecureMessagePack) Read(din *io.DataInputX) {
	this.AbstractPack.Read(din)

	this.Hash = din.ReadTextShortLength()
	this.Value = din.ReadTextShortLength()
	this.Desc = din.ReadTextShortLength()
}
func (this *UdpTxSecureMessagePack) Process() {
}

func ParseParameter(m map[string][]string, maxCount, keyMaxSize, valueMaxSize int) string {
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
