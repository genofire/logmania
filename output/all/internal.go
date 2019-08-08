package all

import (
	"github.com/bdlm/log"
	"time"

	"dev.sum7.eu/sum7/logmania/bot"
	"dev.sum7.eu/sum7/logmania/database"
	"dev.sum7.eu/sum7/logmania/output"
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
				logger.Warnf("notify not send to %s: [%d] %s", to.Address(), c.Level, c.Message)
			}
		}
	}
}

func (out *Output) Send(e *log.Entry, to *database.Notify) bool {
	before := time.Now()

	logger := log.WithFields(e.Data)
	logger = logger.WithField("msg", e.Message)

	logger.Debugf("starting forward message")

	out.channelNotify <- e

	after := time.Now()
	delta := after.Sub(before)
	logger.WithField("ms", float64(delta)/float64(time.Millisecond)).Debugf("end forward message")

	return true
}

func (out *Output) Close() {
	for _, item := range out.list {
		item.Close()
	}
}
