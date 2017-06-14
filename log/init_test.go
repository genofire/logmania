package log

import (
	"net/http"
	"testing"

	"golang.org/x/net/websocket"

	"github.com/stretchr/testify/assert"
)

func TestGetIP(t *testing.T) {
	assert := assert.New(t)

	req, _ := http.NewRequest("GET", "https://google.com/lola/duda?q=wasd", nil)
	req.RemoteAddr = "127.0.0.1"
	assert.Equal("127.0.0.1", getIP(req))
	req.Header.Set("X-Forwarded-For", "8.8.8.8")
	assert.Equal("8.8.8.8", getIP(req))

}

func TestHTTP(t *testing.T) {
	assert := assert.New(t)

	req, _ := http.NewRequest("GET", "https://google.com/lola/duda?q=wasd", nil)
	entry := HTTP(req)
	_, ok := entry.Fields["remote"]

	assert.NotNil(ok, "remote address not set in logger")
	assert.Equal("GET", entry.Fields["method"], "method not set in logger")
	assert.Equal("/lola/duda?q=wasd", entry.Fields["url"], "path not set in logger")
}

func TestWebsocketX(t *testing.T) {
	assert := assert.New(t)
	ws := &websocket.Conn{}
	entry := WebsocketX(ws)
	_, ok := entry.Fields["remote"]

	assert.NotNil(ok, "remote address not set in logger")
	assert.True(entry.Fields["websocket"].(bool))
}
