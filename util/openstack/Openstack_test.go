package openstack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//TODO Test 환경에 따라 결과값이 달라 상세한 TC 작성이 어려움

func TestOepnStack(t *testing.T) {
	zone, err := GetAvailabilityZone()
	if err != nil {
		assert.Equal(t, zone, "")
		return
	}
	assert.NotEqual(t, zone, "")

}
