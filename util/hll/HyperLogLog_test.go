package hll

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHllCardinality(t *testing.T) {
	assert := assert.New(t)

	hll := NewHyperLogLogDefault()
	hll.Offer(1)
	hll.Offer(2)
	hll.Offer(3)
	hll.Offer(4)
	hll.Offer(5)
	hll.Offer(6)
	hll.Offer(7)
	assert.Equal(int(hll.Cardinality()), 7)

	hll.Offer(2)
	hll.Offer(3)
	hll.Offer(4)
	hll.Offer(5)
	hll.Offer(6)
	hll.Offer(7)
	assert.Equal(int(hll.Cardinality()), 7)

	hll.Offer(8)
	hll.Offer(9)
	hll.Offer(10)
	assert.Equal(int(hll.Cardinality()), 10)
}

func TestHllMerge(t *testing.T) {
	assert := assert.New(t)

	hll := NewHyperLogLogDefault()
	hll2 := NewHyperLogLogDefault()

	hll.Offer(1)
	hll.Offer(2)
	hll.Offer(3)
	hll.Offer(4)
	hll.Offer(5)

	hll2.Offer(3)
	hll2.Offer(4)
	hll2.Offer(5)
	hll2.Offer(6)
	hll2.Offer(7)

	mergedHll := hll.Merge(hll2)
	assert.Equal(int(mergedHll.Cardinality()), 7)

	hll = NewHyperLogLogDefault()
	hll.Offer(8)
	hll.Offer(9)
	hll.Offer(10)

	mergedHll = mergedHll.Merge(hll)
	assert.Equal(int(mergedHll.Cardinality()), 10)

}

func TestHllByte(t *testing.T) {
	assert := assert.New(t)

	hll := NewHyperLogLogDefault()
	hll.Offer(1)
	hll.Offer(2)
	hll.Offer(3)
	hll.Offer(4)
	hll.Offer(5)
	hll.Offer(6)
	hll.Offer(7)

	buildHll := BuildHyperLogLog(hll.GetBytes())
	assert.Equal(int(buildHll.Cardinality()), 7)

}

func TestMain(m *testing.M) {
	m.Run()
}
