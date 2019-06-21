package xmpp

import (
	"encoding/xml"
	"strings"

	xmpp "dev.sum7.eu/genofire/yaja/xmpp"
	"dev.sum7.eu/genofire/yaja/xmpp/base"
	"dev.sum7.eu/genofire/yaja/xmpp/x/muc"
	"github.com/bdlm/log"

	"dev.sum7.eu/genofire/logmania/database"
)

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

func (out *Output) sender() {
	// priority of bot higher: https://groups.google.com/forum/#!topic/golang-nuts/M2xjN_yWBiQ
	for {
		select {
		case el := <-out.botOut:
			out.client.Send(el)
		default:
			select {
			case el := <-out.logOut:
				out.client.Send(el)
			default:
			}
		}
	}
}
func (out *Output) Send(e *log.Entry, to *database.Notify) bool {
	html, text := formatLog(e)
	if html == "" || text == "" {
		logger.Error("during format notify")
		return false
	}
	html = strings.TrimRight(to.RunReplace(html), "\n")
	var body xmpp.XMLElement
	xml.Unmarshal([]byte(html), &body)

	text = strings.TrimRight(to.RunReplace(text), "\n")

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
		out.logOut <- &xmpp.MessageClient{
			Type: xmpp.MessageTypeGroupchat,
			To:   xmppbase.NewJID(to.To),
			Body: text,
			HTML: &xmpp.HTML{Body: xmpp.HTMLBody{Body: body}},
		}
		return true
	}
	if to.Protocol == proto {
		out.logOut <- &xmpp.MessageClient{
			Type: xmpp.MessageTypeChat,
			To:   xmppbase.NewJID(to.To),
			Body: text,
			HTML: &xmpp.HTML{Body: xmpp.HTMLBody{Body: body}},
		}
		return true
	}
	return false
}
