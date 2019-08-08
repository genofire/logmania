package xmpp

import (
	"strings"

	"github.com/bdlm/log"
	"gosrc.io/xmpp"
	"gosrc.io/xmpp/stanza"

	"dev.sum7.eu/sum7/logmania/database"
)

func (out *Output) Join(to string) {
	toJID, err := xmpp.NewJid(to)
	if err != nil {
		logger.Errorf("jid not generate to join muc %s : %s", to, err)
		return
	}
	toJID.Resource = nickname

	if err = out.client.Send(stanza.Presence{Attrs: stanza.Attrs{To: toJID.Full()},
		Extensions: []stanza.PresExtension{
			stanza.MucPresence{
				History: stanza.History{MaxStanzas: stanza.NewNullableInt(0)},
			}},
	}); err != nil {
		logger.Errorf("muc not join %s : %s", toJID.Full(), err)
	} else {
		out.channels[to] = true
	}
}

func (out *Output) Send(e *log.Entry, to *database.Notify) bool {
	if out.client == nil {
		logger.Error("xmpp not connected (yet)")
		return false
	}
	html, text := formatLog(e)
	if html == "" || text == "" {
		logger.Error("during format notify")
		return false
	}
	html = strings.TrimRight(to.RunReplace(html), "\n")
	text = strings.TrimRight(to.RunReplace(text), "\n")

	msg := stanza.Message{
		Attrs: stanza.Attrs{
			To: to.To,
		},
		Body: text,
		Extensions: []stanza.MsgExtension{
			stanza.HTML{Body: stanza.HTMLBody{InnerXML: html}},
		},
	}
	if to.Protocol == protoGroup {
		if _, ok := out.channels[to.To]; ok {
			out.Join(to.To)
		}
		msg.Type = stanza.MessageTypeGroupchat
		if err := out.client.Send(msg); err != nil {
			logger.WithFields(map[string]interface{}{
				"muc":  to.To,
				"text": text,
			}).Errorf("log message not forwarded: %s", err)
		}
		return true
	}
	if to.Protocol == proto {
		msg.Type = stanza.MessageTypeChat
		if err := out.client.Send(msg); err != nil {
			logger.WithFields(map[string]interface{}{
				"user": to.To,
				"text": text,
			}).Errorf("log message not forwarded: %s", err)
		}
		return true
	}
	return false
}
