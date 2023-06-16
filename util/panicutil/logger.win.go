//+build windows

package panicutil

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/natefinch/lumberjack"
)

func getErrorLogger() *lumberjack.Logger {
	if errorLogger == nil {
		ex, err := os.Executable()
		if err != nil {
			panic(err)
		}
		exPath := filepath.Dir(ex)

		errorLogger = &lumberjack.Logger{
			Filename:   fmt.Sprintf("%s/error.log", exPath),
			MaxSize:    10, // megabytes after which new file is created
			MaxBackups: 2,  // number of backups
			MaxAge:     37, //days
		}
	}

	return errorLogger
}
