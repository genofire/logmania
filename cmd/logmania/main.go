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

	"github.com/genofire/logmania/lib"
	log "github.com/genofire/logmania/log"
	"github.com/genofire/logmania/notify"
	allNotify "github.com/genofire/logmania/notify/all"
	"github.com/genofire/logmania/receive"
	allReceiver "github.com/genofire/logmania/receive/all"
)

var (
	configPath string
	config     *lib.Config
	notifier   notify.Notifier
	receiver   receive.Receiver
	logChannel chan *log.Entry
)

func main() {
	flag.StringVar(&configPath, "config", "logmania.conf", "config file")
	flag.Parse()

	log.Info("starting logmania")

	config, err := lib.ReadConfig(configPath)
	if config == nil || err != nil {
		log.Panicf("Could not load '%s' for configuration.", configPath)
	}

	notifier = allNotify.Init(&config.Notify)
	log.Save = notifier.Send
	go func() {
		for a := range logChannel {
			notifier.Send(a)
		}
	}()

	receiver = allReceiver.Init(&config.Receive, logChannel)

	go receiver.Listen()

	// Wait for system signal
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGUSR1)
	for sig := range sigchan {
		switch sig {
		case syscall.SIGTERM:
			log.Panic("terminated of logmania")
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
	receiver.Close()
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
	receiver.Close()
	receiver = allReceiver.Init(&config.Receive, logChannel)
	go receiver.Listen()

	notifier.Close()
	notifier = allNotify.Init(&config.Notify)
}
