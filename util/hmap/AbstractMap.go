package hmap


const (
	// 101
	DEFAULT_CAPACITY    = 101
	// 0.75
	DEFAULT_LOAD_FACTOR = 0.75
)

type PUT_MODE int
const (
	PUT_FORCE_FIRST PUT_MODE = 1
	PUT_FORCE_LAST  PUT_MODE = 2
	PUT_FIRST       PUT_MODE = 3
	PUT_LAST        PUT_MODE = 4
	
	// 1 Map Element Type
	ELEMENT_TYPE_KEYS = 1
	// 2 Map Element Type
	ELEMENT_TYPE_VALUES = 2
	// 3 Map Element Type
	ELEMENT_TYPE_ENTRIES = 3
	
)

type IntEnumer interface {
	HasMoreElements() bool
	NextInt() int32
}
type LongEnumer interface {
	HasMoreElements() bool
	NextLong() int64
}
type FloatEnumer interface {
	HasMoreElements() bool
	NextFloat() float32
}
type StringEnumer interface {
	HasMoreElements() bool
	NextString() string
}
type Enumeration interface {
	HasMoreElements() bool
	NextElement() interface{}
}
