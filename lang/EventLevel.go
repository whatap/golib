package lang

const (
	EVENT_LEVEL_CRITICAL = 30
	EVENT_LEVEL_WARNING  = 20
	EVENT_LEVEL_INFO     = 10
	EVENT_LEVEL_NONE     = 0
)

var EventLevelName = map[int]string{30: "CRITICAL", 20: "WARNING", 10: "INFO", 0: "NONE"}
var EventLevelValue = map[string]int{"CRITICAL": 30, "WARNING": 20, "INFO": 10, "NONE": 0}
