package logrus

import (
	"net/http"

	"dev.sum7.eu/genofire/golang-lib/websocket"
	log "github.com/sirupsen/logrus"

	"dev.sum7.eu/genofire/logmania/input"
)

const inputType = "logrus"
const WS_LOG_ENTRY = "log"

var logger = log.WithField("input", inputType)

type Input struct {
	input.Input
	input         chan *websocket.Message
	exportChannel chan *log.Entry
	serverSocket  *websocket.Server
}

func Init(config interface{}, exportChannel chan *log.Entry) input.Input {
	inputMsg := make(chan *websocket.Message)
	ws := websocket.NewServer(inputMsg, websocket.NewSessionManager())

	http.HandleFunc("/input/"+inputType, ws.Handler)

	in := &Input{
		input:         inputMsg,
		serverSocket:  ws,
		exportChannel: exportChannel,
	}

	logger.Info("init")

	return in
}

func (in *Input) Listen() {
	logger.Info("listen")
	for msg := range in.input {
		if event, ok := msg.Body.(log.Entry); ok {
			in.exportChannel <- &event
		}
	}
}

func (in *Input) Close() {
}

func init() {
	input.Add(inputType, Init)
}
