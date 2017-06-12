package recieve

import (
	"encoding/json"
	"net/http"

	"github.com/genofire/logmania/database"
	"github.com/genofire/logmania/log"
	"github.com/gorilla/websocket"
)

type Handler struct {
	upgrader websocket.Upgrader
}

func NewHandler() *Handler {
	return &Handler{
		upgrader: websocket.Upgrader{},
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logEntry := log.HTTP(r)
	c, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		logEntry.Warn("no webservice upgrade:", err)
		return
	}
	token := ""
	defer c.Close()
	for {
		if token == "" {
			var maybeToken string
			msgType, msg, err := c.ReadMessage()
			if err != nil {
				logEntry.Error("recieving token", err)
				break
			}
			if msgType != websocket.TextMessage {
				logEntry.Warn("recieve no token")
				break
			}
			maybeToken = string(msg)
			logEntry.AddField("token", maybeToken)
			if !database.IsTokenValid(maybeToken) {
				logEntry.Warn("recieve wrong token")
				break
			} else {
				token = maybeToken
				logEntry.Info("recieve valid token")
			}
			continue
		}
		var entry log.Entry
		msgType, msg, err := c.ReadMessage()
		if msgType == -1 {
			c.Close()
			logEntry.Info("connecting closed")
			break
		}
		if err != nil {
			logEntry.Error("recieving log entry:", err)
			break
		}
		err = json.Unmarshal(msg, &entry)
		if err != nil {
			logEntry.Error("umarshal log entry:", err)
			break
		}
		database.InsertEntry(token, &entry)
	}
}
