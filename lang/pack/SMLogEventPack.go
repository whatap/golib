package pack

import (
	"github.com/whatap/golib/io"
)

var (
	EVENTSOURCE_FILE     = byte(1)
	EVENTSOURCE_WINEVENT = byte(2)
	EVENTSOURCE_SCRIPT   = byte(3)
)

type SMLogEvent struct {
	EventSource byte
	Severity    byte
	FilePath    *string
	LogContent  *string

	WinLogFile    *string
	WinType       int32
	WinSourceName *string
	WinEventCode  int32
	WinCreateTime int64

	Keyword *string
	LogRule *string
}

func (this *SMLogEvent) Write(out *io.DataOutputX) {
	dout := io.NewDataOutputX()
	dout.WriteByte(this.EventSource)
	dout.WriteByte(this.Severity)
	if this.FilePath != nil {
		dout.WriteText(*this.FilePath)
	} else {
		dout.WriteText("")
	}
	if this.LogContent != nil {
		dout.WriteText(*this.LogContent)
	} else {
		dout.WriteText("")
	}
	if this.WinLogFile != nil {
		dout.WriteText(*this.WinLogFile)
	} else {
		dout.WriteText("")
	}

	dout.WriteInt(this.WinType)
	if this.WinSourceName != nil {
		dout.WriteText(*this.WinSourceName)
	} else {
		dout.WriteText("")
	}

	dout.WriteInt(this.WinEventCode)
	dout.WriteLong(this.WinCreateTime)
	dout.WriteText(*this.Keyword)
	dout.WriteText(*this.LogRule)

	out.WriteBlob(dout.ToByteArray())
}

func (this *SMLogEvent) Read(in *io.DataInputX) {
	din := io.NewDataInputX(in.ReadBlob())

	this.EventSource = din.ReadByte()
	this.Severity = din.ReadByte()
	filepath := din.ReadText()
	this.FilePath = &filepath
	logcontent := din.ReadText()
	this.LogContent = &logcontent

	winlogfile := din.ReadText()
	this.WinLogFile = &winlogfile
	this.WinType = din.ReadInt()
	winsourcename := din.ReadText()
	this.WinSourceName = &winsourcename
	this.WinEventCode = din.ReadInt()
	this.WinCreateTime = din.ReadLong()
	keyword := din.ReadText()
	this.Keyword = &keyword
	logRule := din.ReadText()
	this.LogRule = &logRule

}

type SMLogEventPack struct {
	AbstractPack
	LogEvent []SMLogEvent
}

func NewSMLogEventPack() *SMLogEventPack {
	p := new(SMLogEventPack)
	return p
}

func (this *SMLogEventPack) GetPackType() int16 {
	return PACK_SM_LOG_EVENT
}

func (this *SMLogEventPack) Write(dout *io.DataOutputX) {
	this.AbstractPack.Write(dout)
	dout.WriteDecimal(int64(len(this.LogEvent)))
	for _, logEvent := range this.LogEvent {
		logEvent.Write(dout)
	}
}

func (this *SMLogEventPack) Read(din *io.DataInputX) {
	this.AbstractPack.Read(din)
	eventCount := din.ReadDecimal()
	if eventCount > 0 {
		this.LogEvent = make([]SMLogEvent, eventCount)
		for i := int64(0); i < eventCount; i++ {
			this.LogEvent[i].Read(din)
		}
	}
}
