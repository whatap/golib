package logger

import (
	"fmt"
	"log"
)

type DefaultLogger struct {
	Level int
}

func NewDefaultLogger() *DefaultLogger {
	p := new(DefaultLogger)
	p.Level = LOG_LEVEL_WARN
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
	if this.Level <= LOG_LEVEL_WARN {
		return
	}
	s := fmt.Sprintf(format, args...)
	log.Println("[Warn] ", s)
}
func (this *DefaultLogger) Warn(args ...interface{}) {
	if this.Level <= LOG_LEVEL_WARN {
		return
	}
	s := fmt.Sprintln(args...)
	log.Println("[Warn] ", s)
}
func (this *DefaultLogger) Infof(format string, args ...interface{}) {
	if this.Level <= LOG_LEVEL_INFO {
		return
	}
	s := fmt.Sprintf(format, args...)
	log.Println("[Info] ", s)
}
func (this *DefaultLogger) Info(args ...interface{}) {
	if this.Level <= LOG_LEVEL_INFO {
		return
	}
	s := fmt.Sprintln(args...)
	log.Println("[Info] ", s)
}
func (this *DefaultLogger) Infoln(args ...interface{}) {
	if this.Level <= LOG_LEVEL_INFO {
		return
	}
	s := fmt.Sprintln(args...)
	log.Println("[Info] ", s)
}
func (this *DefaultLogger) Printf(id string, format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	log.Println(id, s)
}
func (this *DefaultLogger) Println(id string, args ...interface{}) {
	s := fmt.Sprintln(args...)
	log.Println(id, s)
}
func (this *DefaultLogger) Debugf(format string, args ...interface{}) {
	if this.Level <= LOG_LEVEL_DEBUG {
		return
	}
	s := fmt.Sprintf(format, args...)
	log.Println("[Debug] ", s)
}
func (this *DefaultLogger) Debug(args ...interface{}) {
	if this.Level <= LOG_LEVEL_DEBUG {
		return
	}
	s := fmt.Sprintln(args...)
	log.Println("[Debug] ", s)
}
