package all

import (
	log "github.com/sirupsen/logrus"

	"github.com/genofire/logmania/bot"
	"github.com/genofire/logmania/lib"
	"github.com/genofire/logmania/notify"
	configNotify "github.com/genofire/logmania/notify/config"
)

type Notifier struct {
	notify.Notifier
	list          []notify.Notifier
	channelNotify chan *log.Entry
}

func Init(config *lib.NotifyConfig, state *configNotify.NotifyState, bot *bot.Bot) notify.Notifier {
	var list []notify.Notifier
	for _, init := range notify.NotifyRegister {
		notify := init(config, state, bot)

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
			item.Fire(c)
		}
	}
}

func (n *Notifier) Fire(e *log.Entry) error {
	n.channelNotify <- e
	return nil
}

func (n *Notifier) Close() {
	for _, item := range n.list {
		item.Close()
	}
}

func (n *Notifier) Levels() []log.Level {
	return []log.Level{
		log.DebugLevel,
		log.InfoLevel,
		log.WarnLevel,
		log.ErrorLevel,
		log.FatalLevel,
		log.PanicLevel,
	}
}
