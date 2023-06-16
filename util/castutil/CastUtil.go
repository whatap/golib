package castutil

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/whatap/golib/lang/value"
)

func CInt(val interface{}) (rt int32) {
	defer func() {
		if r := recover(); r != nil {
			//logutil.Println("Recover CaseTuil CInt ", r)
			rt = 0
		}
	}()

	if val == nil {
		return 0
	} else {
		switch val.(type) {
		case string:
			v, err := strconv.Atoi(val.(string))
			if err != nil {
				return 0
			} else {
				return int32(v)
			}

		default:
			return int32(val.(int64))
		}
	}
}

func CInteger(val interface{}) int32 {
	return CInt(val)
}

func CLong(val interface{}) (rt int64) {
	defer func() {
		if r := recover(); r != nil {
			// logutil.Println("Recover CaseTuil CLong ", r)
			rt = 0
		}
	}()

	if val == nil {
		return 0
	} else {
		switch val.(type) {
		case string:
			v, err := strconv.Atoi(val.(string))
			if err != nil {
				return 0
			} else {
				return int64(v)
			}

		default:
			return int64(val.(int64))
		}
	}

}

func CFloat(val interface{}) (rt float32) {
	defer func() {
		if r := recover(); r != nil {
			// logutil.Println("Recover CaseTuil CFloat ", r)
			rt = 0
		}
	}()

	if val == nil {
		return 0
	} else {
		switch val.(type) {
		case string:
			v, err := strconv.ParseFloat(val.(string), 32)

			if err != nil {
				return 0
			} else {
				return float32(v)
			}

		default:
			return float32(val.(float64))
		}
	}

}

func CDouble(val interface{}) (rt float64) {
	defer func() {
		if r := recover(); r != nil {
			// logutil.Println("Recover CaseTuil CDouble ", r)
			rt = 0
		}
	}()

	if val == nil {
		return 0
	} else {
		switch val.(type) {
		case string:
			v, err := strconv.ParseFloat(val.(string), 64)
			if err != nil {
				return 0
			} else {
				return float64(v)
			}

		default:
			return float64(val.(float64))
		}
	}

}

func ToString(val interface{}) string {
	return fmt.Sprintln(val)
}

func CString(val interface{}) (rt string) {
	defer func() {
		if r := recover(); r != nil {
			// logutil.Println("CastUtil", "Recover ", r)
			rt = ""
		}
	}()
	if val == nil {
		return ""
	}
	switch val.(type) {
	case string:
		return val.(string)
	case float32, float64:
		return strconv.FormatFloat(val.(float64), 'f', 7, 64)
	default:
		return fmt.Sprintf("%s", val)
	}
}

func CBool(val interface{}) (rt bool) {
	defer func() {
		if r := recover(); r != nil {
			// logutil.Println("CastUtil", "Recover ", r)
			rt = false
		}
	}()

	if val == nil {
		return false
	}
	switch val.(type) {
	case bool:
		return val.(bool)
	case *value.BoolValue:
		return val.(*value.BoolValue).Val
	case string:
		return strings.ToLower(val.(string)) == "true"
	default:
		return false
	}
}
