package mathutil

import (
	"fmt"
	"math"
)

func Round(v float64) float64 {
	return float64(int64(v))
}

func RoundScale(value float64, scale int) float64 {
	if scale == 0 {
		return float64(int64(value))
	}

	var r int64 = Scale(scale)
	var v int64 = int64(value * float64(r))
	return float64(v) / float64(r)
}

func Scale(n int) int64 {
	var r int64 = 0
	switch n {
	case 1:
		r = 10
	case 2:
		r = 100
	case 3:
		r = 1000
	default:
		r = 10000
	}
	return r
}

func RoundString(value float64, scale int) string {
	if scale <= 0 {
		return fmt.Sprintf("%d", int64(value))
	}
	return fmt.Sprintf("%f", RoundScale(value, scale))
}

func Round2(value float64) float64 {
	return RoundScale(value, 2)
}

func Round4(value float64) float64 {
	return RoundScale(value, 4)
}

func GetStandardDeviation(count int, timeSum float64, timeSqrSum float64) float64 {
	if count == 0 {
		return 0
	}
	if timeSqrSum == 0 {
		return 0
	}

	// 제곱의 평균 - 평균의 제곱
	avg := timeSum / float64(count)
	variation := (timeSqrSum / float64(count)) - (avg * avg)
	ret := math.Sqrt(math.Abs(variation))
	if math.IsNaN(ret) {
		return 0
	} else {
		return ret
	}
}

func GetPct90(avg, stdDev float64) float64 {
	return (avg + stdDev*1.282)
}
func GetPct95(avg, stdDev float64) float64 {
	return (avg + stdDev*1.645)
}
