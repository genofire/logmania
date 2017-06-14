package output

import (
	"bytes"
	"os"
	"testing"

	"github.com/genofire/logmania/log"
	"github.com/stretchr/testify/assert"
)

func captureOutput(f func()) (string, string) {
	var bufErrOutput bytes.Buffer
	var bufOutput bytes.Buffer
	errOutput = &bufErrOutput
	output = &bufOutput
	f()
	errOutput = os.Stderr
	output = os.Stdout
	return bufOutput.String(), bufErrOutput.String()
}

// Warning: colors are not tested (it should be in the imported package)
func TestOutput(t *testing.T) {
	assert := assert.New(t)
	assert.True(true)
	out, err := captureOutput(func() {
		log.Info("test")
	})
	assert.Regexp("-.*\\[.{5}Info.{4}\\] test", out)
	assert.Equal("", err)

	ShowTime = false
	out, err = captureOutput(func() {
		log.Warn("test")
	})
	assert.Regexp("\\[.{5}Warn.{4}\\] test", out)
	assert.NotRegexp("-.*\\[.{5}Warn.{4}\\] test", out)
	assert.Equal("", err)

	out, err = captureOutput(func() {
		log.Error("test")
	})
	assert.Equal("", out)
	assert.Regexp("\\[.{5}ERROR.{4}\\] test", err)

	out, err = captureOutput(func() {
		log.Debug("test")
	})
	assert.Equal("", out)
	assert.Equal("", err)

	AboveLevel = log.DebugLevel

	out, err = captureOutput(func() {
		log.New().AddField("a", 3).Debug("test")
	})
	assert.Regexp("\\[.{5}Debug.{4}\\] test .{8}(a=3)", out)
	assert.Equal("", err)

	log.RemoveLogger("output")

	out, err = captureOutput(func() {
		log.Info("test")
	})
	assert.Equal("", out)
	assert.Equal("", err)

}
