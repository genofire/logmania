package xmpp

import (
	"fmt"

	"github.com/genofire/logmania/log"
)

func formatEntry(e *log.Entry) string {
	if e.Hostname != "" && e.Service != "" {
		return fmt.Sprintf("[%s-%s] [%s] %s", e.Hostname, e.Service, e.Level, e.Text)
		
	} else if e.Hostname != "" {
		return fmt.Sprintf("[%s] [%s] %s", e.Hostname, e.Level, e.Text)
		
	} else if e.Service != "" {
		return fmt.Sprintf("[%s] [%s] %s", e.Service, e.Level, e.Text)
		
	}
	return fmt.Sprintf("[%s] %s", e.Level, e.Text)
}
