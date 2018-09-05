package webhook

import (
	"fmt"
	"net/http"

	libHTTP "github.com/genofire/golang-lib/http"
	log "github.com/sirupsen/logrus"

	"dev.sum7.eu/genofire/logmania/input"
)

const InputType = "webhook"

var logger = log.WithField("input", InputType)

type Input struct {
	input.Input
	exportChannel chan *log.Entry
}

func Init(config interface{}, exportChannel chan *log.Entry) input.Input {
	logger.Info("init")

	return &Input{
		exportChannel: exportChannel,
	}
}

func (in *Input) getHTTPHandler(name string, h WebhookHandler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var body interface{}
		libHTTP.Read(r, &body)

		e := h(r.Header, body)
		if e == nil {
			http.Error(w, fmt.Sprintf("no able to generate log for handler-request %s", name), http.StatusInternalServerError)
			return
		}
		in.exportChannel <- e
		http.Error(w, fmt.Sprintf("handler-request %s - ok", name), http.StatusOK)
	}
}

func (in *Input) Listen() {
	for name, h := range handlers {
		http.HandleFunc("/input/"+InputType+"/"+name, in.getHTTPHandler(name, h))
	}
}

func (in *Input) Close() {
}

func init() {
	input.Add(InputType, Init)
}
