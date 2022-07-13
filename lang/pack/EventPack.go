package pack

import (
	"fmt"
	"strconv"

	"github.com/whatap/golib/io"
	"github.com/whatap/golib/util/hmap"
	"github.com/whatap/golib/util/uuidutil"
)

const (
	FATAL   byte = 30
	WARNING byte = 20
	INFO    byte = 10
	NONE    byte = 0

	ESCALATION_KEY string = "_esca_"
	UUID_KEY       string = "_uuid_"
	STATUS_KEY     string = "_status_"
	OTYPE_KEY      string = "_otype_"
)

type EventPack struct {
	AbstractPack

	Uuid       string
	Escalation bool

	Level   byte
	Title   string
	Message string
	Status  int32
	Otype   int32

	Attr *hmap.StringKeyLinkedMap

	Eid int64
}

func NewEventPack() *EventPack {
	p := new(EventPack)
	p.Attr = hmap.NewStringKeyLinkedMap()

	return p
}

func (this *EventPack) GetPackType() int16 {
	return PACK_EVENT
}

func (this *EventPack) Size() int {
	size := 1
	if this.Title != "" {
		size += len(this.Title)
	}
	if this.Message != "" {
		size += len(this.Message)
	}
	size += this.Attr.Size() * 20
	return size
}

func (this *EventPack) Write(dout *io.DataOutputX) {
	this.AbstractPack.Write(dout)

	dout.WriteByte(this.Level)
	dout.WriteText(this.Title)
	dout.WriteText(this.Message)

	if this.Uuid != "" {
		this.Attr.Put(UUID_KEY, this.Uuid)
	}
	if this.Escalation {
		this.Attr.Put(ESCALATION_KEY, "true")
	} else {
		this.Attr.Put(ESCALATION_KEY, "false")
	}
	this.Attr.Put(STATUS_KEY, fmt.Sprintf("%d", this.Status))
	this.Attr.Put(OTYPE_KEY, fmt.Sprintf("%d", this.Otype))

	sz := this.Attr.Size()
	dout.WriteByte(byte(sz))
	en := this.Attr.Entries()
	for i := 0; i < sz; i++ {
		e := en.NextElement().(*hmap.StringKeyLinkedEntry)
		dout.WriteText(e.GetKey())
		dout.WriteText(e.GetValue().(string))
	}
}
func (this *EventPack) Read(din *io.DataInputX) {
	this.AbstractPack.Read(din)
	this.Level = din.ReadByte()
	this.Title = din.ReadText()
	this.Message = din.ReadText()

	sz := int(din.ReadByte())
	for i := 0; i < sz; i++ {
		key := din.ReadText()
		value := din.ReadText()
		this.Attr.Put(key, value)
	}

	//attr에서는 삭제한다. 데이터를 시리얼라이즈할때만 사용한다.
	val := this.Attr.Remove(ESCALATION_KEY)
	if val != nil {
		if val.(string) == "true" {
			this.Escalation = true
		} else {
			this.Escalation = false
		}
	}

	val = this.Attr.Remove(UUID_KEY)
	if val != nil {
		this.Uuid = val.(string)
	} else {
		val = ""
	}
	val = this.Attr.Remove(OTYPE_KEY)
	if val != nil {
		v, err := strconv.Atoi(val.(string))
		if err != nil {
			this.Otype = int32(v)
		} else {
			this.Otype = 0
		}
	}
	val = this.Attr.Remove(STATUS_KEY)
	if val != nil {
		v, err := strconv.Atoi(val.(string))
		if err != nil {
			this.Otype = int32(v)
		} else {
			this.Otype = 0
		}
	}
}

func (this *EventPack) SetUuid() {
	if this.Uuid != "" {
		return
	}
	this.Uuid = uuidutil.Generate()
}

func (this *EventPack) ToString() string {
	return fmt.Sprintf("ALERT %s %s %s %s %s ", this.AbstractPack.ToString(), string(this.Level), this.Title, this.Message, this.Attr.ToString())

}
