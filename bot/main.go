package bot

import (
	"fmt"
	"strings"

	"github.com/genofire/logmania/database"
)

type Bot struct {
	db       *database.DB
	commands map[string]commandFunc
}

func NewBot(db *database.DB) *Bot {
	b := &Bot{
		db: db,
	}
	b.commands = map[string]commandFunc{
		"help":          b.help,
		"send-to":       b.sendTo,
		"send-list":     b.sendList,
		"send-rm":       b.sendRemove,
		"hostname-set":  b.setHostname,
		"hostname-list": b.listHostname,
		"filter-set":    b.setMaxfilter,
		"filter-list":   b.listMaxfilter,
		"regex-add":     b.addRegex,
		"regex-list":    b.listRegex,
		"regex-rm":      b.delRegex,
	}
	return b
}

func (b *Bot) Handle(answer func(string), from, msg string) {
	msgParts := strings.Split(msg, " ")
	if len(msgParts[0]) <= 0 || msgParts[0][0] != '!' {
		return
	}
	cmdName := msgParts[0][1:]
	if cmd, ok := b.commands[cmdName]; ok {
		cmd(answer, from, msgParts[1:])
	} else {
		answer(fmt.Sprintf("not found command: !%s", cmdName))
	}
}
