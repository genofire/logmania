package lib

import (
	"net/http"
	"sync"

	"github.com/genofire/logmania/log"
)

// little httpserver to handle reconnect and http.Handlers
type HTTPServer struct {
	srv               *http.Server
	Addr              string
	Handler           http.Handler
	errorNoPanic      bool
	errorNoPanicAsync sync.Mutex
}

// start httpserver
func (hs *HTTPServer) Start() {
	hs.srv = &http.Server{
		Addr:    hs.Addr,
		Handler: hs.Handler,
	}
	go func() {
		log.Debug("startup of http listener")
		if err := hs.srv.ListenAndServe(); err != nil {
			if hs.errorNoPanic {
				log.Debug("httpserver shutdown without panic")
				return
			}
			log.Panic(err)
		}
	}()
}

// rebind httpserver to a new address (e.g. new configuration)
func (hs *HTTPServer) Rebind(addr string) bool {
	if addr == hs.Addr {
		return false
	}
	hs.errorNoPanicAsync.Lock()
	hs.errorNoPanic = true
	hs.Close()
	hs.Addr = addr
	hs.Start()
	hs.errorNoPanic = false
	hs.errorNoPanicAsync.Unlock()
	return true
}

// close/stop current httpserver
func (hs *HTTPServer) Close() {
	log.Debug("startup of http listener")
	hs.srv.Close()
}
