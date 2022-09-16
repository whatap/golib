package logger

const (
	LOG_LEVEL_ERROR = 0
	LOG_LEVEL_WARN  = 1
	LOG_LEVEL_INFO  = 2
	LOG_LEVEL_DEBUG = 3
)

type Logger interface {
	// Errorf logs an error message, patterned after log.Printf.
	Errorf(format string, args ...interface{})
	// Error logs an error message, patterned after log.Print.
	Error(args ...interface{})
	// Warnf logs a warning message, patterned after log.Printf.
	Warnf(format string, args ...interface{})
	// Warn logs a warning message, patterned after log.Print.
	Warn(args ...interface{})
	// Infof logs an information message, patterned after log.Printf.
	Infof(format string, args ...interface{})
	// Info logs an information message, patterned after log.Print.
	Info(args ...interface{})
	Infoln(args ...interface{})
	// Debugf logs a debug message, patterned after log.Printf.
	Debugf(format string, args ...interface{})
	// Debug logs a debug message, patterned after log.Print.
	Debug(args ...interface{})

	// whatap cache log
	Printf(id string, format string, args ...interface{})
	Println(id string, args ...interface{})
}

type EmptyLogger struct{}

func (el *EmptyLogger) Errorf(format string, args ...interface{})            {}
func (el *EmptyLogger) Error(args ...interface{})                            {}
func (el *EmptyLogger) Warnf(format string, args ...interface{})             {}
func (el *EmptyLogger) Warn(args ...interface{})                             {}
func (el *EmptyLogger) Infof(format string, args ...interface{})             {}
func (el *EmptyLogger) Info(args ...interface{})                             {}
func (el *EmptyLogger) Infoln(args ...interface{})                           {}
func (el *EmptyLogger) Debugf(format string, args ...interface{})            {}
func (el *EmptyLogger) Debug(args ...interface{})                            {}
func (el *EmptyLogger) Printf(id string, format string, args ...interface{}) {}
func (el *EmptyLogger) Println(id string, args ...interface{})               {}
