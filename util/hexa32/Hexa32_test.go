package hexa32

import (
	// "fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToSting32(t *testing.T) {
	s := int64(-743752992412427445)
	v := ToString32(s)
	s1 := "zkkincvom0p5l"
	assert.Equal(t, s1, v)

}

func TestToLong32(t *testing.T) {
	s := "zkkincvom0p5l"
	s1 := int64(-743752992412427445)
	v := ToLong32(s)
	assert.Equal(t, s1, v)
}
