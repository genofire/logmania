// logmania Server
//
// reload config with SIGUSR1
//
//   Usage of logmania:
//    -config string
//      	config file (default "logmania.conf")
//    -debug
//      	enable debuging
package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/genofire/logmania/api/receive"
	"github.com/genofire/logmania/database"
	"github.com/genofire/logmania/lib"
	"github.com/genofire/logmania/log"
	logOutput "github.com/genofire/logmania/log/hook/output"
	"github.com/genofire/logmania/notify"
	"github.com/genofire/logmania/notify/all"
)

var (
	configPath string
	config     *lib.Config
	api        *lib.HTTPServer
	apiNoPanic *bool
	notifier   notify.Notifier
	debug      bool
)

func main() {
	flag.StringVar(&configPath, "config", "logmania.conf", "config file")
	flag.BoolVar(&debug, "debug", false, "enable debuging")
	flag.Parse()
	logger := NewSelfLogger()

	if debug {
		logger.AboveLevel = log.DebugLevel
		logOutput.AboveLevel = log.DebugLevel
	}

	log.Info("starting logmania")

	config, err := lib.ReadConfig(configPath)
	if config == nil || err != nil {
		log.Panicf("Could not load '%s' for configuration.", configPath)
	}

	database.Connect(config.Database.Type, config.Database.Connect)
	log.AddLogger("selflogger", logger)

	notifier = all.NotifyInit(&config.Notify)

	api = &lib.HTTPServer{
		Addr:    config.API.Bind,
		Handler: receive.NewHandler(notifier),
	}
	api.Start()

	// Wait for system signal
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGUSR1)
	for sig := range sigchan {
		switch sig {
		case syscall.SIGTERM:
			log.Warn("terminated of logmania")
			os.Exit(0)
		case syscall.SIGQUIT:
			quit()
		case syscall.SIGHUP:
			quit()
		case syscall.SIGUSR1:
			reload()
		}
	}
}

func quit() {
	notifier.Close()
	log.Info("quit of logmania")
	os.Exit(0)
}

func reload() {
	log.Info("reload config file")
	config, err := lib.ReadConfig(configPath)
	if config == nil || err != nil {
		log.Errorf("reload: could not load '%s' for new configuration. Skip reload.", configPath)
		return
	}
	if api.Rebind(config.API.Bind) {
		log.Info("reload: new api bind")
	}
	if database.ReplaceConnect(config.Database.Type, config.Database.Connect) {
		log.Info("reload: new database connection establish")
	}
}
