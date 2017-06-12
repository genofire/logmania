// log package with entry as a lib in other go applications
package log

import "fmt"

// a struct with all information of a log entry
type Entry struct {
	Level  LogLevel               `json:"level"`
	Fields map[string]interface{} `json:"fields"`
	Text   string                 `json:"text"`
}

// save/out current state of log entry
func (e *Entry) Log(level LogLevel, v ...interface{}) {
	e.Text = fmt.Sprint(v...)
	e.Level = level
	save(e)
}

// save/out current state of log entry with formation
func (e *Entry) Logf(level LogLevel, format string, v ...interface{}) {
	e.Text = fmt.Sprintf(format, v...)
	e.Level = level
	save(e)
}

// init new log entry
func New() *Entry {
	return &Entry{Fields: make(map[string]interface{})}
}

// add extra value to entry (log entry with context)
func (e *Entry) AddField(key string, value interface{}) *Entry {
	e.Fields[key] = value
	return e
}

// add multi extra values to entry (log entry with context)
func (e *Entry) AddFields(fields map[string]interface{}) *Entry {
	for key, value := range fields {
		e.Fields[key] = value
	}
	return e
}

// create a readable string of extra values (log entry with context)
func (e *Entry) FieldString() string {
	text := ""
	for key, value := range e.Fields {
		text = fmt.Sprintf("%s %s=%v", text, key, value)
	}
	return text[1:]
}
