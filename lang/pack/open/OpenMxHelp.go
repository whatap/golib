package open

import (
	"github.com/whatap/golib/io"
	langvalue "github.com/whatap/golib/lang/value"
	"github.com/whatap/golib/util/stringutil"
)

type OpenMxHelp struct {
	Metric   string
	property *langvalue.MapValue
}

func NewOpenMxHelp() *OpenMxHelp {
	p := new(OpenMxHelp)
	p.property = langvalue.NewMapValue()

	return p
}

func (this *OpenMxHelp) Put(k, v string) *OpenMxHelp {
	this.property.PutString(k, v)
	return this
}

func (this *OpenMxHelp) Write(o *io.DataOutputX) {
	//version
	o.WriteByte(0)
	o.WriteText(this.Metric)
	this.property.Write(o)
}

func (this *OpenMxHelp) Read(in *io.DataInputX) *OpenMxHelp {
	//version
	// ver := in.ReadByte()
	_ = in.ReadByte()
	this.Metric = in.ReadText()
	m := langvalue.NewMapValue()
	m.Read(in)
	this.property = m
	return this
}

func (this *OpenMxHelp) ToString() string {
	return this.String()
}

func (this *OpenMxHelp) String() string {
	sb := stringutil.NewStringBuffer()
	sb.Append(" ").Append(this.Metric)
	sb.Append(" ").Append(this.property.String())
	return sb.ToString()
}
