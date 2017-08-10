package xmpp

import (
	"github.com/genofire/logmania/lib"
	"github.com/genofire/logmania/log"
	"github.com/genofire/logmania/notify"
	xmpp "github.com/mattn/go-xmpp"
)

type Notifier struct {
	notify.Notifier
	client *xmpp.Client
	state  *notify.NotifyState
}

func Init(config *lib.NotifyConfig, state *notify.NotifyState) notify.Notifier {
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
	return &Notifier{client: client, state: state}
}

func (n *Notifier) Send(e *log.Entry) {
	to := n.state.SendTo(e)
	if to == nil {
		return
	}
	for _, to := range to {
		n.client.SendHtml(xmpp.Chat{Remote: to, Type: "chat", Text: formatEntry(e)})
	}
}

func (n *Notifier) Close() {}

func init() {
	notify.AddNotifier(Init)
}
