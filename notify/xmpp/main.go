package xmpp

import (
	"errors"
	"fmt"
	"strings"

	xmpp "github.com/mattn/go-xmpp"
	log "github.com/sirupsen/logrus"

	"github.com/genofire/logmania/bot"
	"github.com/genofire/logmania/lib"
	"github.com/genofire/logmania/notify"
	configNotify "github.com/genofire/logmania/notify/config"
)

var logger = log.WithField("notify", "xmpp")

type Notifier struct {
	notify.Notifier
	client    *xmpp.Client
	state     *configNotify.NotifyState
	formatter *log.TextFormatter
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
		logger.Error(err)
		return nil
	}
	go func() {
		for {
			chat, err := client.Recv()
			if err != nil {
				logger.Warn(err)
			}
			switch v := chat.(type) {
			case xmpp.Chat:
				bot.Handle(func(answer string) {
					client.SendHtml(xmpp.Chat{Remote: v.Remote, Type: "chat", Text: answer})
				}, fmt.Sprintf("xmpp:%s", strings.Split(v.Remote, "/")[0]), v.Text)
			}
		}
	}()
	logger.Info("startup")
	return &Notifier{
		client: client,
		state:  state,
		formatter: &log.TextFormatter{
			DisableTimestamp: true,
		},
	}
}

func (n *Notifier) Fire(e *log.Entry) error {
	to := n.state.SendTo(e)
	if to == nil {
		return errors.New("no reciever found")
	}
	text, err := n.formatter.Format(e)
	if err != nil {
		return err
	}
	for _, toAddr := range to {
		to := strings.TrimPrefix(toAddr, "xmpp:")
		if strings.Contains(toAddr, "conference") || strings.Contains(toAddr, "irc") {
			n.client.JoinMUCNoHistory(to, "logmania")
			_, err = n.client.SendHtml(xmpp.Chat{Remote: to, Type: "groupchat", Text: string(text)})
			if err != nil {
				logger.Error("xmpp to ", to, " error:", err)
			}
		} else {
			_, err := n.client.SendHtml(xmpp.Chat{Remote: to, Type: "chat", Text: string(text)})
			if err != nil {
				logger.Error("xmpp to ", to, " error:", err)
			}
		}
	}
	return nil
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

func (n *Notifier) Close() {}

func init() {
	notify.AddNotifier(Init)
}
