package xmpp

import "github.com/genofire/logmania/database"

func FormatEntry(e *database.Entry) string {
	return e.Text
}
