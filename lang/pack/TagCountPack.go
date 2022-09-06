package pack

import (
	"fmt"

	"github.com/whatap/golib/io"
	"github.com/whatap/golib/lang/value"
	"github.com/whatap/golib/util/hash"
)

type TagCountPack struct {
	AbstractPack
	Category string
	tagHash  int64
	Tags     *value.MapValue
	Data     *value.MapValue
}

func NewTagCountPack() *TagCountPack {
	p := new(TagCountPack)
	p.Tags = value.NewMapValue()
	p.Data = value.NewMapValue()
	return p
}

func (this *TagCountPack) GetPackType() int16 {
	return TAG_COUNT
}

func (this *TagCountPack) ToString() string {
	return fmt.Sprintln(this.AbstractPack.ToString(), " [Category=", this.Category, ", tags=", this.Tags, ", data=", this.Data, "]")
}

func (this *TagCountPack) GetTag(name string) string {
	return this.Tags.GetString(name)
}

func (this *TagCountPack) GetTagHash() int64 {
	return this.tagHash
}

func (this *TagCountPack) PutTag(name, val string) {
	this.Tags.PutString(name, val)
}

func (this *TagCountPack) Put(name string, v interface{}) {
	switch v.(type) {
	case value.Value:
		this.Data.Put(name, v.(value.Value))
	case int:
		this.Data.Put(name, value.NewDecimalValue(int64(v.(int))))
	case int16:
		this.Data.Put(name, value.NewDecimalValue(int64(v.(int16))))
	case int32:
		this.Data.Put(name, value.NewDecimalValue(int64(v.(int32))))
	case int64:
		this.Data.Put(name, value.NewDecimalValue(v.(int64)))
	case uint:
		this.Data.Put(name, value.NewDecimalValue(int64(v.(uint))))
	case uint32:
		this.Data.Put(name, value.NewDecimalValue(int64(v.(uint32))))
	case uint64:
		this.Data.Put(name, value.NewDecimalValue(int64(v.(uint64))))
	case float32:
		this.Data.Put(name, value.NewFloatValue(v.(float32)))
	case float64:
		this.Data.Put(name, value.NewDoubleValue(v.(float64)))
	case string:
		this.Data.Put(name, value.NewTextValue(v.(string)))
	default:
		panic(fmt.Sprintf("Panic, Not supported type %T. available type: value.Value, int, int32, int64, float32, float64, string ", v))
	}
}

func (this *TagCountPack) Get(name string) value.Value {
	return this.Data.Get(name)
}

func (this *TagCountPack) GetFloat(name string) float64 {
	val := this.Data.Get(name)
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

func (this *TagCountPack) GetLong(name string) int64 {
	val := this.Data.Get(name)
	if val == nil {
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

func (this *TagCountPack) Write(dout *io.DataOutputX) {
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
	value.WriteValue(dout, this.Data)
}

func (this *TagCountPack) Read(din *io.DataInputX) {
	this.AbstractPack.Read(din)
	//ver := din.ReadByte()
	din.ReadByte()
	this.Category = din.ReadText()
	this.tagHash = din.ReadDecimal()
	this.Tags = value.ReadValue(din).(*value.MapValue)
	this.Data = value.ReadValue(din).(*value.MapValue)
}

func (this *TagCountPack) IsEmpty() bool {
	return this.Data.IsEmpty()
}
func (this *TagCountPack) Size() int {
	return this.Data.Size()
}

func (this *TagCountPack) Clear() {
	this.Data.Clear()
}
