package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/genofire/logmania/api/recieve"
	"github.com/genofire/logmania/database"
	"github.com/genofire/logmania/lib"
	"github.com/genofire/logmania/log"
	logOutput "github.com/genofire/logmania/log/hook/output"
)

var (
	configPath string
	config     *lib.Config
	api        *lib.HTTPServer
	apiNoPanic *bool
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
	log.AddLogger(logger)

	api = &lib.HTTPServer{
		Addr:    config.API.Bind,
		Handler: recieve.NewHandler(),
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
	if config.API.Bind != api.Addr {
		api.ErrorNoPanic = true
		api.Close()
		api.Addr = config.API.Bind
		api.Start()
		log.Info("reload: new api bind")
	}
	if database.ReplaceConnect(config.Database.Type, config.Database.Connect) {
		log.Info("reload: new database connection establish")
	}
}
