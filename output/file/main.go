package file

import (
	"os"
	"path"
	"regexp"

	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"

	"dev.sum7.eu/genofire/logmania/bot"
	"dev.sum7.eu/genofire/logmania/database"
	"dev.sum7.eu/genofire/logmania/output"
)

const (
	proto = "file"
)

var logger = log.WithField("output", proto)

type Output struct {
	output.Output
	defaults  []*database.Notify
	files     map[string]*os.File
	formatter log.Formatter
	path      string
}

type OutputConfig struct {
	Directory string `mapstructure:"directory"`
	Default   string `mapstructure:"default"`
}

func Init(configInterface interface{}, db *database.DB, bot *bot.Bot) output.Output {
	var config OutputConfig
	if err := mapstructure.Decode(configInterface, &config); err != nil {
		logger.Warnf("not able to decode data: %s", err)
		return nil
	}
	if config.Directory == "" {
		return nil
	}
	logger.WithField("directory", config.Directory).Info("startup")

	var defaults []*database.Notify
	if config.Default != "" {
		defaults = append(defaults, &database.Notify{
			Protocol:  proto,
			To:        config.Default,
			RegexIn:   make(map[string]*regexp.Regexp),
			MaxPrioIn: log.DebugLevel,
		})
	}

	return &Output{
		defaults:  defaults,
		files:     make(map[string]*os.File),
		formatter: &log.JSONFormatter{},
		path:      config.Directory,
	}
}

func (out *Output) Default() []*database.Notify {
	return out.defaults
}

func (out *Output) getFile(name string) *os.File {
	if file, ok := out.files[name]; ok {
		return file
	}
	if m, err := regexp.MatchString(`^[0-9A-Za-z_-]*$`, name); err != nil || !m {
		logger.Errorf("not allowed to use '%s:%s'", proto, name)
		return nil
	}
	filename := path.Join(out.path, name+".json")
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		logger.Errorf("could not open file: %s", err.Error())
		return nil
	}
	out.files[name] = file
	return file
}

func (out *Output) Send(e *log.Entry, to *database.Notify) bool {
	if to.Protocol != proto {
		return false
	}
	byteText, err := out.formatter.Format(e)
	if err != nil {
		return false
	}
	text := to.RunReplace(string(byteText))
	file := out.getFile(to.To)
	if file == nil {
		return false
	}
	_, err = file.WriteString(text)
	return err == nil
}

func init() {
	output.Add(proto, Init)
}
