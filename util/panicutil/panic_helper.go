package panicutil

import (
	"context"
	"fmt"
	"runtime"
	"runtime/debug"
	"runtime/pprof"
	"strings"
	"time"

	"github.com/natefinch/lumberjack"
	"github.com/whatap/golib/util/ansi"
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

// getCallerName returns the caller function name (skip levels up the call stack)
func getCallerName(skip int) string {
	pc, _, _, ok := runtime.Caller(skip)
	if !ok {
		return ""
	}
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return ""
	}
	// Extract short name: "github.com/whatap/opsgent/filelog.(*Collector).collect" -> "filelog.(*Collector).collect"
	name := fn.Name()
	if idx := strings.LastIndex(name, "/"); idx >= 0 {
		name = name[idx+1:]
	}
	return name
}

//Debug Debug
func Debug(args ...interface{}) {
	if IsDebug {
		errorLogger := getErrorLogger()
		t := time.Now()
		caller := getCallerName(2)
		message := fmt.Sprintf("%s [DEBUG] %s %s\n", t.Format("2006/01/02 15:04:05"), ansi.Yellow("["+caller+"]"), fmt.Sprint(args...))
		errorLogger.Write([]byte(message))
	}

}

//Info Info
func Info(args ...interface{}) {
	errorLogger := getErrorLogger()
	t := time.Now()
	caller := getCallerName(2)
	message := fmt.Sprintf("%s [INFO] %s %s\n", t.Format("2006/01/02 15:04:05"), ansi.Yellow("["+caller+"]"), fmt.Sprint(args...))
	errorLogger.Write([]byte(message))
}

//Error Error
func Error(args ...interface{}) {
	errorLogger := getErrorLogger()
	t := time.Now()
	caller := getCallerName(2)
	message := fmt.Sprintf("%s [Error] %s %s\n", t.Format("2006/01/02 15:04:05"), ansi.Yellow("["+caller+"]"), fmt.Sprint(args...))
	errorLogger.Write([]byte(message))

}

//SendEvent SendEvent
func SendEvent(title string, message string) {

}
