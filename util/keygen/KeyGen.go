//github.com/whatap/golib/util/keygen
package keygen

import (
	"math"
	"math/rand"
	"time"
)

// Import 에서 초기화
func init() {
	Seed()
}

func Seed() {
	rand.Seed(time.Now().UnixNano())
}

func AddSeed(v ...interface{}) {
	seed := time.Now().UnixNano()
	for _, it := range v {
		switch it.(type) {
		case int8, int16, int32, int64:
		case uint8, uint16, uint32, uint64:
		case float32, float64:
			seed += it.(int64)
		}
	}
	rand.Seed(seed)
}
func SetSeed(i int64) {
	rand.Seed(i)
}
func Next() int64 {
	v := rand.NormFloat64()
	return int64(math.Float64bits(v))
}
func RandInt(i int32) int32 {
	return rand.Int31n(i)
}
func RandLong(i int64) int64 {
	return rand.Int63n(i)
}
