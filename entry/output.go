package entry

import (
	"fmt"
	"os"
	"time"
)

var TimeFormat = "2006-01-02 15:04:05"

var InternelSend = func(e *Entry) {
	format := "%s [%s] %s\n"
	v := []interface{}{time.Now().Format(TimeFormat), e.Level.String(), e.Text}
	if len(e.Fields) > 0 {
		format = "%s [%s] %s (%s)\n"
		v = append(v, e.FieldString())
	}
	text := fmt.Sprintf(format, v...)

	if e.Level == PanicLevel {
		panic(text)
	} else if e.Level > WarnLevel {
		os.Stderr.WriteString(text)
	} else {
		os.Stdout.WriteString(text)
	}
}

var FieldOutput = func(fields map[string]interface{}) string {
	text := ""
	for key, value := range fields {
		text = fmt.Sprintf("%s %s=%v", text, key, value)
	}
	return text[1:]
}

func save(e *Entry) {
	InternelSend(e)
}
