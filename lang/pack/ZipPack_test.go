package pack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeZipPack(t *testing.T) {
	zp := NewZipPack()

	pList := make([]Pack, 0)

	for i := 0; i < 10; i++ {
		tp := NewTagCountPack()
		tp.PutTag("tag", "Value")
		tp.Put("field", i)

		pList = append(pList, tp)
	}

	zp = zp.SetRecords(pList)

	tpList := zp.GetRecords()
	for i, tp := range tpList {
		assert.Equal(t, tp.(*TagCountPack).GetLong("field"), (int64)(i))
	}
}
