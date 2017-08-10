package syslog

import (
	"regexp"
	"strconv"

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
	re := regexp.MustCompile("<([0-9]*)>(.*)")
	match := re.FindStringSubmatch(string(msg))

	if len(match) <= 1 {
		return &log.Entry{
			Level:    log.DebugLevel,
			Text:     string(msg),
			Hostname: from,
		}
	}
	v, _ := strconv.Atoi(match[1])
	prio := v % 8
	text := match[2]

	return &log.Entry{
		Level:    SyslogPriorityMap[prio],
		Text:     text,
		Hostname: from,
	}
}
