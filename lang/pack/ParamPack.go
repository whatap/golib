package pack

import (
	"bytes"
	"fmt"

	"github.com/whatap/golib/io"
	val "github.com/whatap/golib/lang/value"
	"github.com/whatap/golib/util/hmap"
)

type ParamPack struct {
	AbstractPack
	Id       int32
	table    *hmap.StringKeyLinkedMap
	Request  int64
	Response int64
}

func NewParamPack() *ParamPack {
	p := new(ParamPack)
	p.table = hmap.NewStringKeyLinkedMap()
	return p
}

func (this *ParamPack) Get(key string) val.Value {
	o := this.table.Get(key)
	if o == nil {
		return val.NewNullValue()
	}
	return o.(val.Value)
}

func (this *ParamPack) GetMap(key string) *val.MapValue {
	o := this.table.Get(key)
	if o == nil {
		return nil
	}
	return o.(*val.MapValue)
}

func (this *ParamPack) GetString(key string) string {
	o := this.Get(key)
	if o.GetValueType() == val.VALUE_TEXT {
		t := o.(*val.TextValue)
		return t.Val
	}
	return ""
}

func (this *ParamPack) GetLong(key string) int64 {
	o := this.Get(key)
	if o.GetValueType() == val.VALUE_DECIMAL {
		v := o.(*val.DecimalValue)
		return v.Val
	}
	return 0
}

func (this *ParamPack) PutString(key string, value string) {
	this.Put(key, val.NewTextValue(value))
}
func (this *ParamPack) PutLong(key string, value int64) {
	this.Put(key, val.NewDecimalValue(value))
}
func (this *ParamPack) Put(key string, value val.Value) {
	this.table.Put(key, value)
}
func (this *ParamPack) Clear() {
	this.table.Clear()
}
func (this *ParamPack) Remove(key string) {
	this.table.Remove(key)
}
func (this *ParamPack) Size() {
	this.table.Size()
}
func (this *ParamPack) Keys() hmap.StringEnumer {
	return this.table.Keys()
}

func (this *ParamPack) GetPackType() int16 {
	return PACK_PARAMETER
}
func (this *ParamPack) Write(dout *io.DataOutputX) {
	this.AbstractPack.Write(dout)
	dout.WriteInt(this.Id)
	dout.WriteDecimal(this.Request)
	dout.WriteDecimal(this.Response)
	dout.WriteDecimal(int64(this.table.Size()))

	keys := this.Keys()
	for keys.HasMoreElements() {
		key := keys.NextString()
		value := this.table.Get(key).(val.Value)

		dout.WriteText(key)
		val.WriteValue(dout, value)
	}
}
func (this *ParamPack) Read(din *io.DataInputX) {
	this.AbstractPack.Read(din)
	this.Id = din.ReadInt()
	this.Request = din.ReadDecimal()
	this.Response = din.ReadDecimal()
	count := int(din.ReadDecimal())
	for t := 0; t < count; t++ {
		key := din.ReadText()
		value := val.ReadValue(din)
		this.table.Put(key, value)
	}
}
func (this *ParamPack) SetMapValue(mapValue *val.MapValue) *ParamPack {
	if mapValue == nil {
		return this
	}
	keys := mapValue.Keys()
	for keys.HasMoreElements() {
		key := keys.NextString()
		value := mapValue.Get(key)
		this.table.Put(key, value)
	}
	return this
}
func (this *ParamPack) ToResponse() *ParamPack {
	if this.Request == 0 {
		return this
	}
	this.Response = this.Request
	this.Request = 0
	return this
}
func (this *ParamPack) ToString() string {
	var buffer bytes.Buffer
	buffer.WriteString("ParamPack")
	buffer.WriteString(fmt.Sprintf("%s", this.AbstractPack.ToString()))
	buffer.WriteString("\nId=")
	buffer.WriteString(fmt.Sprintf("%d", this.Id))
	buffer.WriteString("\nRequest=")
	buffer.WriteString(fmt.Sprintf("%d", this.Request))
	buffer.WriteString("\nResponse=")
	buffer.WriteString(fmt.Sprintf("%d", this.Response))
	buffer.WriteString("\n")
	buffer.WriteString(this.table.ToString())
	return buffer.String()
}
