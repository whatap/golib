package topology

import (
	// "fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateLINK(t *testing.T) {
	k := NewLINK()

	k.IP = make([]byte, 16)
	k.Port = 123
	// fmt.Println(k)
	assert.Equal(t, "::", k.IP.String())
	assert.Equal(t, 123, k.Port)

	k1 := CreateLINK("0:0:0:0:0:0:0:0", 123)
	// fmt.Println(k1)
	assert.Equal(t, "0.0.0.0", k1.IP.String())
	assert.Equal(t, 123, k1.Port)
}
