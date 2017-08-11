package xmpp

import (
	"fmt"
	"strings"

	xmpp "github.com/mattn/go-xmpp"

	"github.com/genofire/logmania/bot"
	"github.com/genofire/logmania/lib"
	"github.com/genofire/logmania/log"
	"github.com/genofire/logmania/notify"
	configNotify "github.com/genofire/logmania/notify/config"
)

type Notifier struct {
	notify.Notifier
	client *xmpp.Client
	state  *configNotify.NotifyState
}

func Init(config *lib.NotifyConfig, state *configNotify.NotifyState, bot *bot.Bot) notify.Notifier {
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
	go func() {
		for {
			chat, err := client.Recv()
			if err != nil {
				log.Warn(err)
			}
			switch v := chat.(type) {
			case xmpp.Chat:
				bot.Handle(func(answer string) {
					client.SendHtml(xmpp.Chat{Remote: v.Remote, Type: "chat", Text: answer})
				}, fmt.Sprintf("xmpp:%s", strings.Split(v.Remote, "/")[0]), v.Text)
			}
		}
	}()
	log.Info("xmpp startup")
	return &Notifier{client: client, state: state}
}

func (n *Notifier) Send(e *log.Entry) {
	to := n.state.SendTo(e)
	if to == nil {
		return
	}
	for _, toAddr := range to {
		to := strings.TrimPrefix(toAddr, "xmpp:")
		if strings.Contains(toAddr, "conference") || strings.Contains(toAddr, "irc") {
			n.client.JoinMUCNoHistory(to, "logmania")
			_, err := n.client.SendHtml(xmpp.Chat{Remote: to, Type: "groupchat", Text: formatEntry(e)})
			if err != nil {
				fmt.Println("xmpp to ", to, " error:", err)
			}
		} else {
			_, err := n.client.SendHtml(xmpp.Chat{Remote: to, Type: "chat", Text: formatEntry(e)})
			if err != nil {
				fmt.Println("xmpp to ", to, " error:", err)
			}
		}
	}
}

func (n *Notifier) Close() {}

func init() {
	notify.AddNotifier(Init)
}
