package bot

import (
	"fmt"

	"dev.sum7.eu/genofire/logmania/database"
)

func NewFilter(db *database.DB) *Command {
	return &Command{
		Name:        "filter",
		Description: "list and configurate regex filter for channel by message content",
		Commands: []*Command{
			{
				Name:        "add",
				Description: "add regex filter for channel:  [channel] regex",
				Action: func(from string, params []string) string {
					if len(params) < 1 {
						return "invalid: [channel] regex"
					}
					of := from
					regex := params[0]
					if len(params) > 1 {
						of = params[0]
						regex = params[1]
					}

					n := db.NotifiesByAddress[of]
					if err := n.AddRegex(regex); err != nil {
						return fmt.Sprintf("\"%s\" is no valid regex expression: %s", regex, err)
					}
					return fmt.Sprintf("add regex for \"%s\" to %s", of, regex)
				},
			},
			{
				Name:        "del",
				Description: "del regex filter for channel:  [channel] regex",
				Action: func(from string, params []string) string {
					if len(params) < 1 {
						return "invalid: [channel] regex"
					}
					of := from
					regex := params[0]
					if len(params) > 1 {
						of = params[0]
						regex = params[1]
					}
					n := db.NotifiesByAddress[of]
					delete(n.RegexIn, regex)
					return "deleted"
				},
			},
			{
				Name:        "all",
				Description: "list of all channels",
				Action: func(from string, params []string) string {
					msg := "filter:\n"
					for _, n := range db.Notifies {
						msg = fmt.Sprintf("%s%s\n-------------\n", msg, n.Address())
						for expression := range n.RegexIn {
							msg = fmt.Sprintf("%s - %s\n", msg, expression)
						}
					}
					return msg
				},
			},
			{
				Name:        "channel",
				Description: "list of given channel: channel",
				Action: func(from string, params []string) string {
					msg := "filter:\n"
					if len(params) != 1 {
						return "invalid: no channel given"
					}
					of := params[0]
					if n, ok := db.NotifiesByAddress[of]; ok {
						msg = fmt.Sprintf("%s%s\n-------------\n", msg, of)
						for expression := range n.RegexIn {
							msg = fmt.Sprintf("%s - %s\n", msg, expression)
						}
					}
					return msg
				},
			},
		},
		Action: func(from string, params []string) string {
			msg := "filter:\n"
			if n, ok := db.NotifiesByAddress[from]; ok {
				msg = fmt.Sprintf("%s%s\n-------------\n", msg, from)
				for expression := range n.RegexIn {
					msg = fmt.Sprintf("%s - %s\n", msg, expression)
				}
			}
			return msg
		},
	}
}
