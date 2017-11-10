package xmpp

import (
	"errors"
	"fmt"
	"io"
	"strings"

	xmpp "github.com/mattn/go-xmpp"
	log "github.com/sirupsen/logrus"

	"github.com/genofire/logmania/bot"
	"github.com/genofire/logmania/database"
	"github.com/genofire/logmania/lib"
	"github.com/genofire/logmania/notify"
)

const (
	proto      = "xmpp:"
	protoGroup = "xmpp-muc:"
	nickname   = "logmania"
)

var logger = log.WithField("notify", proto)

type Notifier struct {
	notify.Notifier
	client    *xmpp.Client
	channels  map[string]bool
	db        *database.DB
	formatter *log.TextFormatter
}

func Init(config *lib.NotifyConfig, db *database.DB, bot *bot.Bot) notify.Notifier {
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
				if err == io.EOF {
					client, err = options.NewClient()
					log.Warn("reconnect")
					if err != nil {
						log.Panic(err)
					}
					continue
				}
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
	for _, toAddresses := range db.HostTo {
		for to, _ := range toAddresses {
			toAddr := strings.TrimPrefix(to, protoGroup)
			client.JoinMUCNoHistory(toAddr, nickname)
		}
	}
	logger.Info("startup")
	return &Notifier{
		client: client,
		db:     db,
		formatter: &log.TextFormatter{
			DisableTimestamp: true,
		},
	}
}

func (n *Notifier) Send(e *log.Entry) error {
	e, to := n.db.SendTo(e)
	if to == nil {
		return errors.New("no reciever found")
	}
	text, err := n.formatter.Format(e)
	if err != nil {
		return err
	}
	for _, toAddr := range to {
		if strings.HasPrefix(toAddr, protoGroup) {
			toAddr = strings.TrimPrefix(toAddr, protoGroup)
			if _, ok := n.channels[toAddr]; ok {
				n.client.JoinMUCNoHistory(toAddr, nickname)
			}
			_, err = n.client.SendHtml(xmpp.Chat{Remote: toAddr, Type: "groupchat", Text: string(text)})
			if err != nil {
				logger.Error("xmpp to ", to, " error:", err)
			}
		} else {
			toAddr = strings.TrimPrefix(toAddr, proto)
			_, err := n.client.SendHtml(xmpp.Chat{Remote: toAddr, Type: "chat", Text: string(text)})
			if err != nil {
				logger.Error("xmpp to ", to, " error:", err)
			}
		}
	}
	return nil
}

func (n *Notifier) Close() {
	for jid := range n.channels {
		n.client.LeaveMUC(jid)
	}
	n.client.Close()
}

func init() {
	notify.AddNotifier(Init)
}
