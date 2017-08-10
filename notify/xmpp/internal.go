package xmpp

import (
	"fmt"

	"github.com/genofire/logmania/log"
)

func formatEntry(e *log.Entry) string {
	return fmt.Sprintf("[%s] [%s] %s", e.Hostname, e.Level, e.Text)
}
