package notify

import (
	"github.com/genofire/logmania/bot"
	"github.com/genofire/logmania/lib"
	"github.com/genofire/logmania/log"
	configNotify "github.com/genofire/logmania/notify/config"
)

var NotifyRegister []NotifyInit

type Notifier interface {
	Send(entry *log.Entry)
	Close()
}

type NotifyInit func(*lib.NotifyConfig, *configNotify.NotifyState, *bot.Bot) Notifier

func AddNotifier(n NotifyInit) {
	NotifyRegister = append(NotifyRegister, n)
}
