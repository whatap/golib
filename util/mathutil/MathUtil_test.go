package mathutil

import (
	// "fmt"
	// "math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRound(t *testing.T) {
	// s := "hello world"
	// hO := int32(222957957)
	// hR := HashStr(s)
	// assert.Equal(t, hO, hR)
}

func TestRoundScale(t *testing.T) {
	// s := "hello world"
	// hO := int64(-281470736525980)
	// hR := Hash64Str(s)
	// assert.Equal(t, hO, hR)
}

func TestGetStandardDeviation(t *testing.T) {
	rawData := []float64{
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
	}

	sum, mean, sqrSum := func() (float64, float64, float64) {
		var sum float64 = 0
		var sqrSum float64 = 0

		for _, v := range rawData {
			sum += v
			sqrSum += v * v
		}
		return sum, sum / float64(len(rawData)), sqrSum
	}()

	n := len(rawData)
	stdDev := GetStandardDeviation(n, sum, sqrSum)
	assert.Equal(t, stdDev, 2.8722813232690143)

	assert.Equal(t, GetPct90(mean, stdDev), mean+stdDev*1.282)
	assert.Equal(t, GetPct95(mean, stdDev), mean+stdDev*1.645)
}

func TestGetPct90(t *testing.T) {
	rawData := []float64{
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
		11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
		21, 22, 23, 24, 25,

		//,26,27,28,29,30,31,32,33,34,35,36,37,38,39,40,41,42,43,44,45,46,47,48,49,50,51,52,53,54,55,56,57,58,59,60,61,62,63,64,65,66,67,68,69,70,71,72,73,74,75,76,77,78,79,80,81,82,83,84,85,86,87,88,89,90,91,92,93,94,95,96,97,98,99,100
	}

	sum, mean, sqrSum := func() (float64, float64, float64) {
		var sum float64 = 0
		var sqrSum float64 = 0

		for _, v := range rawData {
			sum += v
			sqrSum += v * v
		}
		return sum, sum / float64(len(rawData)), sqrSum
	}()

	n := len(rawData)
	stdDev := GetStandardDeviation(n, sum, sqrSum)
	// assert.Equal(t, stdDev, 2.8722813232690143)

	assert.Equal(t, GetPct90(mean, stdDev), mean+stdDev*1.282)
	assert.Equal(t, GetPct95(mean, stdDev), mean+stdDev*1.645)
}

func TestGetPct95(t *testing.T) {
	rawData := []float64{
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
	}

	sum, mean, sqrSum := func() (float64, float64, float64) {
		var sum float64 = 0
		var sqrSum float64 = 0

		for _, v := range rawData {
			sum += v
			sqrSum += v * v
		}
		return sum, sum / float64(len(rawData)), sqrSum
	}()

	n := len(rawData)
	stdDev := GetStandardDeviation(n, sum, sqrSum)
	assert.Equal(t, stdDev, 2.8722813232690143)

	assert.Equal(t, GetPct90(mean, stdDev), mean+stdDev*1.282)
	assert.Equal(t, GetPct95(mean, stdDev), mean+stdDev*1.645)

}
