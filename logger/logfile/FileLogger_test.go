package logfile

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/whatap/golib/logger"
)

func TestFileGetLogger(t *testing.T) {
	Log := GetFileLogger()
	f := Log.GetLogFile()
	assert.NotNil(t, f)
	Log.Println("aaa")
	Log.Error("Error ")
}

func TestFileLoggerOptionOnameLogID(t *testing.T) {
	Log := NewFileLogger(WithOnameLogID("rumctl", "RUM"))
	f := Log.GetLogFile()
	assert.NotNil(t, f)
	Log.Println("tytttt")
	Log.Error("Error ")
}

func TestFileLoggerGetLogFiles(t *testing.T) {
	Log := NewFileLogger(WithOnameLogID("rumctl", "RUM"))
	f := Log.GetLogFile()
	if !assert.NotNil(t, f) {
		return
	}

	Log.Println("tytttt")
	Log.Error("Error ")

	mapValue := Log.GetLogFiles()
	if !assert.NotNil(t, mapValue) {
		return
	}

	en := mapValue.Keys()

	for en.HasMoreElements() {
		next := en.NextString()
		assert.True(t, strings.HasPrefix(next, "RUM-rumctl"))
	}
	//fmt.Println(mapValue.ToString())
}

func TestFileLoggerReadLog(t *testing.T) {
	Log := NewFileLogger(WithOnameLogID("rumctl", "RUM"), WithLevel(logger.LOG_LEVEL_DEBUG))

	f := Log.GetLogFile()
	if !assert.NotNil(t, f) {
		return
	}

	Log.Println("tytttt")
	Log.Error("Error ")
	name := filepath.Base(f.Name())
	data := Log.Read(name, 0, 1024)
	if !assert.NotNil(t, data) {
		return
	}
}
