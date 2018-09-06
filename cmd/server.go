package cmd

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NYTimes/gziphandler"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"dev.sum7.eu/genofire/golang-lib/file"
	"dev.sum7.eu/genofire/golang-lib/worker"
	"dev.sum7.eu/genofire/logmania/bot"
	"dev.sum7.eu/genofire/logmania/database"
	"dev.sum7.eu/genofire/logmania/input"
	allInput "dev.sum7.eu/genofire/logmania/input/all"
	"dev.sum7.eu/genofire/logmania/lib"
	"dev.sum7.eu/genofire/logmania/output"
	allOutput "dev.sum7.eu/genofire/logmania/output/all"
)

var (
	configPath   string
	config       *lib.Config
	db           *database.DB
	dbSaveWorker *worker.Worker
	out          output.Output
	in           input.Input
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
		if config.Debug {
			log.SetLevel(log.DebugLevel)
		}

		db = database.ReadDBFile(config.DB)
		go func() { dbSaveWorker = file.NewSaveJSONWorker(time.Minute, config.DB, db) }()

		logmaniaBot = bot.NewBot(db)

		out = allOutput.Init(config.Output, db, logmaniaBot)
		logChannel = make(chan *log.Entry)

		go func() {
			for a := range logChannel {
				out.Send(a, nil)
			}
		}()

		if config.AlertCheck.Duration > time.Duration(time.Second) {
			go db.Alert(config.AlertCheck.Duration, out.Send)
		}

		log.WithField("defaults", len(db.DefaultNotify)).Info("starting logmania")

		if config.HTTPAddress != "" {
			if config.Webroot != "" {
				http.Handle("/", gziphandler.GzipHandler(http.FileServer(http.Dir(config.Webroot))))
			}

			srv := &http.Server{
				Addr: config.HTTPAddress,
			}

			go func() {
				if err := srv.ListenAndServe(); err != http.ErrServerClosed {
					log.Panic(err)
				}
			}()
		}

		in = allInput.Init(config.Input, logChannel)

		go in.Listen()

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
	in.Close()
	out.Close()
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
	if config.Debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
	in.Close()
	in = allInput.Init(config.Input, logChannel)
	go in.Listen()

	out.Close()
	out = allOutput.Init(config.Output, db, logmaniaBot)
}

func init() {
	RootCmd.AddCommand(serverCmd)
	serverCmd.Flags().StringVarP(&configPath, "config", "c", "logmania.conf", "Path to configuration file")
}
