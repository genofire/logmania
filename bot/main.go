package bot

import (
	"fmt"
	"strings"

	"dev.sum7.eu/genofire/logmania/database"
)

type Bot struct {
	db          *database.DB
	commandsMap map[string]commandFunc
	commands    []string
}

func NewBot(db *database.DB) *Bot {
	b := &Bot{
		db: db,
	}
	b.commandsMap = map[string]commandFunc{
		"help":          b.help,
		"send-add":      b.addSend,
		"send-list":     b.listSend,
		"send-del":      b.delSend,
		"hostname-set":  b.addHostname,
		"hostname-list": b.listHostname,
		"hostname-del":  b.delHostname,
		"filter-set":    b.setMaxfilter,
		"filter-list":   b.listMaxfilter,
		"regex-add":     b.addRegex,
		"regex-list":    b.listRegex,
		"regex-del":     b.delRegex,
	}
	for k, _ := range b.commandsMap {
		b.commands = append(b.commands, k)
	}
	return b
}

func (b *Bot) Handle(answer func(string), from, msg string) {
	msgParts := strings.Split(msg, " ")
	if len(msgParts[0]) <= 0 || msgParts[0][0] != '!' {
		return
	}
	cmdName := msgParts[0][1:]
	if cmd, ok := b.commandsMap[cmdName]; ok {
		cmd(answer, from, msgParts[1:])
	} else {
		answer(fmt.Sprintf("not found command: !%s", cmdName))
	}
}
