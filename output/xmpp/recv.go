package xmpp

import (
	"github.com/bdlm/log"
	"gosrc.io/xmpp"
)

func (out *Output) recvMessage(s xmpp.Sender, p xmpp.Packet) {
	msg, ok := p.(xmpp.Message)
	if !ok {
		log.Errorf("blame gosrc.io/xmpp for routing: %s", p)
		return
	}
	from, err := xmpp.NewJid(msg.From)
	if err != nil {
		log.Errorf("blame gosrc.io/xmpp for jid encoding: %s", msg.From)
		return
	}
	fromBare := from.Bare()
	fromLogmania := ""
	if msg.Type == xmpp.MessageTypeGroupchat {
		fromLogmania = protoGroup + ":" + fromBare
	} else {
		fromLogmania = proto + ":" + fromBare
	}

	answer := out.bot.Handle(fromLogmania, msg.Body)
	if answer == "" {
		return
	}
	reply := xmpp.Message{Attrs: xmpp.Attrs{To: fromBare, Type: msg.Type}, Body: answer}
	s.Send(reply)
}

func (out *Output) recvPresence(s xmpp.Sender, p xmpp.Packet) {
	pres, ok := p.(xmpp.Presence)
	if !ok {
		log.Errorf("blame gosrc.io/xmpp for routing: %s", p)
		return
	}
	from, err := xmpp.NewJid(pres.From)
	if err != nil {
		log.Errorf("blame gosrc.io/xmpp for jid encoding: %s", pres.From)
		return
	}
	fromBare := from.Bare()
	logPres := logger.WithField("from", from)

	switch pres.Type {
	case xmpp.PresenceTypeSubscribe:
		logPres.Debugf("recv presence subscribe")
		s.Send(xmpp.Presence{Attrs: xmpp.Attrs{
			Type: xmpp.PresenceTypeSubscribed,
			To:   fromBare,
			Id:   pres.Id,
		}})
		logPres.Debugf("accept new subscribe")

		s.Send(xmpp.Presence{Attrs: xmpp.Attrs{
			Type: xmpp.PresenceTypeSubscribe,
			To:   fromBare,
		}})
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
		return
	default:
		logPres.Warnf("recv presence unsupported: %s -> %s", pres.Type, pres)
	}
}
