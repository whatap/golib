package pack

import (
	//	"log"

	"github.com/whatap/golib/io"
	"github.com/whatap/golib/util/hmap"
)

const (
	STAT_USERAGENT_TABLE_MAX_SIZE = 500
)

type StatUserAgentPack struct {
	AbstractPack
	UserAgents *hmap.IntIntLinkedMap
}

func NewStatUserAgentPack() *StatUserAgentPack {
	p := new(StatUserAgentPack)
	p.UserAgents = hmap.NewIntIntLinkedMap().SetMax(STAT_USERAGENT_TABLE_MAX_SIZE)
	return p
}

func (this *StatUserAgentPack) GetPackType() int16 {
	return PACK_STAT_USER_AGENT
}

func (this *StatUserAgentPack) Write(dout *io.DataOutputX) {
	this.AbstractPack.Write(dout)
	sz := this.UserAgents.Size()

	dout.WriteDecimal(int64(sz))
	//DEBUG
	//fmt.Println("StatUserAgentPack sz=", sz)
	index := 0

	en := this.UserAgents.Entries()
	for en.HasMoreElements() {
		e := en.NextElement().(*hmap.IntIntLinkedEntry)
		dout.WriteInt(e.GetKey())
		dout.WriteInt(e.GetValue())

		//DEBUG
		//fmt.Println("StatUserAgentPack =", e.ToString())
		index++
	}

}
func (this *StatUserAgentPack) Read(din *io.DataInputX) {
	this.AbstractPack.Read(din)

	sz := int(din.ReadDecimal())
	for i := 0; i < sz; i++ {
		key := din.ReadInt()
		value := din.ReadInt()
		this.UserAgents.Put(key, value)
	}
}
