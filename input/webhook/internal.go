package webhook

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

type WebhookHandler func(http.Header, interface{}) *log.Entry

var handlers = make(map[string]WebhookHandler)

func AddHandler(name string, f WebhookHandler) {
	handlers[name] = f
}
