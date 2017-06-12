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

func HTTP(r *http.Request) *Entry {
	return New().AddFields(map[string]interface{}{
		"remote": getIP(r),
		"method": r.Method,
		"url":    r.URL.RequestURI(),
	})
}

func WebsocketX(ws *websocket.Conn) *Entry {
	r := ws.Request()
	return New().AddFields(map[string]interface{}{
		"remote":    getIP(r),
		"websocket": true,
		"url":       r.URL.RequestURI(),
	})
}
func WebsocketGozilla(ws *wsGozilla.Conn) *Entry {
	return New().AddFields(map[string]interface{}{
		"remote":    ws.RemoteAddr().String(),
		"websocket": true,
	})
}
