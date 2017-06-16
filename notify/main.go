package notify

import (
	"github.com/genofire/logmania/database"
	"github.com/genofire/logmania/lib"
)

var NotifyRegister []NotifyInit

type Notifier interface {
	Send(entry *database.Entry)
	Close()
}

type NotifyInit func(*lib.NotifyConfig) Notifier

func AddNotifier(n NotifyInit) {
	NotifyRegister = append(NotifyRegister, n)
}

func Start(config *lib.NotifyConfig) {
}
