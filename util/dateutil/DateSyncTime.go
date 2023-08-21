package dateutil

import (
	//	"fmt"
	"sync"
	"time"
)

// #include <unistd.h>
// //#include <errno.h>
// //int usleep(useconds_t usec);
import "C"

const (
	TIME_SYNC_INTERVAL = 5000
)

var SyncTimeMillis int64
var SyncTime time.Time
var lastSyncTime int64
var syncTimeTicker *time.Ticker
var syncTimeLock sync.Mutex

func StartSyncTime() {
	syncTimeLock.Lock()
	defer syncTimeLock.Unlock()
	if syncTimeTicker != nil {
		return
	}
	syncTimeTicker = time.NewTicker(time.Millisecond)
	clock()
	SyncTimeMillis = systemMillis()
	lastSyncTime = SyncTimeMillis

}

func StopSyncTime() {
	syncTimeLock.Lock()
	defer syncTimeLock.Unlock()
	func() {
		defer func() {
			if r := recover(); r != nil {
			}
		}()
		if syncTimeTicker != nil {
			// can panic
			syncTimeTicker.Stop()
		}
	}()
}

func IsSyncTime() bool {
	if syncTimeTicker != nil {
		return true
	}
	return false
}

func clock() {
	go func() {
		cnt := 0
		// per millisecond
		for t := range syncTimeTicker.C {
			cnt++
			SyncTime = t
			if cnt > TIME_SYNC_INTERVAL {
				cnt = 0
				SyncTime = time.Now()
				now := systemMillis()
				//fmt.Println(">>>> sync gap =", (now - SyncTimeMillis))
				lastSyncTime = now
				SyncTimeMillis = now
			}
		}
	}()
}

func clockUsleep() {
	go func() {
		cnt := 0
		for {
			cnt++
			SyncTime = time.Now()
			if cnt > TIME_SYNC_INTERVAL {
				cnt = 0
				SyncTime = time.Now()
				now := systemMillis()
				//fmt.Println(">>>> sync gap =", (now - SyncTimeMillis))
				lastSyncTime = now
				SyncTimeMillis = now
			}

			C.usleep(1000)

		}
	}()
}

func systemMillis() int64 {
	return time.Now().UnixNano() / 1000000
}
