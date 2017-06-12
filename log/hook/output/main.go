// logger to print log entry (with color)
// this logger would be bind by importing
package output

import (
	"fmt"
	"os"
	"time"

	"github.com/bclicn/color"

	"github.com/genofire/logmania/log"
)

var (
	TimeFormat = "2006-01-02 15:04:05"
	ShowTime   = true
	AboveLevel = log.InfoLevel
)

// logger for output
type Logger struct {
	log.Logger
	TimeFormat string
	ShowTime   bool
	AboveLevel log.LogLevel
}

// CurrentLogger (for override settings e.g. AboveLevel,ShowTime or TimeFormat)
var CurrentLogger *Logger

// create a new output logger
func NewLogger() *Logger {
	return &Logger{
		TimeFormat: "2006-01-02 15:04:05",
		ShowTime:   true,
		AboveLevel: log.InfoLevel,
	}
}

// handle a log entry (print it on the terminal with color)
func (l *Logger) Hook(e *log.Entry) {
	if e.Level < AboveLevel {
		return
	}
	v := []interface{}{}
	format := "[%s] %s"

	if ShowTime {
		format = "%s [%s] %s"
		v = append(v, color.LightBlue(time.Now().Format(TimeFormat)))
	}
	lvl := e.Level.String()
	switch e.Level {
	case log.DebugLevel:
		lvl = color.DarkGray(lvl)
	case log.InfoLevel:
		lvl = color.Green(lvl)
	case log.WarnLevel:
		lvl = color.Yellow(lvl)
	case log.ErrorLevel:
		lvl = color.Red(lvl)
	case log.PanicLevel:
		lvl = color.BRed(lvl)
	}

	v = append(v, lvl, e.Text)

	if len(e.Fields) > 0 {
		v = append(v, color.Purple(e.FieldString()))
		format = fmt.Sprintf("%s (%%s)\n", format)
	} else {
		format = fmt.Sprintf("%s\n", format)
	}

	text := fmt.Sprintf(format, v...)

	if e.Level > log.WarnLevel {
		os.Stderr.WriteString(text)
	} else {
		os.Stdout.WriteString(text)
	}
}

// do nothing - terminal did not need something to close
func (l *Logger) Close() {
}

func init() {
	CurrentLogger = NewLogger()
	log.AddLogger(CurrentLogger)
}
