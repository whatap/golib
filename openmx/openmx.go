package openmx

import (
	"github.com/whatap/golib/lang/pack"
	"github.com/whatap/golib/lang/pack/open"
)

const (
	OPENMX_TYPE_GAUGE     = "gauge"
	OPENMX_TYPE_COUNTER   = "counter"
	OPENMX_TYPE_SUMMARY   = "summary"
	OPENMX_TYPE_HISTOGRAM = "histogram"
)

var (
	helps   []*open.OpenMxHelp
	metrics []*open.OpenMx
	// default label . add always
	labels []*open.Label
)

type OpenMetric struct {
	helps   []*open.OpenMxHelp
	metrics []*open.OpenMx
	labels  []*open.Label
}

func New() *OpenMetric {
	p := new(OpenMetric)
	p.helps = make([]*open.OpenMxHelp, 0)
	p.metrics = make([]*open.OpenMx, 0)
	return p
}

func (this *OpenMetric) AddHelp(m, h, t string) {
	help := open.NewOpenMxHelp()
	help.Metric = m
	help.Put("help", h)
	help.Put("type", t)
	this.helps = append(this.helps, help)
}
func (this *OpenMetric) AddHelpObject(help *open.OpenMxHelp) {
	this.helps = append(this.helps, help)
}
func (this *OpenMetric) GetOpenMxHelpPack() pack.Pack {
	h := this.resetOpenMxHelp()
	p := open.NewOpenMxHelpPack()
	p.SetRecords(h)
	return p
}
func (this *OpenMetric) resetOpenMxHelp() []*open.OpenMxHelp {
	sz := len(this.helps)
	result := make([]*open.OpenMxHelp, sz)
	copy(result, this.helps)
	this.helps = make([]*open.OpenMxHelp, 0)
	return result
}

func (this *OpenMetric) AddMetric(m *open.OpenMx) {
	// default labels
	if len(this.labels) > 0 {
		for _, it := range this.labels {
			m.AddLabel(it.Key, it.Value)
		}
	}
	this.metrics = append(this.metrics, m)
}

func (this *OpenMetric) GetOpenMxPack() pack.Pack {
	m := this.resetOpenMx()
	p := open.NewOpenMxPack()
	p.SetRecords(m)
	return p
}
func (this *OpenMetric) resetOpenMx() []*open.OpenMx {
	sz := len(this.metrics)
	result := make([]*open.OpenMx, sz)
	copy(result, this.metrics)
	this.metrics = make([]*open.OpenMx, 0)
	return result
}

func (this *OpenMetric) AddDefaultLabel(k, v string) {
	this.labels = append(this.labels, open.NewLabel(k, v))
}

func (this *OpenMetric) SetDefaultLabels(lbs ...*open.Label) {
	this.labels = append(this.labels, lbs...)
}

func (this *OpenMetric) ResetDefaultLabels(lbs ...*open.Label) []*open.Label {
	sz := len(this.labels)
	result := make([]*open.Label, sz)
	copy(result, this.labels)
	this.labels = make([]*open.Label, 0)
	if len(lbs) > 0 {
		this.labels = append(this.labels, lbs...)
	}
	return result
}

func (this *OpenMetric) Send(tm int64, send func(pack.Pack)) {
	this.sendOpenMxHelpPack(tm, send)
	this.sendOpenMxPack(tm, send)
}
func (this *OpenMetric) sendOpenMxHelpPack(tm int64, send func(pack.Pack)) {
	h := this.helps
	p := open.NewOpenMxHelpPack()
	p.SetRecords(h)
	p.SetTime(tm)
	send(p)
	this.helps = make([]*open.OpenMxHelp, 0)
}
func (this *OpenMetric) sendOpenMxPack(tm int64, send func(pack.Pack)) {
	m := this.metrics
	p := open.NewOpenMxPack()
	p.SetRecords(m)
	p.SetTime(tm)
	send(p)
	this.metrics = make([]*open.OpenMx, 0)
}
