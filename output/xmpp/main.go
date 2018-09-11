package xmpp

import (
	"regexp"
	"strings"

	xmpp_client "dev.sum7.eu/genofire/yaja/client"
	xmpp "dev.sum7.eu/genofire/yaja/xmpp"
	"dev.sum7.eu/genofire/yaja/xmpp/base"
	"dev.sum7.eu/genofire/yaja/xmpp/x/muc"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"

	"dev.sum7.eu/genofire/logmania/bot"
	"dev.sum7.eu/genofire/logmania/database"
	"dev.sum7.eu/genofire/logmania/output"
)

const (
	proto      = "xmpp"
	protoGroup = "xmpp-muc"
	nickname   = "logmania"
)

var historyMaxChars = 0

var logger = log.WithField("output", proto)

type Output struct {
	output.Output
	defaults  []*database.Notify
	client    *xmpp_client.Client
	channels  map[string]bool
	formatter log.Formatter
}

type OutputConfig struct {
	JID      string          `mapstructure:"jid"`
	Password string          `mapstructure:"password"`
	Defaults map[string]bool `mapstructure:"default"`
}

func Init(configInterface interface{}, db *database.DB, bot *bot.Bot) output.Output {
	var config OutputConfig
	if err := mapstructure.Decode(configInterface, &config); err != nil {
		logger.Warnf("not able to decode data: %s", err)
		return nil
	}
	channels := make(map[string]bool)

	jid := xmppbase.NewJID(config.JID)
	client := &xmpp_client.Client{
		JID:     jid,
		Logging: logger.WithField("jid", jid.String()),
	}
	err := client.Connect(config.Password)

	if err != nil {
		logger.Error(err)
		return nil
	}
	go func() {
		for {
			if err := client.Start(); err != nil {
				log.Warn("close connection, try reconnect")
				client.Connect(config.Password)
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
				log.Warn("could not receive new message, try later")
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
					from = protoGroup + ":" + from
				} else {
					from = proto + ":" + from
				}

				answer := bot.Handle(from, msg.Body)
				if answer == "" {
					continue
				}
				to := msg.From
				if msg.Type == xmpp.MessageTypeGroupchat && !to.IsBare() {
					to = to.Bare()
				}
				err := client.Send(&xmpp.MessageClient{
					Type: msg.Type,
					To:   to,
					Body: answer,
				})
				if err != nil {
					logger.Error("xmpp to ", msg.From.String(), " error:", err)
				}
			}
		}
	}()

	logger.WithField("jid", config.JID).Info("startup")

	out := &Output{
		channels: channels,
		client:   client,
		formatter: &log.TextFormatter{
			DisableTimestamp: true,
		},
	}

	for to, muc := range config.Defaults {
		def := &database.Notify{
			Protocol:  proto,
			To:        to,
			RegexIn:   make(map[string]*regexp.Regexp),
			MaxPrioIn: log.DebugLevel,
		}
		if muc {
			def.Protocol = protoGroup
			out.Join(to)
		}
		out.defaults = append(out.defaults, def)
	}
	for _, toAddresses := range db.NotifiesByAddress {
		if toAddresses.Protocol == protoGroup {
			out.Join(toAddresses.To)
		}
	}
	return out
}

func (out *Output) Join(to string) {
	toJID := xmppbase.NewJID(to)
	toJID.Resource = nickname
	err := out.client.Send(&xmpp.PresenceClient{
		To: toJID,
		MUC: &xmuc.Base{
			History: &xmuc.History{
				MaxChars: &historyMaxChars,
			},
		},
	})
	if err != nil {
		logger.Error("xmpp could not join ", toJID.String(), " error:", err)
	} else {
		out.channels[to] = true
	}
}

func (out *Output) Default() []*database.Notify {
	return out.defaults
}

func (out *Output) Send(e *log.Entry, to *database.Notify) bool {
	textByte, err := out.formatter.Format(e)
	if err != nil {
		logger.Error("during format notify", err)
		return false
	}
	text := strings.TrimRight(to.RunReplace(string(textByte)), "\n")
	if to.Protocol == protoGroup {
		if _, ok := out.channels[to.To]; ok {
			toJID := xmppbase.NewJID(to.To)
			toJID.Resource = nickname
			err := out.client.Send(&xmpp.PresenceClient{
				To: toJID,
				MUC: &xmuc.Base{
					History: &xmuc.History{
						MaxChars: &historyMaxChars,
					},
				},
			})
			if err != nil {
				logger.Error("xmpp could not join ", toJID.String(), " error:", err)
			} else {
				out.channels[to.To] = true
			}
		}
		err := out.client.Send(&xmpp.MessageClient{
			Type: xmpp.MessageTypeGroupchat,
			To:   xmppbase.NewJID(to.To),
			Body: text,
		})
		if err != nil {
			logger.Error("xmpp to ", to.To, " error:", err)
		}
		return true
	}
	if to.Protocol == proto {
		err := out.client.Send(&xmpp.MessageClient{
			Type: xmpp.MessageTypeChat,
			To:   xmppbase.NewJID(to.To),
			Body: text,
		})
		if err != nil {
			logger.Error("xmpp to ", to, " error:", err)
		}
		return true
	}
	return false
}

func (out *Output) Close() {
	for jid := range out.channels {
		toJID := xmppbase.NewJID(jid)
		toJID.Resource = nickname
		err := out.client.Send(&xmpp.PresenceClient{
			To:   toJID,
			Type: xmpp.PresenceTypeUnavailable,
		})
		if err != nil {
			logger.Error("xmpp could not leave ", toJID.String(), " error:", err)
		}
	}
	out.client.Close()
}

func init() {
	output.Add(proto, Init)
}
