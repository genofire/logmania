package all

import (
	"github.com/genofire/logmania/lib"
	"github.com/genofire/logmania/log"
	"github.com/genofire/logmania/notify"
)

type Notifier struct {
	notify.Notifier
	list          []notify.Notifier
	channelNotify chan *log.Entry
}

func Init(config *lib.NotifyConfig, state *notify.NotifyState) notify.Notifier {
	var list []notify.Notifier
	for _, init := range notify.NotifyRegister {
		notify := init(config, state)

		if notify == nil {
			continue
		}
		list = append(list, notify)
	}

	n := &Notifier{
		list:          list,
		channelNotify: make(chan *log.Entry),
	}
	go n.sender()
	return n
}

func (n *Notifier) sender() {
	for c := range n.channelNotify {
		for _, item := range n.list {
			item.Send(c)
		}
	}
}

func (n *Notifier) Send(e *log.Entry) {
	n.channelNotify <- e
}

func (n *Notifier) Close() {
	for _, item := range n.list {
		item.Close()
	}
}
