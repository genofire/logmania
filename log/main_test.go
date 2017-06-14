package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLog(t *testing.T) {
	assert := assert.New(t)
	//save := func(e *Entry) {}
	entry := New()
	assert.Equal(0, int(entry.Level))
	assert.Equal("", entry.Text)

	entry.Log(WarnLevel, "blub")
	assert.Equal(1, int(entry.Level))
	assert.Equal("blub", entry.Text)

	entry.Logf(ErrorLevel, "lola %.1f", 13.13431)
	assert.Equal(2, int(entry.Level))
	assert.Equal("lola 13.1", entry.Text)
}

func TestAddFields(t *testing.T) {
	assert := assert.New(t)
	entry := New()
	assert.Len(entry.Fields, 0)

	entry.AddField("a", "lola")
	assert.Len(entry.Fields, 1)
	assert.Equal("lola", entry.Fields["a"])

	entry.AddFields(map[string]interface{}{"a": 232., "foo": "bar"})
	assert.Len(entry.Fields, 2)
	assert.Equal(232.0, entry.Fields["a"])
}

func TestFieldString(t *testing.T) {
	assert := assert.New(t)
	entry := New()
	entry.AddFields(map[string]interface{}{"a": 232., "foo": "bar"})
	str := entry.FieldString()
	assert.Contains(str, "a=232")
	assert.Contains(str, "foo=bar")
}
