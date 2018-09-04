package notify

import (
	log "github.com/sirupsen/logrus"

	"dev.sum7.eu/genofire/logmania/bot"
	"dev.sum7.eu/genofire/logmania/database"
	"dev.sum7.eu/genofire/logmania/lib"
)

var NotifyRegister []NotifyInit

type Notifier interface {
	Default() []*database.Notify
	Send(entry *log.Entry, to *database.Notify) bool
	Close()
}

type NotifyInit func(*lib.NotifyConfig, *database.DB, *bot.Bot) Notifier

func AddNotifier(n NotifyInit) {
	NotifyRegister = append(NotifyRegister, n)
}
