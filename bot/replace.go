package bot

import (
	"fmt"

	"dev.sum7.eu/genofire/logmania/database"
)

func NewReplace(db *database.DB) *Command {
	return &Command{
		Name:        "replace",
		Description: "list and configurate replace content of message for channel",
		Commands: []*Command{
			{
				Name:        "add",
				Description: "add regex replace for channel:  [channel] regex replace",
				Action: func(from string, params []string) string {
					if len(params) < 1 {
						return "invalid: [channel] regex replace"
					}
					of := from
					regex := params[0]
					value := params[1]
					if len(params) > 2 {
						of = params[0]
						regex = params[1]
						value = params[2]
					}

					n := db.NotifiesByAddress[of]
					if err := n.AddRegexReplace(regex, value); err != nil {
						return fmt.Sprintf("\"%s\" to \"%s\" is no valid regex replace expression: %s", regex, value, err)
					}
					return fmt.Sprintf("add replace in \"%s\" for \"%s\" to \"%s\"", of, regex, value)
				},
			},
			{
				Name:        "del",
				Description: "del regex replace for channel:  [channel] regex replace",
				Action: func(from string, params []string) string {
					if len(params) < 1 {
						return "invalid: [channel] regex replace"
					}
					of := from
					regex := params[0]
					if len(params) > 1 {
						of = params[0]
						regex = params[1]
					}
					n := db.NotifiesByAddress[of]

					delete(n.RegexReplace, regex)
					return "deleted"
				},
			},

			{
				Name:        "all",
				Description: "list of all channels",
				Action: func(from string, params []string) string {
					msg := "replaces:\n"
					for _, n := range db.Notifies {
						msg = fmt.Sprintf("%s%s\n-------------\n", msg, n.Address())
						for expression, value := range n.RegexReplace {
							msg = fmt.Sprintf("%s - \"%s\" : \"%s\"\n", msg, expression, value)
						}
					}
					return msg
				},
			},
			{
				Name:        "channel",
				Description: "list of given channel: channel",
				Action: func(from string, params []string) string {
					if len(params) != 1 {
						return "invalid: no channel given"
					}
					of := params[0]
					msg := "replaces:\n"
					if n, ok := db.NotifiesByAddress[of]; ok {
						msg = fmt.Sprintf("%s%s\n-------------\n", msg, of)
						for expression, value := range n.RegexReplace {
							msg = fmt.Sprintf("%s - \"%s\" : \"%s\"\n", msg, expression, value)
						}
					}
					return msg
				},
			},
		},
		Action: func(from string, params []string) string {
			msg := "replaces:\n"
			if n, ok := db.NotifiesByAddress[from]; ok {
				msg = fmt.Sprintf("%s%s\n-------------\n", msg, from)
				for expression, value := range n.RegexReplace {
					msg = fmt.Sprintf("%s - \"%s\" : \"%s\"\n", msg, expression, value)
				}
			}
			return msg
		},
	}
}
