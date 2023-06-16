//+build !windows

package panicutil

import (
	"github.com/natefinch/lumberjack"
)

func getErrorLogger() *lumberjack.Logger {
	if errorLogger == nil {
		errorLogger = &lumberjack.Logger{
			Filename:   "/var/log/whatap_infrad.log",
			MaxSize:    10, // megabytes after which new file is created
			MaxBackups: 2,  // number of backups
			MaxAge:     37, //days
		}
	}

	return errorLogger
}
