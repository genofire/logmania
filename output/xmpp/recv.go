package xmpp

import (
	xmpp "dev.sum7.eu/genofire/yaja/xmpp"
	"github.com/bdlm/log"
)

func (out *Output) receiver() {
	for {
		element, more := out.client.Recv()
		if !more {
			log.Warn("could not receive new message, try later")
			continue
		}
		out.recv(element)
	}
}
func (out *Output) recv(element interface{}) {

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
			out.botOut <- pres
			logPres.Debugf("accept new subscribe")

			pres.Type = xmpp.PresenceTypeSubscribe
			pres.ID = ""
			out.botOut <- pres
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

		answer := out.bot.Handle(from, msg.Body)
		if answer == "" {
			return
		}
		to := msg.From
		if msg.Type == xmpp.MessageTypeGroupchat && !to.IsBare() {
			to = to.Bare()
		}
		out.botOut <- &xmpp.MessageClient{
			Type: msg.Type,
			To:   to,
			Body: answer,
		}
	}
}
