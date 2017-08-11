package log

// definition of loglevel
type LogLevel int32

// accepted LogLevels and his internal int values
const (
	DebugLevel = LogLevel(-1)
	InfoLevel  = LogLevel(0)
	WarnLevel  = LogLevel(1)
	ErrorLevel = LogLevel(2)
	PanicLevel = LogLevel(3)
)

// string of loglevel
func (l LogLevel) String() string {
	switch l {
	case DebugLevel:
		return "Debug"
	case InfoLevel:
		return "Info"
	case WarnLevel:
		return "Warn"
	case ErrorLevel:
		return "ERROR"
	case PanicLevel:
		return "PANIC"

	}
	return "NOT VALID"
}

func NewLoglevel(by string) LogLevel {
	switch by {
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "warn":
		return WarnLevel
	case "error":
		return ErrorLevel
	case "panic":
		return PanicLevel

	}
	return DebugLevel
}

/**
 * log command
 */

// close logentry with debug
func (e *Entry) Debug(v ...interface{}) {
	e.Log(DebugLevel, v...)
}

// close logentry with formated debug
func (e *Entry) Debugf(format string, v ...interface{}) {
	e.Logf(DebugLevel, format, v...)
}

// close logentry with info
func (e *Entry) Info(v ...interface{}) {
	e.Log(InfoLevel, v...)
}

// close logentry with formated info
func (e *Entry) Infof(format string, v ...interface{}) {
	e.Logf(InfoLevel, format, v...)
}

// close logentry with warning
func (e *Entry) Warn(v ...interface{}) {
	e.Log(WarnLevel, v...)
}

// close logentry with formated warning
func (e *Entry) Warnf(format string, v ...interface{}) {
	e.Logf(WarnLevel, format, v...)
}

// close logentry with error
func (e *Entry) Error(v ...interface{}) {
	e.Log(ErrorLevel, v...)
}

// close logentry with formated error
func (e *Entry) Errorf(format string, v ...interface{}) {
	e.Logf(ErrorLevel, format, v...)
}

// close logentry with panic
func (e *Entry) Panic(v ...interface{}) {
	e.Log(PanicLevel, v...)
}

// close logentry with formated panic
func (e *Entry) Panicf(format string, v ...interface{}) {
	e.Logf(PanicLevel, format, v...)
}

/**
 * Direct log command
 */

//  direct log with debug
func Debug(v ...interface{}) {
	New().Log(DebugLevel, v...)
}

//  direct log with formated debug
func Debugf(format string, v ...interface{}) {
	New().Logf(DebugLevel, format, v...)
}

// direct log with info
func Info(v ...interface{}) {
	New().Log(InfoLevel, v...)
}

// direct log with formated info
func Infof(format string, v ...interface{}) {
	New().Logf(InfoLevel, format, v...)
}

// direct log with warning
func Warn(v ...interface{}) {
	New().Log(WarnLevel, v...)
}

// direct log with formated warning
func Warnf(format string, v ...interface{}) {
	New().Logf(WarnLevel, format, v...)
}

// direct log with error
func Error(v ...interface{}) {
	New().Log(ErrorLevel, v...)
}

// direct log with formated error
func Errorf(format string, v ...interface{}) {
	New().Logf(ErrorLevel, format, v...)
}

// direct log with panic
func Panic(v ...interface{}) {
	New().Log(PanicLevel, v...)
}

// direct log with formated panic
func Panicf(format string, v ...interface{}) {
	New().Logf(PanicLevel, format, v...)
}
