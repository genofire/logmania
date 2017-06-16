package all

import (
	"github.com/genofire/logmania/database"
	"github.com/genofire/logmania/lib"
	"github.com/genofire/logmania/notify"
)

type Notifier struct {
	notify.Notifier
	list []notify.Notifier
}

func NotifyInit(config *lib.NotifyConfig) notify.Notifier {
	var list []notify.Notifier
	for _, init := range notify.NotifyRegister {
		notify := init(config)

		if notify == nil {
			continue
		}
		list = append(list, notify)
	}
	return &Notifier{
		list: list,
	}
}

func (n *Notifier) Send(entry *database.Entry) {
	for _, item := range n.list {
		item.Send(entry)
	}
}

func (n *Notifier) Close() {
	for _, item := range n.list {
		item.Close()
	}
}
