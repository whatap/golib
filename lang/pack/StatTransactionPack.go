package pack

import (
	"container/list"
	"fmt"

	"github.com/whatap/golib/io"
	"github.com/whatap/golib/util/hmap"
)

type StatTransactionPack struct {
	AbstractPack
	// [] byte
	Records []byte
	// int
	RecordCount int

	// byte
	Version byte
}

func NewStatTransactionPack() *StatTransactionPack {
	p := new(StatTransactionPack)
	// 2021.06.28
	p.Version = 2
	return p
}

func (this *StatTransactionPack) GetPackType() int16 {
	return PACK_STAT_SERVICE
}

func (this *StatTransactionPack) ToString() string {
	//	sb.Append(",bytes=" + ArrayUtil.len(records));

	return fmt.Sprintln("StatService ", this.Oid, ",", this.Pcode, ",", this.Time, ",records=", this.RecordCount, ",bytes=", len(this.Records))
}

func (this *StatTransactionPack) Write(dout *io.DataOutputX) {
	this.AbstractPack.Write(dout)
	dout.WriteBlob(this.Records)
	dout.WriteDecimal(int64(this.RecordCount))
}

func (this *StatTransactionPack) Read(din *io.DataInputX) {
	this.AbstractPack.Read(din)
	this.Records = din.ReadBlob()
	this.RecordCount = int(din.ReadDecimal())
	//return this
}

func (this *StatTransactionPack) SetRecords(size int, items hmap.Enumeration) *StatTransactionPack {
	o := io.NewDataOutputX()
	o.WriteShort(int16(size))
	for i := 0; i < size; i++ {
		//fmt.Println("StatTransactionPack:SetRecords i=", i)
		WriteTransactionRec(o, items.NextElement().(*TransactionRec), this.Version)
	}
	this.Records = o.ToByteArray()
	this.RecordCount = size
	return this
}

func (this *StatTransactionPack) SetRecordsList(items *list.List) *StatTransactionPack {
	o := io.NewDataOutputX()
	size := items.Len()
	o.WriteShort(int16(size))
	for e := items.Front(); e != nil; e = e.Next() {
		WriteTransactionRec(o, e.Value.(*TransactionRec), this.Version)
	}
	this.Records = o.ToByteArray()
	this.RecordCount = size
	return this
}

func (this *StatTransactionPack) GetRecords() *list.List {
	items := list.New()

	if this.Records == nil {
		return nil
	}
	in := io.NewDataInputX(this.Records)
	size := int(in.ReadShort()) & 0xffff
	for i := 0; i < size; i++ {
		items.PushBack(ReadTransactionRec(in))
	}
	return items
}
