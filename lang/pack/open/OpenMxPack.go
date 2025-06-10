package open

import (
	"fmt"

	"github.com/whatap/golib/io"
	"github.com/whatap/golib/util/compressutil"
	"github.com/whatap/golib/util/stringutil"
)

type OpenMxPack struct {
	AbstractPack
	zip     byte
	bytes   []byte
	records []*OpenMx
}

func NewOpenMxPack() *OpenMxPack {
	p := new(OpenMxPack)
	p.bytes = make([]byte, 0)
	p.records = make([]*OpenMx, 0)
	return p
}
func (this *OpenMxPack) GetPackType() int16 {
	return PACK_OPEN_MX_HELP_PACK
}

func (this *OpenMxPack) ToStrint() string {
	return this.String()
}
func (this *OpenMxPack) String() string {
	sb := stringutil.NewStringBuffer()
	sb.Append("pack:").Append("OpenMxPack")
	sb.Append("zip:").Append(fmt.Sprintf("%d", this.zip))
	sb.Append("bytes").Append(fmt.Sprintf("%d", len(this.bytes)))
	return sb.ToString()
}

func (this *OpenMxPack) Write(dout *io.DataOutputX) {
	this.AbstractPack.Write(dout)
	if this.bytes == nil || len(this.bytes) == 0 {
		this.reset(this.records)
	}
	dout.WriteByte(this.zip)
	dout.WriteBlob(this.bytes)
}

func (this *OpenMxPack) Read(din *io.DataInputX) {
	this.AbstractPack.Read(din)
	this.zip = din.ReadByte()
	this.bytes = din.ReadBlob()
}

func (this *OpenMxPack) SetRecords(items []*OpenMx) *OpenMxPack {
	this.records = items
	return this.reset(items)
}

func (this *OpenMxPack) reset(items []*OpenMx) *OpenMxPack {
	o := io.NewDataOutputX()
	o.WriteByte(0) //version
	if items == nil || len(items) == 0 {
		o.WriteShort(0)
	} else {
		o.WriteShort(int16(len(items)))
		for _, it := range items {
			it.Write(o)
		}
	}
	this.bytes = o.ToByteArray()
	if len(this.bytes) > 100 {
		if data, err := compressutil.DoZip(this.bytes); err == nil {
			this.zip = 1
			this.bytes = data
		}
	}
	return this
}

func (this *OpenMxPack) GetUnpack() []*OpenMx {
	if this.bytes == nil {
		return make([]*OpenMx, 0)
	}

	var in *io.DataInputX
	if this.zip == 1 {
		if data, err := compressutil.UnZip(this.bytes); err == nil {
			in = io.NewDataInputX(data)
		} else {
			in = io.NewDataInputX(this.bytes)
		}
	} else {
		in = io.NewDataInputX(this.bytes)
	}

	this.records = make([]*OpenMx, 0)
	_ = in.ReadByte()
	size := int(in.ReadShort())
	for i := 0; i < size; i++ {
		this.records = append(this.records, NewOpenMx().Read(in))
	}

	return this.records
}

func (this *OpenMxPack) GetRecords() []*OpenMx {
	if this.bytes == nil {
		return make([]*OpenMx, 0)
	}
	if this.records == nil {
		return this.GetUnpack()
	}
	return this.records
}
