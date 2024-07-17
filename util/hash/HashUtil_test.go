package hash

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHash(t *testing.T) {
	s := "hello world"
	hO := int32(222957957)
	hR := HashStr(s)
	assert.Equal(t, hO, hR)
}

func TestHash64(t *testing.T) {
	s := "hello world"
	hO := int64(-281470736525980)
	hR := Hash64Str(s)
	assert.Equal(t, hO, hR)
}

func TestHash64v2(t *testing.T) {
	s := "hello world"
	hO := int64(-2739238572885903238)
	hR := GetLongHash(s)
	assert.Equal(t, hO, hR)

	hV2 := Hash64StrV2(s)
	fmt.Println("hR=", hR, ",HV2=", hV2)
	assert.Equal(t, hR, hV2)
}
