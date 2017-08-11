package syslog

import (
	libSyslog "github.com/genofire/logmania/lib/syslog"
	"github.com/genofire/logmania/log"
)

var SyslogPriorityMap = map[int]log.LogLevel{
	0: log.PanicLevel,
	1: log.PanicLevel,
	2: log.PanicLevel,
	3: log.ErrorLevel,
	4: log.WarnLevel,
	5: log.InfoLevel,
	6: log.InfoLevel,
	7: log.DebugLevel,
}

func toLogEntry(msg []byte, from string) *log.Entry {
	syslogMsg := libSyslog.Parse(msg)

	return &log.Entry{
		Level:    SyslogPriorityMap[syslogMsg.Severity],
		Text:     syslogMsg.Content,
		Hostname: from,
	}
}
