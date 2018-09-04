package output

import (
	log "github.com/sirupsen/logrus"

	"dev.sum7.eu/genofire/logmania/bot"
	"dev.sum7.eu/genofire/logmania/database"
)

var Register = make(map[string]Init)

type Output interface {
	Default() []*database.Notify
	Send(entry *log.Entry, to *database.Notify) bool
	Close()
}

type Init func(interface{}, *database.DB, *bot.Bot) Output

func Add(name string, init Init) {
	Register[name] = init
}
