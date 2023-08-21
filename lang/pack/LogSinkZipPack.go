package pack

import (
	"github.com/whatap/golib/io"
	"github.com/whatap/golib/util/compressutil"
	"github.com/whatap/golib/util/stringutil"
)

const (
	ZIPPED    = 1
	UN_ZIPPED = 0
)

type LogSinkZipPack struct {
	AbstractPack

	Records     []byte
	RecordCount int
	Status      byte
}

func NewLogSinkZipPack() *LogSinkZipPack {
	p := new(LogSinkZipPack)
	return p
}

func (this *LogSinkZipPack) GetPackType() int16 {
	return PACK_LOGSINK_ZIP
}

func (this *LogSinkZipPack) Write(dout *io.DataOutputX) {
	this.AbstractPack.Write(dout)

	dout.WriteByte(this.Status)
	dout.WriteDecimal(int64(this.RecordCount))
	dout.WriteBlob(this.Records)
}
func (this *LogSinkZipPack) Read(din *io.DataInputX) {
	this.AbstractPack.Read(din)

	this.Status = din.ReadByte()
	this.RecordCount = int(din.ReadDecimal())
	this.Records = din.ReadBlob()
}

func (this *LogSinkZipPack) SetRecords(records []byte, zipMinSize int) {
	this.Records, _ = this.doZip(records, zipMinSize)
}

func (this *LogSinkZipPack) doZip(records []byte, zipMinSize int) ([]byte, error) {
	if this.Status != UN_ZIPPED {
		return records, nil
	}
	if len(records) < zipMinSize {
		return records, nil
	}
	this.Status = ZIPPED
	return compressutil.DoZip(records)
}

func (this *LogSinkZipPack) doUnZip() ([]byte, error) {
	if this.Status != ZIPPED {
		return this.Records, nil
	}
	return compressutil.UnZip(this.Records)
}

func (this *LogSinkZipPack) GetRecords() []*LogSinkPack {
	items := make([]*LogSinkPack, 0)
	if this.Records == nil {
		return items
	}
	data, err := this.doUnZip()
	if err != nil {
		return items
	}

	in := io.NewDataInputX(data)
	for i := 0; i < this.RecordCount; i++ {
		tmp := ReadPack(in)
		if tmp != nil {
			if p, ok := tmp.(*LogSinkPack); ok {
				p.Pcode = this.Pcode
				p.Oid = this.Oid
				p.Okind = this.Okind
				p.Onode = this.Onode
				// time 은 자기시간 사용
				items = append(items, p)
			}
		}
	}

	return items

}

func (this *LogSinkZipPack) ToString() string {
	sb := stringutil.NewStringBuffer()
	sb.Append("LogSinkZipPack ")
	sb.Append(this.AbstractPack.ToString())
	sb.AppendFormat("records=%d bytes", len(this.Records))
	return sb.ToString()
}
