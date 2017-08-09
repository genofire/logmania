package syslog

import "github.com/genofire/logmania/log"

var SyslogPriorityMap = map[uint]log.LogLevel{
	0: log.PanicLevel,
	1: log.PanicLevel,
	2: log.PanicLevel,
	3: log.ErrorLevel,
	4: log.WarnLevel,
	5: log.InfoLevel,
	6: log.InfoLevel,
	7: log.DebugLevel,
}

func toLogEntry(logParts map[string]interface{}) *log.Entry {
	severityID := uint(logParts["severity"].(int))
	level := SyslogPriorityMap[severityID]

	if _, ok := logParts["content"]; ok {
		return &log.Entry{
			Level:    level,
			Hostname: logParts["hostname"].(string),
			Service:  logParts["tag"].(string),
			Text:     logParts["content"].(string),
		}
	}

	return &log.Entry{
		Level:    level,
		Hostname: logParts["hostname"].(string),
		Service:  logParts["app_name"].(string),
		Text:     logParts["message"].(string),
	}
}
