package syslog

import (
	"testing"

	"github.com/stretchr/testify/assert"

	log "github.com/sirupsen/logrus"
)

func TestToEntry(t *testing.T) {
	assert := assert.New(t)
	entry := toLogEntry([]byte("<11>Aug 17 11:43:33 Msg"), "::1")
	assert.Equal("Msg", entry.Message)
	assert.Equal(log.ErrorLevel, entry.Level)

	hostname, ok := entry.Data["hostname"]
	assert.True(ok)
	assert.Equal("::1", hostname)

}
