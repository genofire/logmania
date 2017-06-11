package client

import (
	"fmt"

	"github.com/gorilla/websocket"

	"github.com/genofire/logmania/log"
)

type Logger struct {
	AboveLevel log.LogLevel
	conn       *websocket.Conn
}

func (l *Logger) hook(e *log.Entry) {
	if e.Level < l.AboveLevel {
		return
	}
	err := l.conn.WriteJSON(e)
	if err != nil {
		log.Panic("[logmania] could not send token")
	}
}
func (l *Logger) Close() {
	l.conn.Close()
}

func Init(url, token string) *Logger {
	logger := &Logger{
		AboveLevel: log.InfoLevel,
	}
	c, _, err := websocket.DefaultDialer.Dial(fmt.Sprint(url, "/logger"), nil)
	if err != nil {
		log.Panic("[logmania] error on connect")
		return nil
	}
	err = c.WriteJSON(token)
	if err != nil {
		log.Panic("[logmania] could not send token")
		return nil
	}
	logger.conn = c
	log.AddHook(logger.hook)
	return logger
}
