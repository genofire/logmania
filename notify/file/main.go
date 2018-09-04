package file

import (
	"os"
	"path"
	"regexp"

	log "github.com/sirupsen/logrus"

	"dev.sum7.eu/genofire/logmania/bot"
	"dev.sum7.eu/genofire/logmania/database"
	"dev.sum7.eu/genofire/logmania/lib"
	"dev.sum7.eu/genofire/logmania/notify"
)

const (
	proto = "file"
)

var logger = log.WithField("notify", proto)

type Notifier struct {
	notify.Notifier
	defaults  []*database.Notify
	files     map[string]*os.File
	formatter log.Formatter
	path      string
}

func Init(config *lib.NotifyConfig, db *database.DB, bot *bot.Bot) notify.Notifier {
	if config.File.Directory == "" {
		return nil
	}
	logger.WithField("directory", config.File.Directory).Info("startup")

	var defaults []*database.Notify
	if config.File.Default != "" {
		defaults = append(defaults, &database.Notify{
			Protocol: proto,
			To:       config.File.Default,
		})
	}

	return &Notifier{
		defaults:  defaults,
		files:     make(map[string]*os.File),
		formatter: &log.JSONFormatter{},
		path:      config.File.Directory,
	}
}

func (n *Notifier) Default() []*database.Notify {
	return n.defaults
}

func (n *Notifier) getFile(name string) *os.File {
	if file, ok := n.files[name]; ok {
		return file
	}
	if m, err := regexp.MatchString(`^[0-9A-Za-z_-]*$`, name); err != nil || !m {
		logger.Errorf("not allowed to use '%s:%s'", proto, name)
		return nil
	}
	filename := path.Join(n.path, name+".json")
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		logger.Errorf("could not open file: %s", err.Error())
		return nil
	}
	n.files[name] = file
	return file
}

func (n *Notifier) Send(e *log.Entry, to *database.Notify) bool {
	if to.Protocol != proto {
		return false
	}
	byteText, err := n.formatter.Format(e)
	if err != nil {
		return false
	}
	text := to.RunReplace(string(byteText))
	file := n.getFile(to.To)
	if file == nil {
		return false
	}
	_, err = file.WriteString(text)
	return err == nil
}

func init() {
	notify.AddNotifier(Init)
}
