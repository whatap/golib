package pack

import (
	"fmt"

	"github.com/whatap/golib/io"
)

const ()

type ZipPack struct {
	AbstractPack

	Records      []byte
	RecountCount int
	Status       byte
}

func NewZipPack() *ZipPack {
	p := new(ZipPack)
	return p
}

func (this *ZipPack) GetPackType() int16 {
	return PACK_ZIP
}

func (this *ZipPack) ToString() string {
	return fmt.Sprintln("ZipPack ", this.AbstractPack.ToString(),
		" records=", len(this.Records), " bytes")
}

func (this *ZipPack) Write(dout *io.DataOutputX) {
	this.AbstractPack.Write(dout)

	dout.WriteByte(this.Status)
	dout.WriteDecimal(int64(this.RecountCount))
	dout.WriteBlob(this.Records)
}

func (this *ZipPack) Read(din *io.DataInputX) {
	this.AbstractPack.Read(din)

	this.Status = din.ReadByte()
	this.RecountCount = int(din.ReadDecimal())
	this.Records = din.ReadBlob()
}

//public ZipPack setRecords(int size, Enumeration<AbstractPack> items) {
//		this.recountCount=size;
//		DataOutputX o = new DataOutputX();
//		for (int i = 0; i < size; i++) {
//			o.writePack(items.nextElement());
//		}
//		records = o.toByteArray();
//		return this;
//	}

func (this *ZipPack) SetRecords(items []*AbstractPack) *ZipPack {
	this.RecountCount = len(items)
	o := io.NewDataOutputX()
	for _, it := range items {
		it.Write(o)
	}
	this.Records = o.ToByteArray()
	return this
}

func (this *ZipPack) GetRecords() []*AbstractPack {
	items := make([]*AbstractPack, 0)
	if this.Records == nil {
		return nil
	}
	in := io.NewDataInputX(this.Records)
	for i := 0; i < this.RecountCount; i++ {
		p := new(AbstractPack)
		p.Read(in)

		p.Pcode = this.Pcode
		p.Oid = this.Oid
		// time은 자기 시간을 사용한다.
		p.Okind = this.Okind
		p.Onode = this.Onode

		items = append(items, p)
	}
	return items
}
