package pack

import (
	"fmt"

	"github.com/whatap/golib/io"
	"github.com/whatap/golib/lang/value"
	"github.com/whatap/golib/util/hash"
)

const ()

type LogSinkPack struct {
	AbstractPack

	Category string
	TagHash  int64
	Tags     *value.MapValue
	Line     int64
	Content  string
	Fields   *value.MapValue
}

func NewLogSinkPack() *LogSinkPack {
	p := new(LogSinkPack)
	p.Tags = value.NewMapValue()
	p.Fields = value.NewMapValue()
	return p
}

func (this *LogSinkPack) GetPackType() int16 {
	return PACK_LOGSINK
}

func (this *LogSinkPack) ToString() string {
	return fmt.Sprintln("LogSinkPack ", this.AbstractPack.ToString(),
		" [category=", this.Category, ", tagHash=", this.TagHash, ", tags=", this.Tags,
		", content=", this.Content, ", fields=", this.Fields, "]")
}

func (this *LogSinkPack) GetTabAsBytes() []byte {
	out := io.NewDataOutputX()
	value.WriteMapValue(out, this.Tags)
	return out.ToByteArray()
}

func (this *LogSinkPack) GetContentBytes() []byte {
	out := io.NewDataOutputX()
	out.WriteByte(1)
	out.WriteText(this.Content)
	out.WriteDecimal(this.Line)
	return out.ToByteArray()
}

func (this *LogSinkPack) GetContent() string {
	if this.Content == "" {
		return ""
	} else {
		return this.Content
	}
}
func (this *LogSinkPack) SetContent(str string) {
	this.Content = str
}

func (this *LogSinkPack) SetContentBytes(d []byte) {
	defer func() {
		if r := recover(); r != nil {

		}
	}()

	if d == nil || len(d) < 1 {
		return
	}

	in := io.NewDataInputX(d)
	ver := in.ReadByte()
	if ver == 1 {
		this.Content = in.ReadText()
		this.Line = in.ReadDecimal()
	}
}

func (this *LogSinkPack) Write(dout *io.DataOutputX) {
	this.AbstractPack.Write(dout)

	dout.WriteByte(0)
	dout.WriteText(this.Category)

	if this.TagHash == 0 && this.Tags.Size() > 0 {
		tagBytes := this.ResetTagHash()
		dout.WriteDecimal(this.TagHash)
		dout.WriteBytes(tagBytes)
	} else {
		dout.WriteDecimal(this.TagHash)
		value.WriteMapValue(dout, this.Tags)
	}
	dout.WriteDecimal(this.Line)
	dout.WriteText(this.Content)
	if this.Fields != nil && this.Fields.Size() > 0 {
		dout.WriteBool(true)
		value.WriteMapValue(dout, this.Fields)
	} else {
		dout.WriteBool(false)
	}
}

func (this *LogSinkPack) Read(din *io.DataInputX) {
	this.AbstractPack.Read(din)
	//	ver := din.ReadByte()
	din.ReadByte()
	this.Category = din.ReadText()
	this.TagHash = din.ReadDecimal()
	this.Tags = value.ReadMapValue(din)
	this.Line = din.ReadDecimal()
	this.Content = din.ReadText()
	if din.ReadBool() {
		this.Fields = value.ReadMapValue(din)
	}
}

func (this *LogSinkPack) ResetTagHash() []byte {
	out := io.NewDataOutputX()
	value.WriteMapValue(out, this.Tags)
	tagBytes := out.ToByteArray()
	this.TagHash = hash.Hash64(tagBytes)
	return tagBytes
}

func (this *LogSinkPack) TransferOidToTag() {
	if this.Oid != 0 && this.Tags.ContainsKey("oid") == false {
		this.Tags.PutLong("oid", int64(this.Oid))
		this.TagHash = 0
	}
	if this.Okind != 0 && this.Tags.ContainsKey("okind") == false {
		this.Tags.PutLong("okind", int64(this.Okind))
		this.TagHash = 0
	}
	if this.Onode != 0 && this.Tags.ContainsKey("onode") == false {
		this.Tags.PutLong("onode", int64(this.Onode))
		this.TagHash = 0
	}
}
