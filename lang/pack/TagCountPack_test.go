package pack

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/whatap/golib/io"
)

func TestTagCoutPackDefault(t *testing.T) {
	p := NewTagCountPack()
	assert.NotNil(t, p)
	assert.Equal(t, p.GetPackType(), int16(TAG_COUNT))

	p.PutTag("TAGTEST", "tagtest")

	assert.Equal(t, p.GetTag("TAGTEST"), "tagtest")

	p.Put("TESTLONG1", int(1))
	p.Put("TESTLONG2", int16(2))
	p.Put("TESTLONG3", int32(3))
	p.Put("TESTLONG4", int64(4))
	p.Put("TESTLONG5", uint(5))
	p.Put("TESTLONG6", uint32(6))
	p.Put("TESTLONG7", uint64(7))
	p.Put("TESTFLOAT1", float32(1.5))
	p.Put("TESTFLOAT2", float64(2.5))

	assert.Equal(t, p.GetLong("TESTLONG1"), int64(1))
	assert.Equal(t, p.GetLong("TESTLONG2"), int64(2))
	assert.Equal(t, p.GetLong("TESTLONG3"), int64(3))
	assert.Equal(t, p.GetLong("TESTLONG4"), int64(4))
	assert.Equal(t, p.GetLong("TESTLONG5"), int64(5))
	assert.Equal(t, p.GetLong("TESTLONG6"), int64(6))
	assert.Equal(t, p.GetLong("TESTLONG7"), int64(7))
	assert.Equal(t, p.GetFloat("TESTFLOAT1"), float64(1.5))
	assert.Equal(t, p.GetFloat("TESTFLOAT2"), float64(2.5))

	dout := io.NewDataOutputX()
	p.Write(dout)

	din := io.NewDataInputX(dout.ToByteArray())
	p2 := NewTagCountPack()
	p2.Read(din)

	assert.Equal(t, p.GetTagHash(), p2.GetTagHash())
	assert.True(t, p.Tags.Equals(p2.Tags))

}
