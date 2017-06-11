package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/genofire/logmania/lib"
	"github.com/genofire/logmania/log"
	_ "github.com/genofire/logmania/log/hook/output"
)

var (
	configPath string
	config     *lib.Config
)

func main() {
	flag.StringVar(&configPath, "config", "logmania.conf", "config file")
	log.Info("starting logmania")
	config, err := lib.ReadConfig(configPath)
	if config == nil || err != nil {
		log.Panicf("Could not load '%s' for configuration.", configPath)
	}

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
		log.Errorf("Could not load '%s' for new configuration. Skip reload.", configPath)
		return
	}
}
