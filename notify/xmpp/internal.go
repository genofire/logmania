package xmpp

import "github.com/genofire/logmania/database"

func formatEntry(e *database.Entry) string {
	return e.Text
}
