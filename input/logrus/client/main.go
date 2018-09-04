package client

import (
	"io"

	websocketLib "dev.sum7.eu/genofire/golang-lib/websocket"
	"dev.sum7.eu/genofire/logmania/input/logrus"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

// client logger
type Logmania struct {
	URL    string
	Token  uuid.UUID
	Levels []log.Level
	quere  chan *log.Entry
	conn   *websocket.Conn
}

func NewClient(url string, token uuid.UUID, lvls ...log.Level) *Logmania {
	logger := &Logmania{
		URL:    url,
		Token:  token,
		Levels: lvls,
		quere:  make(chan *log.Entry),
	}
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Error("[logmania] error on connect: ", err)
		return nil
	}
	logger.conn = conn
	go logger.Start()
	return logger
}

// Listen if logmania server want to close the connection
func (l *Logmania) listen() {
	for {
		var msg websocketLib.Message
		err := websocket.ReadJSON(l.conn, &msg)
		if err == io.EOF {
			l.Close()
			log.Warn("[logmania] close listener:", err)
		} else if err != nil {
			log.Println(err)
		} else {
			if msg.Subject == websocketLib.SessionMessageInit {
				l.conn.WriteJSON(&websocketLib.Message{
					Subject: websocketLib.SessionMessageInit,
					ID:      l.Token,
				})
			}
		}
	}
}

func (l *Logmania) writer() {
	for e := range l.quere {
		err := l.conn.WriteJSON(&websocketLib.Message{
			Subject: logrus.WS_LOG_ENTRY,
			Body:    e,
		})
		if err != nil {
			log.Error("[logmania] could not send log entry:", err)
		}
	}
}

func (l *Logmania) Start() {
	go l.listen()
	l.writer()
}

func (l *Logmania) Fire(e *log.Entry) {
	l.quere <- e
}

// close connection to logger
func (l *Logmania) Close() {
	l.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	close(l.quere)
}
