package all

import (
	log "github.com/sirupsen/logrus"

	"dev.sum7.eu/genofire/logmania/bot"
	"dev.sum7.eu/genofire/logmania/database"
	"dev.sum7.eu/genofire/logmania/output"
)

var logger = log.WithField("notify", "all")

type Output struct {
	output.Output
	list          []output.Output
	db            *database.DB
	channelNotify chan *log.Entry
}

func Init(configInterface interface{}, db *database.DB, bot *bot.Bot) output.Output {
	config := configInterface.(map[string]interface{})

	var list []output.Output

	for outputType, init := range output.Register {
		configForItem := config[outputType]
		if configForItem == nil {
			log.Warnf("the input type '%s' has no configuration\n", outputType)
			continue
		}
		notify := init(configForItem, db, bot)

		if notify == nil {
			continue
		}
		list = append(list, notify)
		def := notify.Default()
		if def == nil {
			continue
		}
		db.DefaultNotify = append(db.DefaultNotify, def...)
	}

	out := &Output{
		db:            db,
		list:          list,
		channelNotify: make(chan *log.Entry),
	}
	go out.sender()
	return out
}

func (out *Output) sender() {
	for c := range out.channelNotify {
		e, _, tos := out.db.SendTo(c)
		for _, to := range tos {
			send := false
			for _, item := range out.list {
				send = item.Send(e, to)
				if send {
					break
				}
			}
			if !send {
				logger.Warnf("notify not send to %s: [%s] %s", to.Address(), c.Level.String(), c.Message)
			}
		}
	}
}

func (out *Output) Send(e *log.Entry, to *database.Notify) bool {
	out.channelNotify <- e
	return true
}

func (out *Output) Close() {
	for _, item := range out.list {
		item.Close()
	}
}
