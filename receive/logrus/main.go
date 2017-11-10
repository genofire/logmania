package logrus

import (
	"net/http"

	"github.com/genofire/golang-lib/websocket"
	log "github.com/sirupsen/logrus"

	"github.com/genofire/logmania/lib"
	"github.com/genofire/logmania/receive"
)

const WS_LOG_ENTRY = "log"

var logger = log.WithField("receive", "logrus")

type Receiver struct {
	receive.Receiver
	input         chan *websocket.Message
	exportChannel chan *log.Entry
	serverSocket  *websocket.Server
}

func Init(config *lib.ReceiveConfig, exportChannel chan *log.Entry) receive.Receiver {
	inputMsg := make(chan *websocket.Message)
	ws := websocket.NewServer(inputMsg, websocket.NewSessionManager())

	http.HandleFunc("/receiver", ws.Handler)

	recv := &Receiver{
		input:         inputMsg,
		serverSocket:  ws,
		exportChannel: exportChannel,
	}

	logger.Info("init")

	return recv
}

func (rc *Receiver) Listen() {
	logger.Info("listen")
	for msg := range rc.input {
		if event, ok := msg.Body.(log.Entry); ok {
			rc.exportChannel <- &event
		}
	}
}

func (rc *Receiver) Close() {
}

func init() {
	receive.AddReceiver("websocket", Init)
}
