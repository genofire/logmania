package log

import "fmt"

type Entry struct {
	Level  LogLevel               `json:"level"`
	Fields map[string]interface{} `json:"fields"`
	Text   string                 `json:"text"`
}

func (e *Entry) Log(level LogLevel, v ...interface{}) {
	e.Text = fmt.Sprint(v...)
	e.Level = level
	save(e)
}
func (e *Entry) Logf(level LogLevel, format string, v ...interface{}) {
	e.Text = fmt.Sprintf(format, v...)
	e.Level = level
	save(e)
}

func New() *Entry {
	return &Entry{Fields: make(map[string]interface{})}
}

func (e *Entry) AddField(key string, value interface{}) *Entry {
	e.Fields[key] = value
	return e
}
func (e *Entry) AddFields(fields map[string]interface{}) *Entry {
	for key, value := range fields {
		e.Fields[key] = value
	}
	return e
}

func (e *Entry) FieldString() string {
	text := ""
	for key, value := range e.Fields {
		text = fmt.Sprintf("%s %s=%v", text, key, value)
	}
	return text[1:]
}
