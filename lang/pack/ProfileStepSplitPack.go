package pack

import (
	"github.com/whatap/golib/io"
	"github.com/whatap/golib/lang/step"
	"github.com/whatap/golib/util/dateutil"
	"github.com/whatap/golib/util/stringutil"
)

type ProfileStepSplitPack struct {
	AbstractPack
	Txid  int64
	Inx   int
	Steps []byte
}

func NewProfileStepSplitPack() *ProfileStepSplitPack {
	p := new(ProfileStepSplitPack)
	p.Steps = make([]byte, 0)
	return p
}

func (this *ProfileStepSplitPack) GetPackType() int16 {
	return PACK_PROFILE_STEP_SPLIT
}

func (this *ProfileStepSplitPack) ToString() string {
	sb := stringutil.NewStringBuffer()
	sb.Append("Step ")
	sb.Append(this.AbstractPack.ToString())
	sb.Append(" time=" + dateutil.TimeStamp(this.Time))
	sb.AppendFormat(" inx=%d", this.Inx)
	sb.AppendFormat(" step_bytes=%d", len(this.Steps))

	return sb.ToString()
}

func (this *ProfileStepSplitPack) Write(dout *io.DataOutputX) {
	this.AbstractPack.Write(dout)
	dout.WriteByte(0)
	dout.WriteLong(this.Txid)
	dout.WriteDecimal(int64(this.Inx))
	dout.WriteBlob(this.Steps)
}

func (this *ProfileStepSplitPack) Read(din *io.DataInputX) *ProfileStepSplitPack {
	this.AbstractPack.Read(din)

	// ver
	din.ReadByte()
	this.Txid = din.ReadLong()
	this.Inx = int(din.ReadDecimal())
	this.Steps = din.ReadBlob()
	return this
}

func (this *ProfileStepSplitPack) SetProfile(steps []step.Step) *ProfileStepSplitPack {
	this.Steps = step.ToBytesStep(steps)
	return this
}
