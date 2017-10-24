package syslog

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type SyslogMessage struct {
	Timestemp time.Time
	Hostname  string
	Tag       string
	Service   int
	Content   string
	Facility  int
	Severity  int
}

func Parse(binaryMsg []byte) *SyslogMessage {
	var err error

	msg := &SyslogMessage{}

	re := regexp.MustCompile("<([0-9]*)>(.*)")
	match := re.FindStringSubmatch(string(binaryMsg))

	prio, _ := strconv.Atoi(match[1])
	msg.Facility = prio / 8
	msg.Severity = prio % 8

	timeLength := len(time.RFC3339)
	if len(match[2]) > timeLength {
		msg.Timestemp, err = time.Parse(time.RFC3339, match[2][:timeLength])
		if err != nil {
			timeLength = 0
		}
	}
	timeLength = len(msg.Timestemp.Format(time.Stamp))
	if len(match[2]) > timeLength {
		msg.Timestemp, err = time.Parse(time.Stamp, match[2][:timeLength])
		if err != nil {
			timeLength = 0
		}
	}

	msg.Content = strings.TrimLeft(match[2][timeLength:], " ")

	/*
	 TODO: detect other parts in content
	 - Hostname (if exists)
	 - Tag
	 - Service
	*/

	return msg
}

func (msg *SyslogMessage) Priority() int {
	return msg.Facility*8 + msg.Severity
}

func (msg *SyslogMessage) Dump() []byte {
	result := fmt.Sprintf("<%d>%s %s %s[%d]: %s",
		msg.Priority(),
		msg.Timestemp.Format(time.RFC3339),
		msg.Hostname,
		msg.Tag,
		msg.Service,
		msg.Content,
	)
	if !strings.HasSuffix(msg.Content, "\n") {
		result = fmt.Sprintf("%s\n", result)
	}
	return []byte(result)
}
