package log

import (
	"net/http"

	wsGozilla "github.com/gorilla/websocket"
	"golang.org/x/net/websocket"
)

func getIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.RemoteAddr
	}
	return ip
}

// init log entry with extra fields of interesting http request context
func HTTP(r *http.Request) *Entry {
	return New().AddFields(map[string]interface{}{
		"remote": getIP(r),
		"method": r.Method,
		"url":    r.URL.RequestURI(),
	})
}

// init log entry with extra fields of interesting websocket request context
func WebsocketX(ws *websocket.Conn) *Entry {
	r := ws.Request()
	return New().AddFields(map[string]interface{}{
		"remote":    getIP(r),
		"websocket": true,
		"url":       r.URL.RequestURI(),
	})
}

// init log entry with extra fields of interesting websocket request context
func WebsocketGozilla(ws *wsGozilla.Conn) *Entry {
	return New().AddFields(map[string]interface{}{
		"remote":    ws.RemoteAddr().String(),
		"websocket": true,
	})
}
