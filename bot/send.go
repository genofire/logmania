package bot

import (
	"fmt"

	"dev.sum7.eu/genofire/logmania/database"
)

func NewSend(db *database.DB) *Command {
	return &Command{
		Name:        "send",
		Description: "list and configurate destination for hostnames",
		Commands: []*Command{
			{
				Name:        "add",
				Description: "add a destination for host with: IPAddress/Hostname [to]",
				Action: func(from string, params []string) string {
					if len(params) < 1 {
						return "invalid: IPAddress/Hostname [to]"
					}
					host := params[0]
					to := from
					if len(params) > 1 {
						to = params[1]
					}

					h := db.GetHost(host)
					if h == nil {
						h = db.NewHost(host)
					}
					if h == nil {
						return fmt.Sprintf("could not create host %s", host)
					}
					n, ok := db.NotifiesByAddress[to]
					if !ok {
						n = db.NewNotify(to)
					}
					if n == nil {
						return fmt.Sprintf("could not create notify %s in list of %s", to, host)
					}
					h.AddNotify(n)

					return fmt.Sprintf("added %s in list of %s", to, host)
				},
			},
			{
				Name:        "del",
				Description: "del a destination for host with: IPAddress/Hostname [to]",
				Action: func(from string, params []string) string {
					if len(params) < 1 {
						return "invalid: IPAddress/Hostname [to]"
					}
					host := params[0]
					to := from
					if len(params) > 1 {
						to = params[1]
					}

					if h := db.GetHost(host); h != nil {
						h.DeleteNotify(to)
						return fmt.Sprintf("removed %s in list of %s", to, host)
					}
					return "not found host"
				},
			},
			{
				Name:        "all",
				Description: "list of all hosts with there channels",
				Action: func(from string, params []string) string {
					msg := "sending:\n"
					for _, host := range db.Hosts {
						toList := ""
						for _, to := range host.Notifies {
							toList = fmt.Sprintf("%s , %s", toList, to)
						}
						if len(toList) > 3 {
							toList = toList[3:]
						}
						if host.Name != "" {
							msg = fmt.Sprintf("%s%s (%s): %s\n", msg, host.Address, host.Name, toList)
						} else {
							msg = fmt.Sprintf("%s%s: %s\n", msg, host.Address, toList)
						}
					}
					return msg
				},
			},
			{
				Name:        "channel",
				Description: "list all host of given channel: channel",
				Action: func(from string, params []string) string {
					if len(params) != 1 {
						return "invalid: no channel given"
					}

					of := params[0]
					msg := "sending:\n"

					for _, host := range db.Hosts {
						show := false
						for _, to := range host.Notifies {
							if to == of {
								show = true
								break
							}
						}
						if !show {
							continue
						}
						if host.Name != "" {
							msg = fmt.Sprintf("%s%s (%s)\n", msg, host.Address, host.Name)
						} else {
							msg = fmt.Sprintf("%s%s\n", msg, host.Address)
						}
					}
					return msg
				},
			},
		},
		Action: func(from string, params []string) string {
			msg := "sending:\n"
			for _, host := range db.Hosts {
				show := false
				for _, to := range host.Notifies {
					if to == from {
						show = true
						break
					}
				}
				if !show {
					continue
				}
				if host.Name != "" {
					msg = fmt.Sprintf("%s%s (%s)\n", msg, host.Address, host.Name)
				} else {
					msg = fmt.Sprintf("%s%s\n", msg, host.Address)
				}
			}
			return msg
		},
	}
}
