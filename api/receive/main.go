// receiver of log entry over network (websocket)
package receive

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/genofire/logmania/database"
	"github.com/genofire/logmania/log"
	"github.com/genofire/logmania/notify"
)

// http.Handler for init network
type Handler struct {
	http.Handler
	upgrader websocket.Upgrader
	Notify   notify.Notifier
}

// init new Handler
func NewHandler(notifyHandler notify.Notifier) *Handler {
	return &Handler{
		upgrader: websocket.Upgrader{},
		Notify:   notifyHandler,
	}
}

// server response of handler
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
				logEntry.Error("receiving token", err)
				break
			}
			if msgType != websocket.TextMessage {
				logEntry.Warn("receive no token")
				break
			}
			maybeToken = string(msg)
			logEntry.AddField("token", maybeToken)
			if !database.IsTokenValid(maybeToken) {
				logEntry.Warn("receive wrong token")
				break
			} else {
				token = maybeToken
				logEntry.Info("receive valid token")
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
			logEntry.Error("receiving log entry:", err)
			break
		}
		err = json.Unmarshal(msg, &entry)
		if err != nil {
			logEntry.Error("umarshal log entry:", err)
			break
		}
		dbEntry := database.InsertEntry(token, &entry)
		if dbEntry != nil && h.Notify != nil {
			h.Notify.Send(dbEntry)
		}
	}
}
