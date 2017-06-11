package main

import (
	"time"

	"github.com/genofire/logmania/log"
	logOutput "github.com/genofire/logmania/log/hook/output"
)

func main() {
	log.Info("startup")
	log.New().AddField("answer", 42).AddFields(map[string]interface{}{"answer": 3, "foo": "bar"}).Warn("Some spezial")
	log.Debug("Never shown up")
	logOutput.ShowTime = false
	logOutput.AboveLevel = log.DebugLevel
	log.Debugf("Startup %v", time.Now())
	logOutput.ShowTime = true
	log.Panic("let it crash")
}
