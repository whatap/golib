package pack

import (
	//"log"
	"container/list"
	"fmt"

	"github.com/whatap/golib/io"
	"github.com/whatap/golib/util/hmap"
)

type StatHttpcPack struct {
	AbstractPack
	Records     []byte
	RecordCount int32
}

func NewStatHttpcPack() *StatHttpcPack {
	p := new(StatHttpcPack)
	return p
}

// Implements Pack
func (this *StatHttpcPack) GetPackType() int16 {
	return PACK_STAT_HTTPC
}

// String()
//func (this *StatHttpcPack) String() string {
//	return this.ToString()
//}

func (this *StatHttpcPack) ToString() string {
	return fmt.Sprintln("StatHttpc ", this.AbstractPack.ToString(), ",records=", this.RecordCount, ",bytes=", len(this.Records))
}

func (this *StatHttpcPack) Write(dout *io.DataOutputX) {
	this.AbstractPack.Write(dout)

	dout.WriteBlob(this.Records)
	dout.WriteDecimal(int64(this.RecordCount))

}
func (this *StatHttpcPack) Read(din *io.DataInputX) {
	this.AbstractPack.Read(din)

	this.Records = din.ReadBlob()
	this.RecordCount = int32(din.ReadDecimal())
}

func (this *StatHttpcPack) SetRecords(size int, items hmap.Enumeration) *StatHttpcPack {
	o := io.NewDataOutputX()
	o.WriteShort(int16(size))
	for i := 0; i < size; i++ {
		items.NextElement().(*HttpcRec).Write(o)
	}
	this.Records = o.ToByteArray()
	this.RecordCount = int32(size)
	return this
}

func (this *StatHttpcPack) SetRecordsList(items *list.List) *StatHttpcPack {
	o := io.NewDataOutputX()

	o.WriteShort(int16(items.Len()))
	for e := items.Front(); e != nil; e = e.Next() {
		e.Value.(*HttpcRec).Write(o)
	}
	this.Records = o.ToByteArray()
	this.RecordCount = int32(items.Len())
	return this
}

func (this *StatHttpcPack) GetRecords() *list.List {
	items := list.New()
	if this.Records == nil {
		return items
	}
	in := io.NewDataInputX(this.Records)
	size := int(in.ReadShort()) & 0xffff
	for i := 0; i < size; i++ {
		items.PushBack(NewHttpcRec().Read(in))
	}

	return items
}
