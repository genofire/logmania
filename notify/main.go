package notify

import (
	log "github.com/sirupsen/logrus"

	"github.com/genofire/logmania/bot"
	"github.com/genofire/logmania/database"
	"github.com/genofire/logmania/lib"
)

var NotifyRegister []NotifyInit

type Notifier interface {
	Send(entry *log.Entry) error
	Close()
}

type NotifyInit func(*lib.NotifyConfig, *database.DB, *bot.Bot) Notifier

func AddNotifier(n NotifyInit) {
	NotifyRegister = append(NotifyRegister, n)
}
