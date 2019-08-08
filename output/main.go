package output

import (
	"github.com/bdlm/log"

	"dev.sum7.eu/sum7/logmania/bot"
	"dev.sum7.eu/sum7/logmania/database"
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
