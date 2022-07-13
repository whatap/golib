package pack

import (
	//"log"

	"github.com/whatap/golib/io"
	"github.com/whatap/golib/util/bitutil"
	"github.com/whatap/golib/util/hmap"
)

//type ErrorRec struct {
//	ClassHash int32
//	Service   int32
//	SnapSeq   int64
//	Msg       int32
//	Count     int32
//}
type ErrorRec struct {
	ClassHash int32
	Service   int32
	SnapSeq   int64
	Msg       int32
	Count     int32
}

func NewErrorRec() *ErrorRec {
	p := new(ErrorRec)
	return p
}

func (this *ErrorRec) Merge(o *ErrorRec) {
	this.Count += o.Count
}

func (this *ErrorRec) GetKey() int64 {
	return bitutil.Composite64(this.ClassHash, this.Service)
}

func (this *ErrorRec) SetClassAndTxUrl(classHash int32, txUrl int32) {
	this.ClassHash = classHash
	this.Service = txUrl
}

type StatErrorPack struct {
	AbstractPack
	Records     []byte
	RecordCount int32
}

func NewStatErrorPack() *StatErrorPack {
	p := new(StatErrorPack)
	return p
}

func (this *StatErrorPack) GetPackType() int16 {
	return PACK_STAT_ERROR
}

func (this *StatErrorPack) Write(dout *io.DataOutputX) {
	this.AbstractPack.Write(dout)

	dout.WriteBlob(this.Records)
	dout.WriteDecimal(int64(this.RecordCount))

}
func (this *StatErrorPack) Read(din *io.DataInputX) {
	this.AbstractPack.Read(din)

	this.Records = din.ReadBlob()
	this.RecordCount = int32(din.ReadDecimal())
}
func (this *StatErrorPack) WriteRec(o *io.DataOutputX, m *ErrorRec) {
	o.WriteInt(m.ClassHash)
	o.WriteInt(m.Service)
	o.WriteLong(m.SnapSeq)
	o.WriteDecimal(int64(m.Msg))
	o.WriteDecimal(int64(m.Count))
}

func (this *StatErrorPack) ReadRec(in *io.DataInputX) *ErrorRec {
	m := new(ErrorRec)
	m.ClassHash = in.ReadInt()
	m.Service = in.ReadInt()
	m.SnapSeq = in.ReadLong()
	m.Msg = int32(in.ReadDecimal())
	m.Count = int32(in.ReadDecimal())
	return m
}

func (this *StatErrorPack) GetRecords() []*ErrorRec {
	in := io.NewDataInputX(this.Records)
	sz := int(in.ReadShort())
	items := make([]*ErrorRec, sz)
	for i := 0; i < sz; i++ {
		items[i] = this.ReadRec(in)
	}
	return items
}

func (this *StatErrorPack) SetRecords(size int, items hmap.Enumeration) *StatErrorPack {
	out := io.NewDataOutputX()
	out.WriteShort(int16(size))
	// DEBUG
	//fmt.Println("StatErrorPack size=", size)

	for i := 0; i < size; i++ {
		// DEBUG
		er := items.NextElement().(*ErrorRec)
		//fmt.Println("StatErrorPack ErrorRec=", er)
		this.WriteRec(out, er)

		//this.WriteRec(out, items.NextElement().(*ErrorRec))
	}

	this.Records = out.ToByteArray()
	this.RecordCount = int32(size)

	return this
}

func (this *StatErrorPack) SetRecordsArray(items []*ErrorRec) {
	out := io.NewDataOutputX()
	sz := len(items)
	out.WriteShort(int16(sz))
	for i := 0; i < sz; i++ {
		this.WriteRec(out, items[i])
	}
}
