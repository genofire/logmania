package cmd

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"dev.sum7.eu/genofire/golang-lib/file"
	"dev.sum7.eu/genofire/golang-lib/worker"
	"dev.sum7.eu/genofire/logmania/bot"
	"dev.sum7.eu/genofire/logmania/database"
	"dev.sum7.eu/genofire/logmania/lib"
	"dev.sum7.eu/genofire/logmania/notify"
	allNotify "dev.sum7.eu/genofire/logmania/notify/all"
	"dev.sum7.eu/genofire/logmania/receive"
	allReceiver "dev.sum7.eu/genofire/logmania/receive/all"
)

var (
	configPath   string
	config       *lib.Config
	db           *database.DB
	dbSaveWorker *worker.Worker
	notifier     notify.Notifier
	receiver     receive.Receiver
	logChannel   chan *log.Entry
	logmaniaBot  *bot.Bot
)

// serverCmd represents the serve command
var serverCmd = &cobra.Command{
	Use:     "server",
	Short:   "Runs the logmania server",
	Example: "logmania server --config /etc/yanic.toml",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetFormatter(&log.TextFormatter{
			DisableTimestamp: true,
		})
		config := &lib.Config{}
		err := file.ReadTOML(configPath, config)
		if config == nil || err != nil {
			log.Panicf("Could not load '%s' for configuration.", configPath)
		}

		db = database.ReadDBFile(config.DB)
		go func() { dbSaveWorker = file.NewSaveJSONWorker(time.Minute, config.DB, db) }()

		logmaniaBot = bot.NewBot(db)

		notifier = allNotify.Init(&config.Notify, db, logmaniaBot)
		logChannel = make(chan *log.Entry)

		go func() {
			for a := range logChannel {
				notifier.Send(a, nil)
			}
		}()

		if config.Notify.AlertCheck.Duration > time.Duration(time.Second) {
			go db.Alert(config.Notify.AlertCheck.Duration, notifier.Send)
		}

		log.Info("starting logmania")

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
	},
}

func quit() {
	dbSaveWorker.Close()
	file.SaveJSON(config.DB, db)
	receiver.Close()
	notifier.Close()
	log.Info("quit of logmania")
	os.Exit(0)
}

func reload() {
	log.Info("reload config file")
	var config lib.Config
	err := file.ReadTOML(configPath, &config)
	if err != nil {
		log.Errorf("reload: could not load '%s' for new configuration. Skip reload.", configPath)
		return
	}
	receiver.Close()
	receiver = allReceiver.Init(&config.Receive, logChannel)
	go receiver.Listen()

	notifier.Close()
	notifier = allNotify.Init(&config.Notify, db, logmaniaBot)
}

func init() {
	RootCmd.AddCommand(serverCmd)
	serverCmd.Flags().StringVarP(&configPath, "config", "c", "logmania.conf", "Path to configuration file")
}
