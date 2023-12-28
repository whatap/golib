package pack

import (
	"fmt"

	"github.com/whatap/golib/io"
)

const ()

type ZipPack struct {
	AbstractPack

	Records     []byte
	RecordCount int
	Status      byte
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
	dout.WriteDecimal(int64(this.RecordCount))
	dout.WriteBlob(this.Records)
}

func (this *ZipPack) Read(din *io.DataInputX) {
	this.AbstractPack.Read(din)

	this.Status = din.ReadByte()
	this.RecordCount = int(din.ReadDecimal())
	this.Records = din.ReadBlob()
}

//public ZipPack setRecords(int size, Enumeration<AbstractPack> items) {
//		this.RecordCount=size;
//		DataOutputX o = new DataOutputX();
//		for (int i = 0; i < size; i++) {
//			o.writePack(items.nextElement());
//		}
//		records = o.toByteArray();
//		return this;
//	}

func (this *ZipPack) SetRecords(items []Pack) *ZipPack {
	this.RecordCount = len(items)
	o := io.NewDataOutputX()
	for _, it := range items {
		o = WritePack(o, it)
	}
	this.Records = o.ToByteArray()
	return this
}

func (this *ZipPack) GetRecords() []Pack {
	items := make([]Pack, 0)
	if this.Records == nil {
		return nil
	}
	in := io.NewDataInputX(this.Records)
	for i := 0; i < this.RecordCount; i++ {
		p := ReadPack(in)

		p.SetPCODE(this.Pcode)
		p.SetOID(this.Oid)
		p.SetOKIND(this.Okind)
		p.SetONODE(this.Onode)

		items = append(items, p)
	}
	return items
}
