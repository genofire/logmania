package bot

import (
	"fmt"

	"github.com/bdlm/log"
	logstd "github.com/bdlm/std/logger"

	"dev.sum7.eu/sum7/logmania/database"
)

func NewPriority(db *database.DB) *Command {
	return &Command{
		Name:        "priority",
		Description: "list and configure priority in channel",
		Commands: []*Command{
			{
				Name:        "set",
				Description: "set max priority of channel: [channel] Priority",
				Action: func(from string, params []string) string {
					if len(params) < 1 {
						return "invalid: [channel] Priority"
					}
					to := from
					var max logstd.Level
					var err error

					if len(params) > 1 {
						to = params[0]
						max, err = log.ParseLevel(params[1])
					} else {
						max, err = log.ParseLevel(params[0])
					}
					if err != nil {
						return "invalid: [Channel] Priority"
					}
					n, ok := db.NotifiesByAddress[to]
					if !ok {
						n = db.NewNotify(to)
					}

					n.MaxPrioIn = max

					return fmt.Sprintf("set filter for %s to %s", to, log.LevelString(max))
				},
			},
			{
				Name:        "all",
				Description: "list of all channels",
				Action: func(from string, params []string) string {
					msg := "priority: \n"
					for _, n := range db.Notifies {
						msg = fmt.Sprintf("%s%s - %s\n", msg, n.Address(), log.LevelString(n.MaxPrioIn))
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
					msg := "priority: \n"
					if notify, ok := db.NotifiesByAddress[of]; ok {
						msg = fmt.Sprintf("%s %s is %s", msg, of, log.LevelString(notify.MaxPrioIn))
					}
					return msg
				},
			},
		},
		Action: func(from string, params []string) string {
			msg := "priority: \n"
			if notify, ok := db.NotifiesByAddress[from]; ok {
				msg = fmt.Sprintf("%s %s is %s", msg, from, log.LevelString(notify.MaxPrioIn))
			}
			return msg
		},
	}
}
