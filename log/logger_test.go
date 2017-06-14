package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var entry *Entry

type SaveLogger struct {
	Logger
}

func (*SaveLogger) Hook(e *Entry) {
	entry = e
}
func (*SaveLogger) Close() {}

func init() {
	entry = &Entry{}
	AddLogger("name", &SaveLogger{})
}

func TestLogger(t *testing.T) {
	assert := assert.New(t)
	assert.Len(loggers, 1)

	AddLogger("blub", &SaveLogger{})
	assert.Len(loggers, 2)
	RemoveLogger("blub")
	assert.Len(loggers, 1)

	assert.PanicsWithValue("panic see last log in logmania", func() {
		save(&Entry{Level: PanicLevel})
	})
}
