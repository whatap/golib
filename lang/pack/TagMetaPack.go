package pack

import (
	"github.com/whatap/golib/io"
	"github.com/whatap/golib/lang/value"
	"github.com/whatap/golib/util/compressutil"
)

type TagMetaPack struct {
	AbstractPack
	Property *value.MapValue
	Json     string
}

func NewTagMetaPack() *TagMetaPack {
	p := new(TagMetaPack)
	p.Property = value.NewMapValue()
	return p
}

func (this *TagMetaPack) GetPackType() int16 {
	return PACK_TAG_META
}

func (this *TagMetaPack) Write(dout *io.DataOutputX) {
	this.AbstractPack.Write(dout)
	if this.Json == "" {
		dout.WriteByte(0)
		value.WriteValue(dout, this.Property)
	} else if len(this.Json) < 100 {
		dout.WriteByte(1)
		value.WriteValue(dout, this.Property)
		dout.WriteBlob([]byte(this.Json))
	} else {
		dout.WriteByte(2)
		value.WriteValue(dout, this.Property)
		if data, err := compressutil.DoZip([]byte(this.Json)); err == nil {
			dout.WriteBlob(data)
		} else {
			dout.WriteBlob([]byte(this.Json))
		}
	}
}

func (this *TagMetaPack) Read(din *io.DataInputX) {
	this.AbstractPack.Read(din)
	mode := din.ReadByte()
	this.Property = value.ReadValue(din).(*value.MapValue)
	switch mode {
	case 0:
		// property only
	case 1:
		this.Json = string(din.ReadBlob())
	case 2:
		if data, err := compressutil.UnZip(din.ReadBlob()); err == nil {
			this.Json = string(data)
		}
	}
}

func (this *TagMetaPack) ToString() string {
	return "TagMetaPack json=" + this.Json
}
