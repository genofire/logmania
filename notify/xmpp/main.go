package xmpp

import (
	"github.com/genofire/logmania/database"
	"github.com/genofire/logmania/lib"
	"github.com/genofire/logmania/log"
	"github.com/genofire/logmania/notify"
	xmpp "github.com/mattn/go-xmpp"
)

type Notifier struct {
	notify.Notifier
	client *xmpp.Client
}

func NotifyInit(config *lib.NotifyConfig) notify.Notifier {
	options := xmpp.Options{
		Host:          config.XMPP.Host,
		User:          config.XMPP.Username,
		Password:      config.XMPP.Password,
		NoTLS:         config.XMPP.NoTLS,
		Debug:         config.XMPP.Debug,
		Session:       config.XMPP.Session,
		Status:        config.XMPP.Status,
		StatusMessage: config.XMPP.StatusMessage,
	}
	client, err := options.NewClient()
	if err != nil {
		return nil
	}
	return &Notifier{client: client}
}

func (n *Notifier) Send(e *database.Entry) {
	users := database.UserByApplication(e.ApplicationID)
	for _, user := range users {
		if user.NotifyXMPP && log.LogLevel(e.Level) >= user.NotifyAfterLoglevel {
			n.client.SendHtml(xmpp.Chat{Remote: user.XMPP, Type: "chat", Text: formatEntry(e)})
		}
	}
}

func init() {
	notify.AddNotifier(NotifyInit)
}
