package notify

import (
	log "github.com/sirupsen/logrus"

	"github.com/genofire/logmania/bot"
	"github.com/genofire/logmania/lib"
	configNotify "github.com/genofire/logmania/notify/config"
)

var NotifyRegister []NotifyInit

type Notifier interface {
	Fire(entry *log.Entry) error
	Levels() []log.Level
	Close()
}

type NotifyInit func(*lib.NotifyConfig, *configNotify.NotifyState, *bot.Bot) Notifier

func AddNotifier(n NotifyInit) {
	NotifyRegister = append(NotifyRegister, n)
}
