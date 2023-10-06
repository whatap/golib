package step

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/whatap/golib/io"
	"github.com/whatap/golib/lang/value"
)

func TestGetStepType(t *testing.T) {
	p := NewMessageStepX()
	assert.Equal(t, int(STEP_MESSAGE_X), int(p.GetStepType()))
}

func TestWriteRead(t *testing.T) {
	now := 32
	p := NewMessageStepXWithStartTime(int32(now))

	p.Title = "Title1234"
	p.Desc = "Title1234-Desc1234235"
	p.Ctr = 1

	m := value.NewMapValue()
	m.Put("ttt", value.NewDecimalValue(64))
	m.Put("saaa", value.NewBoolValue(false))
	p.Attr = m

	out := io.NewDataOutputX()
	p.Write(out)
	data := out.ToByteArray()

	fmt.Println("data len=", len(data))

	p2 := NewMessageStepX()

	in := io.NewDataInputX(data)
	p2.Read(in)

	assert.Equal(t, p.GetStepType(), p2.GetStepType())
	assert.Equal(t, int32(now), p2.GetStartTime())
	assert.Equal(t, p.Title, p2.Title)
	assert.Equal(t, p.Desc, p2.Desc)
	assert.Equal(t, p.Ctr, p2.Ctr)
	assert.Equal(t, p.Attr.GetLong("ttt"), p2.Attr.GetLong("ttt"))
	assert.Equal(t, p.Attr.GetBool("saaa"), p2.Attr.GetBool("saaa"))
}

func TestCtrToJson(t *testing.T) {
	now := 32
	p := NewMessageStepXWithStartTime(int32(now))

	p.Title = "Title1234"
	p.Desc = "Title1234-Desc1234235"
	p.Ctr = 1

	fmt.Println(p.isA(SINGLE_LINE_DISPLAY))

	m := value.NewMapValue()
	m.Put("ttt", value.NewDecimalValue(64))
	m.Put("saaa", value.NewBoolValue(false))
	p.Attr = m

	jsonObject := p.CtrToJson()
	jsonData, err := json.Marshal(jsonObject)

	assert.Nil(t, err)
	assert.Equal(t, "{\"SINGLE_LINE_DISPLAY\":true}", string(jsonData))
}
