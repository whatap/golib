package udp

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/whatap/golib/io"
	"github.com/whatap/golib/util/stringutil"
)

type UdpConfigPack struct {
	AbstractPack
	Data    string
	MapData map[string]string
}

func NewUdpConfigPack() *UdpConfigPack {
	p := new(UdpConfigPack)
	p.Ver = UDP_PACK_VERSION
	p.AbstractPack.Flush = true
	p.MapData = make(map[string]string)
	return p
}
func NewUdpConfigPackVer(ver int32) *UdpConfigPack {
	p := new(UdpConfigPack)
	p.Ver = ver
	p.AbstractPack.Flush = true
	p.MapData = make(map[string]string)
	return p
}

func (this *UdpConfigPack) GetPackType() uint8 {
	return CONFIG_INFO
}

func (this *UdpConfigPack) ToString() string {
	return fmt.Sprint("UdpConfigPack", ",data=", this.Data)
}

func (this *UdpConfigPack) Clear() {
	this.AbstractPack.Clear()
	this.AbstractPack.Flush = false
	this.Data = ""
}

func (this *UdpConfigPack) Write(dout *io.DataOutputX) {
	dout.WriteTextShortLength(this.Data)
}

func (this *UdpConfigPack) Read(din *io.DataInputX) {
	this.Data = din.ReadTextShortLength()
}

func (this *UdpConfigPack) Set(m map[string]string) {
	sb := stringutil.NewStringBuffer()
	for k, v := range m {
		sb.Append(k).Append("=").AppendLine(v)
	}
	this.Data = sb.ToString()
}
func (this *UdpConfigPack) Process() {
	br := strings.NewReader(this.Data)
	scanner := bufio.NewScanner(br)
	for scanner.Scan() {
		line := scanner.Text()
		pos := strings.Index(line, "=")
		if pos > -1 {
			this.MapData[line[0:pos]] = line[pos+1:]
		}
	}

	if this.Ver > 60000 {
		// Node.js
	} else if this.Ver > 50000 {
		// Golnag
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
