//github.com/whatap/golib/io
package io

import (
	"bytes"
	"math"
)

const (
	INT3_MIN_VALUE  = -8388608 /*0xff800000*/
	INT3_MAX_VALUE  = 0x007fffff
	LONG5_MIN_VALUE = -549755813888 /*0xffffff8000000000*/
	LONG5_MAX_VALUE = 0x0000007fffffffff
)

type DataOutputX struct {
	buffer  bytes.Buffer
	written int
}

func NewDataOutputX() *DataOutputX {
	in := new(DataOutputX)
	return in
}

func (out *DataOutputX) ToByteArray() []byte {
	return out.buffer.Bytes()
}

func (out *DataOutputX) WriteIntBytes(b []byte) *DataOutputX {
	if b == nil || len(b) == 0 {
		out.WriteInt(0)
	} else {
		out.WriteInt(int32(len(b)))
		out.WriteBytes(b)
	}
	return out
}

func (out *DataOutputX) WriteShortBytes(b []byte) *DataOutputX {
	if b == nil || len(b) == 0 {
		out.WriteShort(0)
	} else {
		out.WriteShort(int16(len(b)))
		out.WriteBytes(b)
	}
	return out
}

func (out *DataOutputX) WriteBlob(value []byte) *DataOutputX {
	if value == nil || len(value) == 0 {
		out.WriteByte(0)
	} else {
		sz := len(value)
		if sz <= 253 {
			out.WriteByte(byte(sz))
			out.WriteBytes(value)
		} else if sz <= 65535 {
			buff := []byte{255, 0, 0}
			out.WriteBytes(SetBytesShort(buff, 1, int16(sz)))
			out.WriteBytes(value)
		} else {
			buff := []byte{254, 0, 0, 0, 0}
			out.WriteBytes(SetBytesInt(buff, 1, int32(sz)))
			out.WriteBytes(value)
		}
	}
	return out
}
func (out *DataOutputX) WriteDecimal(v int64) *DataOutputX {

	switch {
	case v == 0:
		out.WriteByte(0)
	case math.MinInt8 <= v && v <= math.MaxInt8:
		b := []byte{0, 0}
		b[0] = 1
		b[1] = byte(v)
		out.WriteBytes(b)
	case math.MinInt16 <= v && v <= math.MaxInt16:
		b := []byte{0, 0, 0}
		b[0] = 2
		SetBytesShort(b, 1, int16(v))
		out.WriteBytes(b)
	case INT3_MIN_VALUE <= v && v <= INT3_MAX_VALUE:
		b := []byte{0, 0, 0, 0}
		b[0] = 3
		out.WriteBytes(SetBytesInt3(b, 1, int32(v)))
	case math.MinInt32 <= v && v <= math.MaxInt32:
		b := []byte{0, 0, 0, 0, 0}
		b[0] = 4
		out.WriteBytes(SetBytesInt(b, 1, int32(v)))
	case LONG5_MIN_VALUE <= v && v <= LONG5_MAX_VALUE:
		b := []byte{0, 0, 0, 0, 0, 0}
		b[0] = 5
		out.WriteBytes(SetBytesLong5(b, 1, v))
	case math.MinInt64 <= v && v <= math.MaxInt64:
		b := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0}
		b[0] = 8
		out.WriteBytes(SetBytesLong(b, 1, v))
	}
	return out
}
func (out *DataOutputX) WriteText(s string) *DataOutputX {
	if s == "" {
		out.WriteByte(0)
	} else {
		out.WriteBlob([]byte(s))
	}
	return out
}

func (out *DataOutputX) WriteBytes(b []byte) *DataOutputX {
	out.written += len(b)
	out.buffer.Write(b)
	return out
}
func (out *DataOutputX) Write(b []byte, off int, sz int) *DataOutputX {
	out.written += sz
	out.buffer.Write(b[off : off+sz])
	return out
}

func (out *DataOutputX) WriteBool(b bool) *DataOutputX {
	out.WriteBytes(ToBytesBool(b))
	return out
}
func (out *DataOutputX) WriteByte(b byte) *DataOutputX {
	out.written++
	out.buffer.WriteByte(b)
	return out
}
func (out *DataOutputX) WriteShort(b int16) *DataOutputX {
	out.WriteBytes(ToBytesShort(b))
	return out
}
func (out *DataOutputX) WriteUShort(b uint16) *DataOutputX {
	out.WriteBytes(ToBytesUShort(b))
	return out
}
func (out *DataOutputX) WriteInt3(b int32) *DataOutputX {
	out.WriteBytes(ToBytesInt3(b))
	return out
}
func (out *DataOutputX) WriteInt(b int32) *DataOutputX {
	out.WriteBytes(ToBytesInt(b))
	return out
}
func (out *DataOutputX) WriteLong5(b int64) *DataOutputX {
	out.WriteBytes(ToBytesLong5(b))
	return out
}
func (out *DataOutputX) WriteLong(b int64) *DataOutputX {
	out.WriteBytes(ToBytesLong(b))
	return out
}
func (out *DataOutputX) WriteFloat(b float32) *DataOutputX {
	out.WriteBytes(ToBytesFloat(b))
	return out
}
func (out *DataOutputX) WriteDouble(b float64) *DataOutputX {
	out.WriteBytes(ToBytesDouble(b))
	return out
}
func (out *DataOutputX) WriteShortArray(v []int16) {
	if v == nil {
		out.WriteShort(0)
	} else {
		sz := len(v)
		out.WriteShort(int16(sz))
		for i := 0; i < sz; i++ {
			out.WriteShort(v[i])
		}
	}
}
func (out *DataOutputX) WriteIntArray(v []int32) {
	if v == nil {
		out.WriteShort(0)
	} else {
		sz := len(v)
		out.WriteShort(int16(sz))
		for i := 0; i < sz; i++ {
			out.WriteInt(v[i])
		}
	}
}
func (out *DataOutputX) WriteLongArray(v []int64) {
	if v == nil {
		out.WriteShort(0)
	} else {
		sz := len(v)
		out.WriteShort(int16(sz))
		for i := 0; i < sz; i++ {
			out.WriteLong(v[i])
		}
	}
}
func (out *DataOutputX) WriteFloatArray(v []float32) {
	if v == nil {
		out.WriteShort(0)
	} else {
		sz := len(v)
		out.WriteShort(int16(sz))
		for i := 0; i < sz; i++ {
			out.WriteFloat(v[i])
		}
	}
}
func (out *DataOutputX) WriteDoubleArray(v []float64) {
	if v == nil {
		out.WriteShort(0)
	} else {
		sz := len(v)
		out.WriteShort(int16(sz))
		for i := 0; i < sz; i++ {
			out.WriteDouble(v[i])
		}
	}
}
func (out *DataOutputX) WriteTextArray(v []string) {
	if v == nil {
		out.WriteShort(0)
	} else {
		sz := len(v)
		out.WriteShort(int16(sz))
		for i := 0; i < sz; i++ {
			out.WriteText(v[i])
		}
	}
}
func (out *DataOutputX) WriteTextShortLength(v string) {
	if v == "" {
		out.WriteShort(0)
	} else {
		b := []byte(v)
		out.WriteShort(int16(len(b)))
		out.WriteBytes(b)
	}
}

func (out *DataOutputX) Size() int {
	return out.written
}

// Must do after write pack type and pack data.
func (out *DataOutputX) WriteHeader(netSrc, netSrcVer byte, pcode, licenseHash int64) {
	b := out.buffer.Bytes()
	t := make([]byte, len(b))
	copy(t, b)
	out.buffer.Reset()

	out.WriteByte(netSrc)
	out.WriteByte(netSrcVer)
	out.WriteLong(pcode)
	out.WriteLong(licenseHash)
	out.WriteIntBytes(t)
}

// Must do after write pack type and pack data.
func (out *DataOutputX) WriteOneWayHeader(netSrc, netSrcVer byte, pcode, licenseHash int64) {
	b := out.buffer.Bytes()
	t := make([]byte, len(b))
	copy(t, b)
	out.buffer.Reset()

	out.WriteByte(netSrc)
	out.WriteByte(netSrcVer)
	out.WriteLong(pcode)
	out.WriteLong(licenseHash)
	out.WriteIntBytes(t)
}

// Must do after write pack type and pack data.
func (out *DataOutputX) WriteSecureHeader(netSrc, netSrcVer byte, pcode int64, oid, transferKey int32) {
	b := out.buffer.Bytes()
	t := make([]byte, len(b))
	copy(t, b)
	out.buffer.Reset()

	out.WriteByte(netSrc)
	out.WriteByte(netSrcVer)
	out.WriteLong(pcode)
	out.WriteInt(oid)
	out.WriteInt(transferKey)
	out.WriteIntBytes(t)
}

func ToBytesBool(b bool) []byte {
	if b {
		return []byte{1}
	} else {
		return []byte{0}
	}
}
func SetBytesBool(buf []byte, off int, b bool) []byte {
	if b {
		buf[off] = 1
	} else {
		buf[off] = 0
	}
	return buf
}
func ToBytesShort(v int16) []byte {
	buf := []byte{0, 0}
	buf[0] = byte(v >> 8)
	buf[1] = byte(v >> 0)
	return buf
}

func ToBytesUShort(v uint16) []byte {
	buf := []byte{0, 0}
	buf[0] = byte(v >> 8)
	buf[1] = byte(v >> 0)
	return buf
}
func SetBytesShort(buf []byte, off int, v int16) []byte {
	buf[off] = byte(v >> 8)
	buf[off+1] = byte(v >> 0)
	return buf
}
func ToBytesInt(v int32) []byte {
	buf := []byte{0, 0, 0, 0}
	buf[0] = byte(v >> 24)
	buf[1] = byte(v >> 16)
	buf[2] = byte(v >> 8)
	buf[3] = byte(v >> 0)
	return buf
}
func SetBytesInt(buf []byte, off int, v int32) []byte {
	buf[off] = byte(v >> 24)
	buf[off+1] = byte(v >> 16)
	buf[off+2] = byte(v >> 8)
	buf[off+3] = byte(v >> 0)
	return buf
}

func ToBytesInt3(v int32) []byte {
	buf := []byte{0, 0, 0}
	buf[0] = byte(v >> 16)
	buf[1] = byte(v >> 8)
	buf[2] = byte(v >> 0)
	return buf
}
func SetBytesInt3(buf []byte, off int, v int32) []byte {
	buf[off] = byte(v >> 16)
	buf[off+1] = byte(v >> 8)
	buf[off+2] = byte(v >> 0)
	return buf
}
func ToBytesLong(v int64) []byte {
	buf := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	buf[0] = byte(v >> 56)
	buf[1] = byte(v >> 48)
	buf[2] = byte(v >> 40)
	buf[3] = byte(v >> 32)
	buf[4] = byte(v >> 24)
	buf[5] = byte(v >> 16)
	buf[6] = byte(v >> 8)
	buf[7] = byte(v >> 0)
	return buf
}

func SetBytesLong(buf []byte, off int, v int64) []byte {
	buf[off] = byte(v >> 56)
	buf[off+1] = byte(v >> 48)
	buf[off+2] = byte(v >> 40)
	buf[off+3] = byte(v >> 32)
	buf[off+4] = byte(v >> 24)
	buf[off+5] = byte(v >> 16)
	buf[off+6] = byte(v >> 8)
	buf[off+7] = byte(v >> 0)
	return buf
}
func ToBytesLong5(v int64) []byte {
	buf := []byte{0, 0, 0, 0, 0}
	buf[0] = byte(v >> 32)
	buf[1] = byte(v >> 24)
	buf[2] = byte(v >> 16)
	buf[3] = byte(v >> 8)
	buf[4] = byte(v >> 0)
	return buf
}

func SetBytesLong5(buf []byte, off int, v int64) []byte {
	buf[off] = byte(v >> 32)
	buf[off+1] = byte(v >> 24)
	buf[off+2] = byte(v >> 16)
	buf[off+3] = byte(v >> 8)
	buf[off+4] = byte(v >> 0)
	return buf
}
func ToBytesFloat(v float32) []byte {
	return ToBytesInt(int32(math.Float32bits(v)))
}

func SetBytesFloat(buf []byte, off int, v float32) []byte {
	return SetBytesInt(buf, off, int32(math.Float32bits(v)))
}
func ToBytesDouble(v float64) []byte {
	return ToBytesLong(int64(math.Float64bits(v)))
}

func SetBytesDouble(buf []byte, off int, v float64) []byte {
	return SetBytesLong(buf, off, int64(math.Float64bits(v)))
}
func SetBytes(dest []byte, pos int, src []byte) []byte {
	copy(dest[pos:pos+len(src)], src)
	return dest
}
