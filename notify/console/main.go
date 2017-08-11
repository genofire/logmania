package console

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/bclicn/color"

	"github.com/genofire/logmania/bot"
	"github.com/genofire/logmania/lib"
	"github.com/genofire/logmania/log"
	"github.com/genofire/logmania/notify"
	configNotify "github.com/genofire/logmania/notify/config"
)

var (
	errOutput io.Writer = os.Stderr
	output    io.Writer = os.Stdout
)

// logger for output
type Notifier struct {
	notify.Notifier
	TimeFormat string
	ShowTime   bool
}

func Init(config *lib.NotifyConfig, state *configNotify.NotifyState, bot *bot.Bot) notify.Notifier {
	return &Notifier{
		TimeFormat: "2006-01-02 15:04:05",
		ShowTime:   true,
	}
}

// handle a log entry (print it on the terminal with color)
func (n *Notifier) Send(e *log.Entry) {
	if e.Hostname != "" {
		return
	}
	v := []interface{}{}
	format := "[%s]"

	if n.ShowTime {
		format = "%s [%s]"
		v = append(v, color.LightBlue(time.Now().Format(n.TimeFormat)))
	}
	if e.Hostname != "" {
		format = fmt.Sprintf("%s [%%s]", format)
		v = append(v, color.Purple(e.Hostname))
	}
	format = fmt.Sprintf("%s %%s", format)
	lvl := e.Level.String()
	switch e.Level {
	case log.DebugLevel:
		lvl = color.DarkGray(lvl)
	case log.InfoLevel:
		lvl = color.Green(lvl)
	case log.WarnLevel:
		lvl = color.Yellow(lvl)
	case log.ErrorLevel:
		lvl = color.Red(lvl)
	case log.PanicLevel:
		lvl = color.BRed(lvl)
	}

	v = append(v, lvl, e.Text)

	if len(e.Fields) > 0 {
		v = append(v, color.Purple(e.FieldString()))
		format = fmt.Sprintf("%s (%%s)\n", format)
	} else {
		format = fmt.Sprintf("%s\n", format)
	}

	text := fmt.Sprintf(format, v...)

	if e.Level > log.WarnLevel {
		errOutput.Write([]byte(text))
	} else {
		output.Write([]byte(text))
	}
}

func (n *Notifier) Close() {}

func init() {
	notify.AddNotifier(Init)
}
