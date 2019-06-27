package xmpp

import (
	"strings"

	"github.com/bdlm/log"
	"gosrc.io/xmpp"

	"dev.sum7.eu/genofire/logmania/database"
)

func (out *Output) Join(to string) {
	toJID, err := xmpp.NewJid(to)
	if err != nil {
		logger.Error("xmpp could not generate jid to join ", to, " error:", err)
		return
	}
	toJID.Resource = nickname

	err = out.client.Send(xmpp.Presence{Attrs: xmpp.Attrs{To: toJID.Full()},
		Extensions: []xmpp.PresExtension{
			xmpp.MucPresence{
				History: xmpp.History{MaxStanzas: xmpp.NewNullableInt(0)},
			}},
	})
	if err != nil {
		logger.Error("xmpp could not join ", toJID.Full(), " error:", err)
	} else {
		out.channels[to] = true
	}
}

func (out *Output) Send(e *log.Entry, to *database.Notify) bool {
	html, text := formatLog(e)
	if html == "" || text == "" {
		logger.Error("during format notify")
		return false
	}
	html = strings.TrimRight(to.RunReplace(html), "\n")
	text = strings.TrimRight(to.RunReplace(text), "\n")

	msg := xmpp.Message{
		Attrs: xmpp.Attrs{
			To: to.To,
		},
		Body: text,
		Extensions: []xmpp.MsgExtension{
			xmpp.HTML{Body: xmpp.HTMLBody{InnerXML: html}},
		},
	}
	if to.Protocol == protoGroup {
		if _, ok := out.channels[to.To]; ok {
			out.Join(to.To)
		}
		msg.Type = xmpp.MessageTypeGroupchat
		out.client.Send(msg)
		return true
	}
	if to.Protocol == proto {
		msg.Type = xmpp.MessageTypeChat
		out.client.Send(msg)
		return true
	}
	return false
}
