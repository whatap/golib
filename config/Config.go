//github.com/whatap/golib/config
package config

// "log"

const ()

//go:generate mockery --name Config --case underscore --inpackage
type Config interface {
	ApplyDefault()
	GetValue(key string) string
	GetValueDef(key, def string) string
	GetBoolean(key string, def bool) bool
	GetInt(key string, def int) int32
	GetIntSet(key, def, deli string) []int32
	GetLong(key string, def int64) int64
	GetStringArray(key string, def string, deli string) []string
	GetStringHashSet(key, def, deli string) []int32
	GetStringHashCodeSet(key, def, deli string) []int32
	GetFloat(key string, def float32) float32
	SetValues(v map[string]string)
	ToString() string
	String() string
}
