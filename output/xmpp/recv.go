package xmpp

import (
	"time"

	"gosrc.io/xmpp"
	"gosrc.io/xmpp/stanza"
)

func (out *Output) recvMessage(s xmpp.Sender, p stanza.Packet) {
	before := time.Now()

	msg, ok := p.(stanza.Message)
	if !ok {
		logger.Errorf("blame gosrc.io/xmpp for routing: %s", p)
		return
	}
	logger.WithFields(map[string]interface{}{
		"sender":  msg.From,
		"request": msg.Body,
	}).Debug("handling bot message")

	from, err := xmpp.NewJid(msg.From)
	if err != nil {
		logger.Errorf("blame gosrc.io/xmpp for jid encoding: %s", msg.From)
		return
	}

	fromBare := from.Bare()
	fromLogmania := ""
	if msg.Type == stanza.MessageTypeGroupchat {
		fromLogmania = protoGroup + ":" + fromBare
	} else {
		fromLogmania = proto + ":" + fromBare
	}

	answer := out.bot.Handle(fromLogmania, msg.Body)
	if answer == "" {
		return
	}
	if err := s.Send(stanza.Message{Attrs: stanza.Attrs{To: fromBare, Type: msg.Type}, Body: answer}); err != nil {
		logger.WithFields(map[string]interface{}{
			"sender":  fromLogmania,
			"request": msg.Body,
			"answer":  answer,
		}).Errorf("unable to send bot answer: %s", err)
	}

	after := time.Now()
	delta := after.Sub(before)

	logger.WithFields(map[string]interface{}{
		"sender":  fromLogmania,
		"request": msg.Body,
		"answer":  answer,
		"ms":      float64(delta) / float64(time.Millisecond),
	}).Debug("handled xmpp bot message")
}

func (out *Output) recvPresence(s xmpp.Sender, p stanza.Packet) {
	pres, ok := p.(stanza.Presence)
	if !ok {
		logger.Errorf("blame gosrc.io/xmpp for routing: %s", p)
		return
	}
	from, err := xmpp.NewJid(pres.From)
	if err != nil {
		logger.Errorf("blame gosrc.io/xmpp for jid encoding: %s", pres.From)
		return
	}
	fromBare := from.Bare()
	logPres := logger.WithField("from", from)

	switch pres.Type {
	case stanza.PresenceTypeSubscribe:
		logPres.Debugf("recv presence subscribe")
		if err := s.Send(stanza.Presence{Attrs: stanza.Attrs{
			Type: stanza.PresenceTypeSubscribed,
			To:   fromBare,
			Id:   pres.Id,
		}}); err != nil {
			logPres.WithField("user", pres.From).Errorf("answer of subscribe not send: %s", err)
			return
		}
		logPres.Debugf("accept new subscribe")

		if err := s.Send(stanza.Presence{Attrs: stanza.Attrs{
			Type: stanza.PresenceTypeSubscribe,
			To:   fromBare,
		}}); err != nil {
			logPres.WithField("user", pres.From).Errorf("request of subscribe not send: %s", err)
			return
		}
		logPres.Info("request also subscribe")
	case stanza.PresenceTypeSubscribed:
		logPres.Info("recv presence accepted subscribe")
	case stanza.PresenceTypeUnsubscribe:
		logPres.Info("recv presence remove subscribe")
	case stanza.PresenceTypeUnsubscribed:
		logPres.Info("recv presence removed subscribe")
	case stanza.PresenceTypeUnavailable:
		logPres.Debug("recv presence unavailable")
	case "":
		logPres.Debug("recv empty presence, maybe from joining muc")
		return
	default:
		logPres.Warnf("recv presence unsupported: %s -> %v", pres.Type, pres)
	}
}
