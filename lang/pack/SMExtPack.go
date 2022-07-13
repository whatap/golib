package pack

import (
	"github.com/whatap/golib/io"
	"github.com/whatap/golib/lang/value"
)

type SMExtension struct {
	AbstractPack
	ver           byte
	isProjectwide bool
	values        *value.IntMapValue
	meta          *value.IntMapValue
	header        *value.IntMapValue
}

const (
	VERSION = 0x01
)

func (this *SMExtension) Write(dout *io.DataOutputX) {
	this.AbstractPack.Write(dout)
	bufHeader := io.NewDataOutputX()
	bufHeader.WriteBool(this.isProjectwide)
	this.header.WriteValue(bufHeader)
	dout.WriteByte(this.ver)
	dout.WriteBytes(bufHeader.ToByteArray())
	this.values.WriteValue(dout)
	this.meta.WriteValue(dout)
}

func (this *SMExtension) Read(din *io.DataInputX) {
	this.AbstractPack.Read(din)
	this.header = value.NewIntMapValue()
	this.values = value.NewIntMapValue()
	this.meta = value.NewIntMapValue()
	this.ver = din.ReadByte()
	this.isProjectwide = din.ReadBool()
	this.header.Read(din)
	din.ReadByte()
	this.values.Read(din)
	this.meta.Read(din)
}

func NewSMExtensionPack() *SMExtension {
	p := new(SMExtension)
	p.ver = VERSION
	return p
}

func (this *SMExtension) GetPackType() int16 {
	return PACK_SM_EXTENSION
}

func (this *SMExtension) SetValues(values *value.IntMapValue) {
	this.values = values
}

func (this *SMExtension) SetMetaValues(values *value.IntMapValue) {
	this.meta = values
}

func (this *SMExtension) SetHeader(header *value.IntMapValue) {
	this.header = header
}

func (this *SMExtension) SetIsProjectwide(isProjectwide bool) {
	this.isProjectwide = isProjectwide
}
