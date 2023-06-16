package panicutil

import (
	"context"
	"fmt"
	"runtime/debug"
	"runtime/pprof"
	"time"

	"github.com/natefinch/lumberjack"
)

//PanicLogger panic logger
func PanicLogger() {
	if err := recover(); err != nil {
		errorLogger := getErrorLogger()
		t := time.Now()
		errorLogger.Write([]byte(t.Format("2006/01/02 15:04:05 ")))
		errorLogger.Write([]byte(fmt.Sprint(err)))
		errorLogger.Write([]byte(string(debug.Stack())))
	}
}

var errorLogger *lumberjack.Logger

func LogAllStack() {
	if IsDebug {
		errorLogger := getErrorLogger()
		t := time.Now()
		errorLogger.Write([]byte(t.Format("2006/01/02 15:04:05 ")))
		pprof.Lookup("goroutine").WriteTo(errorLogger, 1)
	}
}

func DoWithTimeout(timeout time.Duration, doFunc func(), msg string, doPanic bool) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	go func() {
		defer PanicLogger()
		doFunc()
		cancel()
	}()

	select {
	case <-time.After(timeout):
		break
	case <-ctx.Done():
		return
	}

	panicmsg := fmt.Sprint("\n", time.Now().Format("2006/01/02 15:04:05 "), "Operation Timeout", msg)
	getErrorLogger().Write([]byte(panicmsg))
	if doPanic {
		panic(panicmsg)
	}
}

func DoWithTimeoutEx(timeout time.Duration, doFunc func(), msg string, onTimeout func()) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	go func() {
		defer PanicLogger()
		doFunc()
		cancel()
	}()

	select {
	case <-time.After(timeout):
		break
	case <-ctx.Done():
		return
	}
	if onTimeout != nil {
		onTimeout()
	}
}

var IsDebug = false

//Debug Debug
func Debug(args ...interface{}) {
	if IsDebug {
		errorLogger := getErrorLogger()
		t := time.Now()
		message := fmt.Sprintln(t.Format("2006/01/02 15:04:05 [DEBUG] "), args)
		errorLogger.Write([]byte(message))
	}

}

//Info Info
func Info(args ...interface{}) {
	errorLogger := getErrorLogger()
	t := time.Now()
	message := fmt.Sprintln(t.Format("2006/01/02 15:04:05 [INFO] "), args)
	errorLogger.Write([]byte(message))
}

//Error Error
func Error(args ...interface{}) {
	errorLogger := getErrorLogger()
	t := time.Now()
	message := fmt.Sprintln(t.Format("2006/01/02 15:04:05 [Error] "), args)
	errorLogger.Write([]byte(message))

}

//SendEvent SendEvent
func SendEvent(title string, message string) {

}
