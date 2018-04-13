package xmpp

import (
	"os"

	log "github.com/sirupsen/logrus"

	"dev.sum7.eu/genofire/logmania/bot"
	"dev.sum7.eu/genofire/logmania/database"
	"dev.sum7.eu/genofire/logmania/lib"
	"dev.sum7.eu/genofire/logmania/notify"
)

const (
	proto = "file:"
)

var logger = log.WithField("notify", proto)

type Notifier struct {
	notify.Notifier
	formatter log.Formatter
	file      *os.File
	path      string
}

func Init(config *lib.NotifyConfig, db *database.DB, bot *bot.Bot) notify.Notifier {
	logger.Info("startup")
	if config.File == "" {
		return nil
	}
	file, err := os.OpenFile(config.File, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		logger.Errorf("could not open file: %s", err.Error())
		return nil
	}

	return &Notifier{
		formatter: &log.JSONFormatter{},
		file:      file,
		path:      config.File,
	}
}

func (n *Notifier) Send(e *log.Entry) error {
	text, err := n.formatter.Format(e)
	if err != nil {
		return err
	}
	_, err = n.file.Write(text)
	if err != nil {
		logger.Warnf("could not write to logfile: %s - try to reopen it", err.Error())
		file, err := os.OpenFile(n.path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			return err
		}
		n.file = file
		_, err = n.file.Write(text)
	}
	return err
}

func init() {
	notify.AddNotifier(Init)
}
