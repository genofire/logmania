package main

import (
	"github.com/genofire/logmania/database"
	"github.com/genofire/logmania/log"

	logOutput "github.com/genofire/logmania/log/hook/output"
)

type SelfLogger struct {
	log.Logger
	AboveLevel log.LogLevel
	lastMsg    string
	lastTime   int
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
	if l.lastMsg == e.Text {
		l.lastTime += 1
	} else {
		l.lastMsg = e.Text
		l.lastTime = 1
	}
	dbEntry := database.InsertEntry("", e)
	if dbEntry != nil && notifier != nil {
		notifier.Send(dbEntry)
	} else {
		l := logOutput.NewLogger()
		e := log.New()
		e.Text = "No notifier found"
		e.Level = log.WarnLevel
		l.Hook(e)
	}
}

func (l *SelfLogger) Close() {
}
