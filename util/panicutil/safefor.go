package panicutil

import (
	"sync"
	"time"
)

var onofflookup map[string]bool
var AllOff bool = false
var OffSleepTime int32 = 10
var PerfMap map[string]int32 = map[string]int32{}
var SendPerfMapEnabled bool = false
var perfMapLock = sync.Mutex{}

//ResetPerfMap ResetPerfMap
func ResetPerfMap() map[string]int32 {
	perfMapLock.Lock()
	m := PerfMap

	PerfMap = map[string]int32{}
	perfMapLock.Unlock()
	return m

}

//SafeFor SafeFor
func SafeFor(name string, callback func()) {
	for {
		start := time.Now()

		if AllOff {
			time.Sleep(time.Duration(OffSleepTime) * time.Millisecond)
			return
		}
		if on, ok := onofflookup[name]; ok {
			if !on {
				time.Sleep(time.Duration(OffSleepTime) * time.Millisecond)
			} else {
				callback()
				perfMapLock.Lock()
				PerfMap[name] = int32(time.Since(start).Nanoseconds() / 1000)
				perfMapLock.Unlock()
			}
		} else {
			callback()
			perfMapLock.Lock()
			PerfMap[name] = int32(time.Since(start).Nanoseconds() / 1000)
			perfMapLock.Unlock()
		}
		time.Sleep(time.Duration(OffSleepTime) * time.Millisecond)
	}
}

func Safe(name string, callback func()) {
	start := time.Now()
	defer func() {
		perfMapLock.Lock()
		PerfMap[name] = int32(time.Since(start).Nanoseconds() / 1000)
		perfMapLock.Unlock()
	}()

	// fmt.Println("Safe step -1")
	if AllOff {
		// fmt.Println("Safe step -2")
		time.Sleep(time.Duration(OffSleepTime) * time.Millisecond)
		return
	}
	// fmt.Println("Safe step -3")
	if on, ok := onofflookup[name]; ok {
		// fmt.Println("Safe step -4")
		if !on {
			// fmt.Println("Safe step -5")
			time.Sleep(time.Duration(OffSleepTime) * time.Millisecond)
		} else {
			// fmt.Println("Safe step -6")
			callback()
		}
	} else {
		// fmt.Println("Safe step -7")
		callback()
	}
	// fmt.Println("Safe step -8")
}

func SetLoopOffMap(m *map[string]bool) {
	onofflookup = *m
	// val, ok := onofflookup["countermanager.poll"]
	// fmt.Println("SetLoopOffMap val:", val, "ok:", ok, "size:", len(onofflookup))
}
func SetOnOff(name string, onoff bool) {
	onofflookup[name] = onoff
}
