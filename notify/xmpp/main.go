package xmpp

import (
	"errors"
	"strings"

	xmpp_client "dev.sum7.eu/genofire/yaja/client"
	xmpp "dev.sum7.eu/genofire/yaja/xmpp"
	"dev.sum7.eu/genofire/yaja/xmpp/base"
	log "github.com/sirupsen/logrus"

	"dev.sum7.eu/genofire/logmania/bot"
	"dev.sum7.eu/genofire/logmania/database"
	"dev.sum7.eu/genofire/logmania/lib"
	"dev.sum7.eu/genofire/logmania/notify"
)

const (
	proto      = "xmpp:"
	protoGroup = "xmpp-muc:"
	nickname   = "logmania"
)

var logger = log.WithField("notify", proto)

type Notifier struct {
	notify.Notifier
	client    *xmpp_client.Client
	channels  map[string]bool
	db        *database.DB
	formatter log.Formatter
}

func Init(config *lib.NotifyConfig, db *database.DB, bot *bot.Bot) notify.Notifier {
	channels := make(map[string]bool)

	client, err := xmpp_client.NewClient(xmppbase.NewJID(config.XMPP.JID), config.XMPP.Password)
	if err != nil {
		logger.Error(err)
		return nil
	}
	go func() {
		for {
			if err := client.Start(); err != nil {
				log.Warn("close connection, try reconnect")
				client.Connect(config.XMPP.Password)
			} else {
				log.Warn("closed connection")
				return
			}
		}
	}()
	go func() {
		for {
			element, more := client.Recv()
			if !more {
				log.Warn("could not recieve new message, try later")
				continue
			}

			switch element.(type) {
			case *xmpp.PresenceClient:
				pres := element.(*xmpp.PresenceClient)
				sender := pres.From
				logPres := logger.WithField("from", sender.Full())
				switch pres.Type {
				case xmpp.PresenceTypeSubscribe:
					logPres.Debugf("recv presence subscribe")
					pres.Type = xmpp.PresenceTypeSubscribed
					pres.To = sender
					pres.From = nil
					client.Send(pres)
					logPres.Debugf("accept new subscribe")

					pres.Type = xmpp.PresenceTypeSubscribe
					pres.ID = ""
					client.Send(pres)
					logPres.Info("request also subscribe")
				case xmpp.PresenceTypeSubscribed:
					logPres.Info("recv presence accepted subscribe")
				case xmpp.PresenceTypeUnsubscribe:
					logPres.Info("recv presence remove subscribe")
				case xmpp.PresenceTypeUnsubscribed:
					logPres.Info("recv presence removed subscribe")
				case xmpp.PresenceTypeUnavailable:
					logPres.Debug("recv presence unavailable")
				case "":
					logPres.Debug("recv empty presence, maybe from joining muc")
					continue
				default:
					logPres.Warnf("recv presence unsupported: %s -> %s", pres.Type, xmpp.XMLChildrenString(pres))
				}
			case *xmpp.MessageClient:
				msg := element.(*xmpp.MessageClient)
				from := msg.From.Bare().String()
				if msg.Type == xmpp.MessageTypeGroupchat {
					from = protoGroup + from
				} else {
					from = proto + from
				}

				bot.Handle(func(answer string) {
					err := client.Send(&xmpp.MessageClient{
						Type: msg.Type,
						To:   msg.From,
						Body: answer,
					})
					if err != nil {
						logger.Error("xmpp to ", msg.From.String(), " error:", err)
					}
				}, from, msg.Body)
			}
		}
	}()
	for _, toAddresses := range db.HostTo {
		for to, _ := range toAddresses {
			toAddr := strings.TrimPrefix(to, protoGroup)
			toJID := xmppbase.NewJID(toAddr)
			toJID.Resource = nickname
			err := client.Send(&xmpp.PresenceClient{
				To: toJID,
			})
			if err != nil {
				logger.Error("xmpp could not join ", toJID.String(), " error:", err)
			} else {
				channels[toAddr] = true
			}
		}
	}
	logger.Info("startup")
	return &Notifier{
		channels: channels,
		client:   client,
		db:       db,
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
				toJID := xmppbase.NewJID(toAddr)
				toJID.Resource = nickname
				err := n.client.Send(&xmpp.PresenceClient{
					To: toJID,
				})
				if err != nil {
					logger.Error("xmpp could not join ", toJID.String(), " error:", err)
				} else {
					n.channels[toAddr] = true
				}
			}
			err := n.client.Send(&xmpp.MessageClient{
				Type: xmpp.MessageTypeGroupchat,
				To:   xmppbase.NewJID(toAddr),
				Body: string(text),
			})
			if err != nil {
				logger.Error("xmpp to ", toAddr, " error:", err)
			}
		} else {
			toAddr = strings.TrimPrefix(toAddr, proto)
			err := n.client.Send(&xmpp.MessageClient{
				Type: xmpp.MessageTypeChat,
				To:   xmppbase.NewJID(toAddr),
				Body: string(text),
			})
			if err != nil {
				logger.Error("xmpp to ", to, " error:", err)
			}
		}
	}
	return nil
}

func (n *Notifier) Close() {
	for jid := range n.channels {
		toJID := xmppbase.NewJID(jid)
		toJID.Resource = nickname
		err := n.client.Send(&xmpp.PresenceClient{
			To:   toJID,
			Type: xmpp.PresenceTypeUnavailable,
		})
		if err != nil {
			logger.Error("xmpp could not leave ", toJID.String(), " error:", err)
		}
	}
	n.client.Close()
}

func init() {
	notify.AddNotifier(Init)
}
