package logger

import (
	"fmt"
	"log"
)

type DefaultLogger struct{}

func NewDefaultLogger() *DefaultLogger {
	p := new(DefaultLogger)
	log.SetFlags(log.Ldate | log.Ltime)
	return p
}
func (this *DefaultLogger) Errorf(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	log.Println("[Error] ", s)
}
func (this *DefaultLogger) Error(args ...interface{}) {
	s := fmt.Sprintln(args...)
	log.Println("[Error] ", s)
}
func (this *DefaultLogger) Warnf(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	log.Println("[Warn] ", s)
}
func (this *DefaultLogger) Warn(args ...interface{}) {
	s := fmt.Sprintln(args...)
	log.Println("[Warn] ", s)
}
func (this *DefaultLogger) Infof(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	log.Println("[Info] ", s)
}
func (this *DefaultLogger) Info(args ...interface{}) {
	s := fmt.Sprintln(args...)
	log.Println("[Info] ", s)
}
func (this *DefaultLogger) Infoln(args ...interface{}) {
	s := fmt.Sprintln(args...)
	log.Println("[Info] ", s)
}
func (this *DefaultLogger) Printf(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	log.Println("[Print] ", s)
}
func (this *DefaultLogger) Println(args ...interface{}) {
	s := fmt.Sprintln(args...)
	log.Println("[Print] ", s)
}
func (this *DefaultLogger) Debugf(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	log.Println("[Debug] ", s)
}
func (this *DefaultLogger) Debug(args ...interface{}) {
	s := fmt.Sprintln(args...)
	log.Println("[Debug] ", s)
}
