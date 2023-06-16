//+build windows

package panicutil

import (
	"fmt"
	"syscall"
	"time"
)

var prevTime, prevUsage uint64

//GetSelfCPUUsage GetSelfCPUUsage
func GetSelfCPUUsage() (float32, error) {
	handle, err := syscall.GetCurrentProcess()
	if err != nil {
		return 0, err
	}
	// fmt.Println("GetSelfCPUUsage ", prevTime, prevUsage)
	var ctime, etime, ktime, utime syscall.Filetime
	err = syscall.GetProcessTimes(handle, &ctime, &etime, &ktime, &utime)
	if err != nil {
		return 0, err
	}

	curTime := uint64(time.Now().UnixNano())
	curUsage := uint64(ktime.Nanoseconds()) + uint64(utime.Nanoseconds()) // Always overflows

	var cpuUsagePercent float32
	if prevTime > 0 && prevUsage > 0 {
		timeDiff := curTime - prevTime
		usageDiff := curUsage - prevUsage
		cpuUsagePercent = float32(100 * float64(usageDiff) / float64(timeDiff))
		err = nil
	} else {
		err = fmt.Errorf("requires to run GetSelfCPUUsage twice at least")
	}

	prevTime = curTime
	prevUsage = curUsage

	return cpuUsagePercent, err
}
