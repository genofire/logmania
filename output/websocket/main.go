package xmpp

import (
	"net/http"

	"dev.sum7.eu/genofire/golang-lib/websocket"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"

	"dev.sum7.eu/genofire/logmania/bot"
	"dev.sum7.eu/genofire/logmania/database"
	"dev.sum7.eu/genofire/logmania/output"
)

const (
	proto = "ws"
)

var logger = log.WithField("output", proto)

type Output struct {
	output.Output
	defaults  []*database.Notify
	ws        *websocket.Server
	formatter log.Formatter
}

type OutputConfig struct {
	Default string `mapstructure:"default"`
}

func Init(configInterface interface{}, db *database.DB, bot *bot.Bot) output.Output {
	var config OutputConfig
	if err := mapstructure.Decode(configInterface, &config); err != nil {
		logger.Warnf("not able to decode data: %s", err)
		return nil
	}
	inputMSG := make(chan *websocket.Message)
	ws := websocket.NewServer(inputMSG, nil)

	http.HandleFunc("/output/ws", ws.Handler)

	go func() {
		for msg := range inputMSG {
			if msg.Subject != "bot" {
				logger.Warnf("receive unknown websocket message: %s", msg.Subject)
				continue
			}
			bot.Handle(func(answer string) {
				msg.Answer("bot", answer)
			}, "", msg.Body.(string))
		}
	}()

	logger.Info("startup")

	var defaults []*database.Notify
	if config.Default != "" {
		defaults = append(defaults, &database.Notify{
			Protocol: proto,
			To:       config.Default,
		})
	}
	return &Output{
		defaults: defaults,
		ws:       ws,
		formatter: &log.TextFormatter{
			DisableTimestamp: true,
		},
	}
}

func (out *Output) Default() []*database.Notify {
	return out.defaults
}

func (out *Output) Send(e *log.Entry, to *database.Notify) bool {
	if to.Protocol != proto {
		return false
	}

	out.ws.SendAll(&websocket.Message{
		Subject: to.Address(),
		Body: &log.Entry{
			Buffer:  e.Buffer,
			Data:    e.Data,
			Level:   e.Level,
			Logger:  e.Logger,
			Message: to.RunReplace(e.Message),
			Time:    e.Time,
		},
	})
	return true
}

func (out *Output) Close() {
}

func init() {
	output.Add("websocket", Init)
}
