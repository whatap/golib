package pack

import (
	//	"log"

	"github.com/whatap/golib/io"
	"github.com/whatap/golib/util/hmap"
)

const (
	STAT_REMOTEIP_TABLE_MAX_SIZE = 10000
)

type StatRemoteIpPack struct {
	AbstractPack
	IpTable *hmap.IntIntLinkedMap
}

func NewStatRemoteIpPack() *StatRemoteIpPack {
	p := new(StatRemoteIpPack)
	p.IpTable = hmap.NewIntIntLinkedMap().SetMax(STAT_REMOTEIP_TABLE_MAX_SIZE)
	return p
}

func (this *StatRemoteIpPack) GetPackType() int16 {
	return PACK_STAT_REMOTE_IP
}

func (this *StatRemoteIpPack) Write(dout *io.DataOutputX) {
	this.AbstractPack.Write(dout)
	sz := this.IpTable.Size()

	//DEBUG
	//fmt.Println("StatRemoteIpPack sz=", sz)
	index := 0

	dout.WriteDecimal(int64(sz))
	en := this.IpTable.Entries()
	for en.HasMoreElements() {
		e := en.NextElement().(*hmap.IntIntLinkedEntry)
		dout.WriteInt(e.GetKey())
		dout.WriteInt(e.GetValue())

		//DEBUG
		//fmt.Println("StatRemoteIpPack =", e.ToString())
		index++
	}

}
func (this *StatRemoteIpPack) Read(din *io.DataInputX) {
	this.AbstractPack.Read(din)

	sz := int(din.ReadDecimal())
	for i := 0; i < sz; i++ {
		key := din.ReadInt()
		value := din.ReadInt()
		this.IpTable.Put(key, value)
	}
}
