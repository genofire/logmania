package all

import (
	log "github.com/sirupsen/logrus"

	"dev.sum7.eu/genofire/logmania/bot"
	"dev.sum7.eu/genofire/logmania/database"
	"dev.sum7.eu/genofire/logmania/lib"
	"dev.sum7.eu/genofire/logmania/notify"
)

var logger = log.WithField("notify", "all")

type Notifier struct {
	notify.Notifier
	list          []notify.Notifier
	db            *database.DB
	channelNotify chan *log.Entry
}

func Init(config *lib.NotifyConfig, db *database.DB, bot *bot.Bot) notify.Notifier {
	var list []notify.Notifier
	for _, init := range notify.NotifyRegister {
		notify := init(config, db, bot)

		if notify == nil {
			continue
		}
		list = append(list, notify)
	}

	n := &Notifier{
		db:            db,
		list:          list,
		channelNotify: make(chan *log.Entry),
	}
	go n.sender()
	return n
}

func (n *Notifier) sender() {
	for c := range n.channelNotify {
		e, _, tos := n.db.SendTo(c)
		for _, to := range tos {
			send := false
			for _, item := range n.list {
				send = item.Send(e, to)
				if send {
					break
				}
			}
			if !send {
				logger.Warn("notify not send to anybody: [%s] %s", c.Level.String(), c.Message)
			}
		}
	}
}

func (n *Notifier) Send(e *log.Entry, to *database.Notify) bool {
	n.channelNotify <- e
	return true
}

func (n *Notifier) Close() {
	for _, item := range n.list {
		item.Close()
	}
}
