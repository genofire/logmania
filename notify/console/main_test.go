package console

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
	assert.Equal("", err)

	out, err = captureOutput(func() {
		log.Warn("test")
	})
	assert.NotRegexp("-.*\\[.{5}Warn.{4}\\] test", out)
	assert.Equal("", err)

	out, err = captureOutput(func() {
		log.Error("test")
	})
	assert.Equal("", out)

	out, err = captureOutput(func() {
		log.Debug("test")
	})
	assert.Equal("", out)
	assert.Equal("", err)

	out, err = captureOutput(func() {
		log.Info("test")
	})
	assert.Equal("", out)
	assert.Equal("", err)

}
