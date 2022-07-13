package value

import (
	"github.com/whatap/golib/io"
)

type Value interface {
	GetValueType() byte
	Equals(o Value) bool
	CompareTo(o Value) int
	Write(out *io.DataOutputX)
	Read(in *io.DataInputX)
}

const (
	VALUE_NULL         = 0
	VALUE_BOOLEAN      = 10
	VALUE_DECIMAL      = 20
	VALUE_DECIMAL_INT  = 21
	VALUE_DECIMAL_LONG = 22
	VALUE_FLOAT        = 30
	VALUE_DOUBLE       = 40

	VALUE_DOUBLE_SUMMARY = 45
	VALUE_LONG_SUMMARY   = 46
	FLOAT_SUMMARY        = 47

	VALUE_TEXT      = 50
	VALUE_TEXT_HASH = 51
	VALUE_BLOB      = 60
	VALUE_IP4ADDR   = 61

	VALUE_LIST  = 70
	ARRAY_INT   = 71
	ARRAY_FLOAT = 72
	ARRAY_TEXT  = 73
	ARRAY_LONG  = 74

	VALUE_MAP     = 80
	INT_VALUE_MAP = 81
)

func CreateValue(code byte) Value {
	switch code {
	case VALUE_NULL:
		return NewNullValue()
	case VALUE_BOOLEAN:
		return NewBoolValue(false)
	case VALUE_DECIMAL:
		return NewDecimalValue(0)
	case VALUE_DECIMAL_INT:
		return NewIntValue(0)
	case VALUE_DECIMAL_LONG:
		return NewLongValue(0)
	case VALUE_FLOAT:
		return NewFloatValue(0.0)
	case VALUE_DOUBLE:
		return NewDoubleValue(0.0)

	case VALUE_DOUBLE_SUMMARY:
		return NewDoubleSummary()
	case VALUE_LONG_SUMMARY:
		return NewLongSummary()
		//	case FLOAT_SUMMARY:
		//		return NewFloatSummary()

	case VALUE_TEXT:
		return NewTextValue("")
	case VALUE_TEXT_HASH:
		return NewTextHashValue(0)
	case VALUE_BLOB:
		return NewBlobValue([]byte{})
	case VALUE_IP4ADDR:
		return NewIP4Value(nil)

	case VALUE_LIST:
		return NewListValue(nil)
	case ARRAY_INT:
		return NewIntArray([]int32{})
	case ARRAY_FLOAT:
		return NewFloatArray([]float32{})
	case ARRAY_TEXT:
		return NewTextArray([]string{})
	case ARRAY_LONG:
		return NewLongArray([]int64{})

	case VALUE_MAP:
		return NewMapValue()
	case INT_VALUE_MAP:
		return NewIntMapValue()

	}
	panic("unknown value : " + string(code))
}

func WriteValue(out *io.DataOutputX, val Value) *io.DataOutputX {
	out.WriteByte(val.GetValueType())
	val.Write(out)
	return out
}

func ReadValue(in *io.DataInputX) Value {
	t := in.ReadByte()
	v := CreateValue(byte(t))
	v.Read(in)
	return v
}

func WriteMapValue(out *io.DataOutputX, val *MapValue) *io.DataOutputX {
	out.WriteByte(val.GetValueType())
	val.Write(out)
	return out
}

func ReadMapValue(in *io.DataInputX) (ret *MapValue) {
	t := in.ReadByte()
	// fmt.Println("Value.ReadMapValue step -1 ", t, t == VALUE_MAP)
	if t == VALUE_MAP {
		// fmt.Println("Value.ReadMapValue step -2")
		ret = NewMapValue()
		ret.Read(in)
	}
	return
}
