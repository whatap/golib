package open

import (
	"github.com/whatap/golib/io"
	"github.com/whatap/golib/util/hash"
)

type Label struct {
	Key   string
	Value string
	_id_  int
}

func NewLabel(k, v string) *Label {
	p := new(Label)
	p.Key = k
	p.Value = v
	return p
}

func (this *Label) Write(o *io.DataOutputX) {
	o.WriteText(this.Key)
	o.WriteText(this.Value)
}

func ReadStatic(in *io.DataInputX) *Label {
	k := in.ReadText()
	v := in.ReadText()
	return NewLabel(k, v)
}

func (this *Label) String() string {
	return this.Key + "=" + this.Value
}

func (this *Label) ID() int {
	if this._id_ == 0 {
		this._id_ = int(hash.HashStr(this.Key)) ^ int(hash.HashStr(this.Value))
	}
	return this._id_
}
