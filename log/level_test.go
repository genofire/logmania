package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	assert := assert.New(t)

	results := map[LogLevel]string{
		DebugLevel:   "Debug",
		InfoLevel:    "Info",
		WarnLevel:    "Warn",
		ErrorLevel:   "ERROR",
		PanicLevel:   "PANIC",
		LogLevel(-2): "NOT VALID",
	}

	for value, expected := range results {
		assert.Equal(expected, value.String())
	}
}

func TestLogLevelFunc(t *testing.T) {
	assert := assert.New(t)

	results := map[LogLevel]func(...interface{}){
		DebugLevel: entry.Debug,
		InfoLevel:  entry.Info,
		WarnLevel:  entry.Warn,
		ErrorLevel: entry.Error,
	}

	for value, function := range results {
		function()
		assert.Equal(value, entry.Level)
	}
	assert.Panics(func() {
		entry.Panic()
		assert.Equal(PanicLevel, entry.Level)
	})
}
func TestLogLevelFormatFunc(t *testing.T) {
	assert := assert.New(t)

	results := map[LogLevel]func(string, ...interface{}){
		DebugLevel: entry.Debugf,
		InfoLevel:  entry.Infof,
		WarnLevel:  entry.Warnf,
		ErrorLevel: entry.Errorf,
	}

	for value, function := range results {
		function("%.1f", 31.121)
		assert.Equal(value, entry.Level)
		assert.Equal("31.1", entry.Text)
	}
	assert.Panics(func() {
		entry.Panicf("")
		assert.Equal(PanicLevel, entry.Level)
	})
}

func TestLogLevelInit(t *testing.T) {
	assert := assert.New(t)

	results := map[LogLevel]func(...interface{}){
		DebugLevel: Debug,
		InfoLevel:  Info,
		WarnLevel:  Warn,
		ErrorLevel: Error,
	}

	for value, function := range results {
		function()
		assert.Equal(value, entry.Level)
	}
	assert.Panics(func() {
		Panic()
		assert.Equal(PanicLevel, entry.Level)
	})
}

func TestLogLevelInitFormatFunc(t *testing.T) {
	assert := assert.New(t)

	results := map[LogLevel]func(string, ...interface{}){
		DebugLevel: Debugf,
		InfoLevel:  Infof,
		WarnLevel:  Warnf,
		ErrorLevel: Errorf,
	}

	for value, function := range results {
		function("%.1f", 31.121)
		assert.Equal(value, entry.Level)
		assert.Equal("31.1", entry.Text)
	}
	assert.Panics(func() {
		Panicf("")
		assert.Equal(PanicLevel, entry.Level)
	})
}
