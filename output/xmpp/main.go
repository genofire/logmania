package xmpp

import (
	"regexp"

	"github.com/bdlm/log"
	"github.com/mitchellh/mapstructure"
	"gosrc.io/xmpp"
	"gosrc.io/xmpp/stanza"

	"dev.sum7.eu/sum7/logmania/bot"
	"dev.sum7.eu/sum7/logmania/database"
	"dev.sum7.eu/sum7/logmania/output"
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
	client   xmpp.Sender
	botOut   chan interface{}
	logOut   chan interface{}
}

type OutputConfig struct {
	Address  string          `mapstructure:"address"`
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
		Address:  config.Address,
		Jid:      config.JID,
		Password: config.Password,
	}, router)

	if err != nil {
		logger.Error(err)
		return nil
	}
	cm := xmpp.NewStreamManager(client, func(c xmpp.Sender) {
		out.client = c

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
		logger.Info("join muc after connect")
	})
	go func() {
		cm.Run()
		log.Panic("closed connection")
	}()

	logger.WithField("jid", config.JID).Info("startup")
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
		if err = out.client.Send(stanza.Presence{Attrs: stanza.Attrs{
			To:   toJID.Full(),
			Type: stanza.PresenceTypeUnavailable,
		}}); err != nil {
			logger.Error("xmpp could not leave ", toJID.Full(), " error:", err)
		}
	}
}

func init() {
	output.Add(proto, Init)
}
