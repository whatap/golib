//+build !windows

package panicutil

import (
	"os"

	"github.com/shirou/gopsutil/process"
)

//GetSelfCPUUsage GetSelfCPUUsage
func GetSelfCPUUsage() (float32, error) {
	if thisProcess == nil {
		p, err := process.NewProcess(int32(os.Getpid()))
		if err != nil {
			return 0, err
		}
		thisProcess = p
	}

	cpuUsagePercent, err := thisProcess.CPUPercent()

	return float32(cpuUsagePercent), err
}
