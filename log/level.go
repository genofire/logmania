package log

type LogLevel int32

const (
	DebugLevel = LogLevel(-1)
	InfoLevel  = LogLevel(0)
	WarnLevel  = LogLevel(1)
	ErrorLevel = LogLevel(2)
	PanicLevel = LogLevel(3)
)

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

/**
 * log command
 */

// debug
func (e *Entry) Debug(v ...interface{}) {
	e.Log(DebugLevel, v...)
}
func (e *Entry) Debugf(format string, v ...interface{}) {
	e.Logf(DebugLevel, format, v...)
}

// info
func (e *Entry) Info(v ...interface{}) {
	e.Log(InfoLevel, v...)
}
func (e *Entry) Infof(format string, v ...interface{}) {
	e.Logf(InfoLevel, format, v...)
}

// warn
func (e *Entry) Warn(v ...interface{}) {
	e.Log(WarnLevel, v...)
}
func (e *Entry) Warnf(format string, v ...interface{}) {
	e.Logf(WarnLevel, format, v...)
}

// error
func (e *Entry) Error(v ...interface{}) {
	e.Log(ErrorLevel, v...)
}
func (e *Entry) Errorf(format string, v ...interface{}) {
	e.Logf(ErrorLevel, format, v...)
}

// panic
func (e *Entry) Panic(v ...interface{}) {
	e.Log(PanicLevel, v...)
}
func (e *Entry) Panicf(format string, v ...interface{}) {
	e.Logf(PanicLevel, format, v...)
}

/**
 * Direct log command
 */

// debug
func Debug(v ...interface{}) {
	New().Log(DebugLevel, v...)
}
func Debugf(format string, v ...interface{}) {
	New().Logf(DebugLevel, format, v...)
}

// info
func Info(v ...interface{}) {
	New().Log(InfoLevel, v...)
}
func Infof(format string, v ...interface{}) {
	New().Logf(InfoLevel, format, v...)
}

// warn
func Warn(v ...interface{}) {
	New().Log(WarnLevel, v...)
}
func Warnf(format string, v ...interface{}) {
	New().Logf(WarnLevel, format, v...)
}

// error
func Error(v ...interface{}) {
	New().Log(ErrorLevel, v...)
}
func Errorf(format string, v ...interface{}) {
	New().Logf(ErrorLevel, format, v...)
}

// panic
func Panic(v ...interface{}) {
	New().Log(PanicLevel, v...)
}
func Panicf(format string, v ...interface{}) {
	New().Logf(PanicLevel, format, v...)
}
