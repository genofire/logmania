// logger to bind at github.com/genofire/logmania/log.AddLogger to send log entries to logmania server
package client

import (
	"fmt"

	"github.com/gorilla/websocket"

	"github.com/genofire/logmania/log"
)

// client logger
type Logger struct {
	log.Logger
	AboveLevel log.LogLevel
	conn       *websocket.Conn
	closed     bool
}

// CurrentLogger (for override settings e.g. AboveLevel)
var CurrentLogger *Logger

// create a new logmania client logger
func NewLogger(url, token string, AboveLevel log.LogLevel) *Logger {
	c, _, err := websocket.DefaultDialer.Dial(fmt.Sprint(url, "/logger"), nil)
	if err != nil {
		log.Error("[logmania] error on connect: ", err)
		return nil
	}
	err = c.WriteMessage(websocket.TextMessage, []byte(token))
	if err != nil {
		log.Error("[logmania] could not send token:", err)
		return nil
	}
	return &Logger{
		AboveLevel: AboveLevel,
		conn:       c,
	}
}

// handle a log entry (send to logmania server)
func (l *Logger) Hook(e *log.Entry) {
	if l.closed {
		return
	}
	if e.Level < l.AboveLevel {
		return
	}
	err := l.conn.WriteJSON(e)
	if err != nil {
		log.Error("[logmania] could not send log entry:", err)
		l.Close()
	}
}

// Listen if logmania server want to close the connection
func (l *Logger) Listen() {
	for {
		msgType, _, err := l.conn.ReadMessage()
		if msgType == -1 {
			l.closed = true
			l.conn.Close()
			return
		}
		if err != nil {
			log.Warn("[logmania] close listener:", err)
		}
	}
}

// close connection to logger
func (l *Logger) Close() {
	l.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	l.closed = true
}

// init logmania's client logger and bind
func Init(url, token string, AboveLevel log.LogLevel) *Logger {
	CurrentLogger = NewLogger(url, token, AboveLevel)
	go CurrentLogger.Listen()
	log.AddLogger(CurrentLogger)
	return CurrentLogger
}
