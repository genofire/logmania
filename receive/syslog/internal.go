package syslog

import (
	log "github.com/sirupsen/logrus"

	libSyslog "github.com/genofire/logmania/lib/syslog"
)

var SyslogPriorityMap = map[int]log.Level{
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

	entry := log.NewEntry(nil)
	entry = entry.WithField("hostname", from)
	entry.Time = syslogMsg.Timestemp
	entry.Level = SyslogPriorityMap[syslogMsg.Severity]
	entry.Message = syslogMsg.Content
	return entry
}
