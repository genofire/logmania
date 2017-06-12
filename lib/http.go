package lib

import (
	"net/http"

	"github.com/genofire/logmania/log"
)

type HTTPServer struct {
	srv          *http.Server
	ErrorNoPanic bool
	Addr         string
	Handler      http.Handler
}

func (hs *HTTPServer) Start() {
	hs.srv = &http.Server{
		Addr:    hs.Addr,
		Handler: hs.Handler,
	}
	go func() {
		log.Debug("startup of http listener")
		if err := hs.srv.ListenAndServe(); err != nil {
			if hs.ErrorNoPanic {
				log.Debug("httpserver shutdown without panic")
				return
			}
			log.Panic(err)
		}
	}()
}
func (hs *HTTPServer) Close() {
	log.Debug("startup of http listener")
	hs.srv.Close()
}
