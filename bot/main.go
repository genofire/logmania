package bot

import (
	"github.com/mattn/go-shellwords"

	"dev.sum7.eu/genofire/logmania/database"
)

type Bot struct {
	Command
}

func NewBot(db *database.DB) *Bot {
	return &Bot{Command{
		Description: "logmania bot, to configurate live all settings",
		Commands: []*Command{
			NewFilter(db),
			NewHostname(db),
			NewPriority(db),
			NewReplace(db),
			NewSend(db),
		},
	}}
}

func (b *Bot) Handle(from, msg string) string {
	msgParts, err := shellwords.Parse(msg)
	if err != nil {
		return ""
	}
	if len(msgParts) <= 0 || msgParts[0][0] != '.' {
		return ""
	}
	msgParts[0] = msgParts[0][1:]
	return b.Run(from, msgParts)
}
