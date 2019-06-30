package xmpp

import (
	"regexp"

	xmpp_client "dev.sum7.eu/genofire/yaja/client"
	xmpp "dev.sum7.eu/genofire/yaja/xmpp"
	"dev.sum7.eu/genofire/yaja/xmpp/base"
	"github.com/bdlm/log"
	"github.com/mitchellh/mapstructure"

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
	defaults []*database.Notify
	channels map[string]bool
	bot      *bot.Bot
	client   *xmpp_client.Client
	botOut   chan interface{}
	logOut   chan interface{}
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

	jid := xmppbase.NewJID(config.JID)
	client, err := xmpp_client.NewClient(jid, config.Password)

	if err != nil {
		logger.Error(err)
		return nil
	}
	go func() {
		client.Start()
		log.Panic("closed connection")
	}()
	out := &Output{
		channels: make(map[string]bool),
		bot:      bot,
		client:   client,
		botOut:   make(chan interface{}),
		logOut:   make(chan interface{}),
	}
	go out.sender()
	go out.receiver()

	logger.WithField("jid", config.JID).Info("startup")

	for to, muc := range config.Defaults {
		var def *database.Notify
		pro := proto
		if muc {
			pro = protoGroup
		}
		if dbNotify, ok := db.NotifiesByAddress[pro+":"+to]; ok {
			def = dbNotify
		} else {
			def = &database.Notify{
				Protocol:  pro,
				To:        to,
				RegexIn:   make(map[string]*regexp.Regexp),
				MaxPrioIn: log.DebugLevel,
			}
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

func (out *Output) Default() []*database.Notify {
	return out.defaults
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
