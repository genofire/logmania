package main

import (
	"github.com/genofire/logmania/database"
	"github.com/genofire/logmania/log"
)

type SelfLogger struct {
	log.Logger
	AboveLevel log.LogLevel
	lastMsg string
	lastTime int
}

func NewSelfLogger() *SelfLogger {
	return &SelfLogger{
		AboveLevel: log.InfoLevel,
	}
}

func (l *SelfLogger) Hook(e *log.Entry) {
	if e.Level >= l.AboveLevel {
		return
	}
	// TODO strange logger
	if l.lastTime > 15 {
		panic("selflogger same log to oftern")
	}
	if l.lastMsg == e.Text{
		l.lastTime += 1
	} else {
		l.lastMsg = e.Text
		l.lastTime = 1
	}
	database.InsertEntry("",e)
}


func (l *SelfLogger) Close() {
}
