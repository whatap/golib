package pack

import (
	"fmt"

	"github.com/whatap/golib/io"
	"github.com/whatap/golib/lang/value"
	"github.com/whatap/golib/util/hash"
)

type TagLogPack struct {
	AbstractPack
	Category string
	tagHash  int64
	Tags     *value.MapValue
	Fields   *value.MapValue
}

func NewTagLogPack() *TagLogPack {
	p := new(TagLogPack)
	p.Tags = value.NewMapValue()
	p.Fields = value.NewMapValue()
	return p
}

func (this *TagLogPack) GetPackType() int16 {
	return TAG_LOG
}

func (this *TagLogPack) ToString() string {
	return fmt.Sprint(this.AbstractPack.ToString(), ",category=", this.Category, ",tags=", this.Tags.ToString(), ",fields=", this.Fields.ToString())
}

func (this *TagLogPack) GetTag(name string) string {
	return this.Tags.GetString(name)
}
func (this *TagLogPack) PutTag(name, val string) {
	this.Tags.PutString(name, val)
}
func (this *TagLogPack) PutTagLong(name string, val int64) {
	this.Tags.PutLong(name, val)
}
func (this *TagLogPack) Put(name string, v interface{}) {
	switch v.(type) {
	case value.Value:
		this.Fields.Put(name, v.(value.Value))
	case int:
		this.Fields.Put(name, value.NewDecimalValue(int64(v.(int))))
	case int32:
		this.Fields.Put(name, value.NewDecimalValue(int64(v.(int32))))
	case int64:
		this.Fields.Put(name, value.NewDecimalValue(v.(int64)))
	case float32:
		this.Fields.Put(name, value.NewFloatValue(v.(float32)))
	case float64:
		this.Fields.Put(name, value.NewDoubleValue(v.(float64)))
	case string:
		this.Fields.Put(name, value.NewTextValue(v.(string)))
	default:
		//panic(fmt.Sprintf("Panic, Not supported type %T. available type: value.Value, int, int32, int64, float32, float64, string ", v))
		this.Fields.Put(name, value.NewTextValue(fmt.Sprintf("%v", v)))
	}
}

func (this *TagLogPack) Get(name string) value.Value {
	return this.Fields.Get(name)
}

func (this *TagLogPack) GetFloat(name string) float64 {
	val := this.Fields.Get(name)
	if val == nil {
		return 0
	}

	switch val.GetValueType() {
	case value.VALUE_DOUBLE_SUMMARY, value.VALUE_LONG_SUMMARY:
		return val.(value.SummaryValue).DoubleAvg()
	case value.VALUE_DECIMAL:
		return float64(val.(*value.DecimalValue).Val)
	case value.VALUE_DECIMAL_INT:
		return float64(val.(*value.IntValue).Val)
	case value.VALUE_DECIMAL_LONG:
		return float64(val.(*value.LongValue).Val)
	case value.VALUE_FLOAT:
		return float64(val.(*value.FloatValue).Val)
	case value.VALUE_DOUBLE:
		return float64(val.(*value.DoubleValue).Val)
	default:
	}

	return 0
}

func (this *TagLogPack) GetLong(name string) int64 {
	val := this.Fields.Get(name)
	if val != nil {
		return 0
	}

	switch val.GetValueType() {
	case value.VALUE_DOUBLE_SUMMARY, value.VALUE_LONG_SUMMARY:
		return val.(value.SummaryValue).LongAvg()
	case value.VALUE_DECIMAL:
		return int64(val.(*value.DecimalValue).Val)
	case value.VALUE_DECIMAL_INT:
		return int64(val.(*value.IntValue).Val)
	case value.VALUE_DECIMAL_LONG:
		return int64(val.(*value.LongValue).Val)
	case value.VALUE_FLOAT:
		return int64(val.(*value.FloatValue).Val)
	case value.VALUE_DOUBLE:
		return int64(val.(*value.DoubleValue).Val)
	default:

		//	default:
		//						if (val instanceof Number) {
		//							return ((Number) val).longValue();
		//						}

	}
	return 0
}

func (this *TagLogPack) Write(dout *io.DataOutputX) {
	this.AbstractPack.Write(dout)
	dout.WriteByte(0)
	dout.WriteText(this.Category)
	if this.tagHash == 0 && this.Tags.Size() > 0 {
		tagIO := io.NewDataOutputX()
		value.WriteValue(tagIO, this.Tags)
		tagBytes := tagIO.ToByteArray()
		this.tagHash = hash.Hash64(tagBytes)
		dout.WriteDecimal(this.tagHash)
		dout.WriteBytes(tagBytes)
	} else {
		dout.WriteDecimal(this.tagHash)
		value.WriteValue(dout, this.Tags)
	}
	value.WriteValue(dout, this.Fields)
}

func (this *TagLogPack) Read(din *io.DataInputX) {
	this.AbstractPack.Read(din)
	//ver := din.ReadByte()
	din.ReadByte()
	this.Category = din.ReadText()
	this.Tags = value.ReadValue(din).(*value.MapValue)
	this.Fields = value.ReadValue(din).(*value.MapValue)
}

func (this *TagLogPack) IsEmpty() bool {
	return this.Fields.IsEmpty()
}
func (this *TagLogPack) Size() int {
	return this.Fields.Size()
}

func (this *TagLogPack) Clear() {
	this.Fields.Clear()
}
