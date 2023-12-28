package dateutil

import (
	//	"fmt"
	"time"
)

const (
	TIME_SYNC_INTERVAL = 5000
)

var SyncTimeMillis int64
var lastSyncTime int64
var syncTimeTicker *time.Ticker

func StartSyncTime() {
	SyncTimeMillis = systemMillis()
	lastSyncTime = SyncTimeMillis
	clock()
}

func StopSyncTime() {
	syncTimeTicker.Stop()
}

func IsSyncTime() bool {
	if syncTimeTicker != nil {
		return true
	}
	return false
}

func clock() {
	syncTimeTicker = time.NewTicker(time.Millisecond)
	go func() {
		for t := range syncTimeTicker.C {
			SyncTimeMillis = t.UnixNano() / 1000000
			if SyncTimeMillis > lastSyncTime+TIME_SYNC_INTERVAL {
				now := systemMillis()
				//fmt.Println(">>>> sync gap =", (now - SyncTimeMillis))
				lastSyncTime = now
				SyncTimeMillis = now
			}
		}
	}()
}
func systemMillis() int64 {
	return time.Now().UnixNano() / 1000000
}
