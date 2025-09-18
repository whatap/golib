package udp

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/whatap/golib/io"
	"github.com/whatap/golib/util/stringutil"
)

type UdpActiveStatsPack struct {
	AbstractPack
	Data        string
	ActiveStats []int16
}

func NewUdpActiveStatsPack() *UdpActiveStatsPack {
	p := new(UdpActiveStatsPack)
	p.Ver = UDP_PACK_VERSION
	p.AbstractPack.Flush = false
	return p
}
func NewUdpActiveStatsPackVer(ver int32) *UdpActiveStatsPack {
	p := new(UdpActiveStatsPack)
	p.Ver = ver
	p.AbstractPack.Flush = false
	return p
}

func (this *UdpActiveStatsPack) GetPackType() uint8 {
	return ACTIVE_STATS
}

func (this *UdpActiveStatsPack) ToString() string {
	return fmt.Sprint(this.AbstractPack.ToString(), ",data=", this.Data, ",active_stats=", this.ActiveStats)
}

func (this *UdpActiveStatsPack) Clear() {
	this.AbstractPack.Clear()
	this.AbstractPack.Flush = false
	this.Data = ""
	this.ActiveStats = nil
}

func (this *UdpActiveStatsPack) Write(dout *io.DataOutputX) {
	this.Data = stringutil.ArrayInt16ToString(this.ActiveStats, ",")
	dout.WriteTextShortLength(this.Data)
}

func (this *UdpActiveStatsPack) Read(din *io.DataInputX) {
	this.Data = din.ReadTextShortLength()
}

func (this *UdpActiveStatsPack) Process() {
	strStats := strings.Split(this.Data, ",")
	if len(strStats) == 5 {
		activeStats := make([]int16, 5, 5)
		for i := 0; i < 5; i++ {
			n, err := strconv.ParseInt(strStats[i], 10, 32)
			if err == nil {
				activeStats[i] = int16(n)
			}
		}
		this.ActiveStats = activeStats
	}

	if this.Ver > 60000 {
		// Node.js
	} else if this.Ver > 50000 {
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
