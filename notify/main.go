package notify

import (
	"github.com/genofire/logmania/lib"
	"github.com/genofire/logmania/log"
)

var NotifyRegister []NotifyInit

type Notifier interface {
	Send(entry *log.Entry)
	Close()
}

type NotifyInit func(*lib.NotifyConfig, *NotifyState) Notifier

func AddNotifier(n NotifyInit) {
	NotifyRegister = append(NotifyRegister, n)
}
