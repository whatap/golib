package lang

const (
	EVENT_STATUS_ON          = 0x0001
	EVENT_STATUS_OFF         = 0x0002
	EVENT_STATUS_DISABLED    = 0x0004
	EVENT_STATUS_ACKNOWLEDGE = 0x0008
)

var EventStatusName = map[int]string{0x0001: "ON", 0x0002: "OFF", 0x0004: "DISABLED", 0x0008: "ACKNOWLEDGE"}
var EventStatusValue = map[string]int{"ON": 0x0001, "OFF": 0x0002, "DISABLED": 0x0004, "ACKNOWLEDGE": 0x0008}
