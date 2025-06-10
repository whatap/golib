package open

import (
	"github.com/whatap/golib/io"
	"github.com/whatap/golib/util/dateutil"
	"github.com/whatap/golib/util/stringutil"
)

type OpenMx struct {
	Metric    string
	label     []*Label
	timestamp int64
	val       float64
}

func NewOpenMx() *OpenMx {
	p := new(OpenMx)
	return p
}

func NewOpenMxWithValue(metric string, timestamp int64, value float64) *OpenMx {
	p := new(OpenMx)
	p.Metric = metric
	p.timestamp = timestamp
	p.val = value
	p.label = make([]*Label, 0)
	return p
}
func NewOpenMxWithLabel(metric string, timestamp int64, value float64, labels ...*Label) *OpenMx {
	p := new(OpenMx)
	p.Metric = metric
	p.timestamp = timestamp
	p.val = value
	p.label = make([]*Label, 0)

	for _, lb := range labels {
		p.label = append(p.label, lb)
	}
	return p
}

func (this *OpenMx) AddLabel(k, v string) {
	if this.label == nil {
		this.label = make([]*Label, 0)
	}
	this.label = append(this.label, NewLabel(k, v))
}

func (this *OpenMx) Write(o *io.DataOutputX) {
	o.WriteByte(0) // version
	o.WriteText(this.Metric)
	labelSize := 0
	if this.label != nil {
		labelSize = len(this.label)
	}
	o.WriteByte(byte(labelSize))
	for _, lb := range this.label {
		lb.Write(o)
	}
	o.WriteLong(this.timestamp)
	o.WriteDouble(this.val)
}

func (this *OpenMx) Read(in *io.DataInputX) *OpenMx {
	// version
	_ = in.ReadByte()
	this.Metric = in.ReadText()
	cnt := int(in.ReadByte())
	if cnt > 0 {
		this.label = make([]*Label, 0)
		for i := 0; i < cnt; i++ {
			this.label = append(this.label, ReadStatic(in))
		}
	}
	this.timestamp = in.ReadLong()
	this.val = in.ReadDouble()

	return this
}

func (this *OpenMx) String() string {
	sb := stringutil.NewStringBuffer()
	sb.Append("OpenMx [")
	sb.Append(this.Metric)
	if this.label != nil && len(this.label) > 0 {
		sb.Append("label:[")
		for _, it := range this.label {
			sb.Append(it.String()).Append(",")
		}
		sb.Append("]")
	}
	sb.Append(" ")
	sb.Append(dateutil.TimeStamp(this.timestamp))

	return sb.ToString()
}
