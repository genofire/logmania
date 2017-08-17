package journald_json

import (
	"encoding/json"
	"strconv"

	"github.com/genofire/logmania/log"
)

type JournalMessage struct {
	Cursor                   string `json:"__CURSOR"`
	RealtimeTimestamp        string `json:"__REALTIME_TIMESTAMP"`
	MonotonicTimestamp       string `json:"__MONOTONIC_TIMESTAMP"`
	TimestampMonotonic       string `json:"TIMESTAMP_MONOTONIC"`
	TimestampBoottime        string `json:"TIMESTAMP_BOOTTIME"`
	SourceMonotonicTimestamp string `json:"_SOURCE_MONOTONIC_TIMESTAMP"`

	UID       string `json:"_UID"`
	GID       string `json:"_GID"`
	Transport string `json:"_TRANSPORT"`

	Priority         string `json:"PRIORITY"`
	SyslogFacility   string `json:"SYSLOG_FACILITY"`
	SyslogIdentifier string `json:"SYSLOG_IDENTIFIER"`

	SystemdCGroup       string `json:"_SYSTEMD_CGROUP"`
	SystemdUnit         string `json:"_SYSTEMD_UNIT"`
	SystemdSlice        string `json:"_SYSTEMD_SLICE"`
	SystemdInvocationID string `json:"_SYSTEMD_INVOCATION_ID"`

	BootID    string `json:"_BOOT_ID"`
	MachineID string `json:"_MACHINE_ID"`
	Hostname  string `json:"_HOSTNAME"`
	Message   string `json:"MESSAGE"`
}

var PriorityMap = map[int]log.LogLevel{
	0: log.PanicLevel, // emerg
	1: log.PanicLevel, // alert
	2: log.PanicLevel, // crit
	3: log.ErrorLevel, // err
	4: log.WarnLevel,  // warn
	5: log.InfoLevel,  // notice
	6: log.InfoLevel,  // info
	7: log.DebugLevel, // debug
}

func toLogEntry(msg []byte, from string) *log.Entry {
	data := &JournalMessage{}
	mapEntry := make(map[string]interface{})
	err := json.Unmarshal(msg, data)
	if err != nil {
		return nil
	}
	err = json.Unmarshal(msg, mapEntry)
	prio, err := strconv.Atoi(data.Priority)
	if err != nil {
		return nil
	}
	e := &log.Entry{
		Level:    PriorityMap[prio],
		Hostname: from,
		Service:  data.SyslogIdentifier,
		Text:     data.Message,
		Fields:   mapEntry,
	}
	if data.SystemdUnit == "" {
		e.Service = data.SystemdUnit
	}
	return e
}
