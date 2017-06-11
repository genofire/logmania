package output

import (
	"fmt"
	"os"
	"time"

	"github.com/genofire/logmania/log"
)

var (
	TimeFormat = "2006-01-02 15:04:05"
	ShowTime   = true
	AboveLevel = log.InfoLevel
)

func hook(e *log.Entry) {
	if e.Level < AboveLevel {
		return
	}
	v := []interface{}{}
	format := "[%s] %s"

	if ShowTime {
		format = "%s [%s] %s"
		v = append(v, time.Now().Format(TimeFormat))
	}

	v = append(v, e.Level.String(), e.Text)

	if len(e.Fields) > 0 {
		v = append(v, e.FieldString())
		format = fmt.Sprintf("%s (%%s)\n", format)
	} else {
		format = fmt.Sprintf("%s\n", format)
	}

	text := fmt.Sprintf(format, v...)

	if e.Level == log.PanicLevel {
		panic(text)
	} else if e.Level > log.WarnLevel {
		os.Stderr.WriteString(text)
	} else {
		os.Stdout.WriteString(text)
	}
}

func init() {
	log.AddHook(hook)
}
