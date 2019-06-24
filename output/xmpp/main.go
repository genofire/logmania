package xmpp

import (
	"regexp"

	"github.com/bdlm/log"
	"github.com/mitchellh/mapstructure"
	"gosrc.io/xmpp"

	"dev.sum7.eu/genofire/logmania/bot"
	"dev.sum7.eu/genofire/logmania/database"
	"dev.sum7.eu/genofire/logmania/output"
)

const (
	proto      = "xmpp"
	protoGroup = "xmpp-muc"
	nickname   = "logmania"
)

var logger = log.WithField("output", proto)

type Output struct {
	output.Output
	defaults []*database.Notify
	channels map[string]bool
	bot      *bot.Bot
	client   *xmpp.Client
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

	out := &Output{
		channels: make(map[string]bool),
		bot:      bot,
	}

	router := xmpp.NewRouter()
	router.HandleFunc("message", out.recvMessage)
	router.HandleFunc("presence", out.recvPresence)

	client, err := xmpp.NewClient(xmpp.Config{
		Jid:      config.JID,
		Password: config.Password,
	}, router)

	if err != nil {
		logger.Error(err)
		return nil
	}
	out.client = client
	cm := xmpp.NewStreamManager(client, nil)
	go func() {
		cm.Run()
		log.Panic("closed connection")
	}()

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
		toJID, err := xmpp.NewJid(jid)
		if err != nil {
			logger.Error("xmpp could generate jid to leave ", jid, " error:", err)
		}
		toJID.Resource = nickname
		err = out.client.Send(xmpp.Presence{Attrs: xmpp.Attrs{
			To:   toJID.Full(),
			Type: xmpp.PresenceTypeUnavailable,
		}})
		if err != nil {
			logger.Error("xmpp could not leave ", toJID.Full(), " error:", err)
		}
	}
}

func init() {
	output.Add(proto, Init)
}
