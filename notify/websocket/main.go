package xmpp

import (
	"net/http"

	"dev.sum7.eu/genofire/golang-lib/websocket"
	log "github.com/sirupsen/logrus"

	"dev.sum7.eu/genofire/logmania/bot"
	"dev.sum7.eu/genofire/logmania/database"
	"dev.sum7.eu/genofire/logmania/lib"
	"dev.sum7.eu/genofire/logmania/notify"
)

const (
	proto = "ws"
)

var logger = log.WithField("notify", proto)

type Notifier struct {
	notify.Notifier
	ws        *websocket.Server
	formatter log.Formatter
}

func Init(config *lib.NotifyConfig, db *database.DB, bot *bot.Bot) notify.Notifier {
	inputMSG := make(chan *websocket.Message)
	ws := websocket.NewServer(inputMSG, nil)

	http.HandleFunc("/ws", ws.Handler)
	http.Handle("/", http.FileServer(http.Dir(config.Websocket.Webroot)))

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

	srv := &http.Server{
		Addr: config.Websocket.Address,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			panic(err)
		}
	}()

	logger.Info("startup")
	return &Notifier{
		ws: ws,
		formatter: &log.TextFormatter{
			DisableTimestamp: true,
		},
	}
}

func (n *Notifier) Send(e *log.Entry, to *database.Notify) bool {
	if to.Protocol != proto {
		return false
	}
	n.ws.SendAll(&websocket.Message{
		Subject: to.Address(),
		Body:    e,
	})
	return true
}

func (n *Notifier) Close() {
}

func init() {
	notify.AddNotifier(Init)
}
