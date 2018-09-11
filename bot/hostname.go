package bot

import (
	"fmt"

	"dev.sum7.eu/genofire/logmania/database"
	timeago "github.com/ararog/timeago"
)

func NewHostname(db *database.DB) *Command {
	return &Command{
		Name:        "hostname",
		Description: "alternative short (host)names for long IP-Addresses or URLs (and time of last recieved input)",
		Commands: []*Command{
			&Command{
				Name:        "set",
				Description: "set or replace a hostname: IPAddress/Hostname NewHostname",
				Action: func(from string, params []string) string {
					if len(params) != 2 {
						return "invalid: IPAddress/Hostname NewHostname"
					}
					addr := params[0]
					name := params[1]

					h := db.GetHost(addr)
					if h == nil {
						h = db.NewHost(addr)
					}
					db.ChangeHostname(h, name)

					return fmt.Sprintf("set for %s the hostname %s", addr, name)
				},
			},
			&Command{
				Name:        "del",
				Description: "delete a hostname: IPAddress/Hostname",
				Action: func(from string, params []string) string {
					if len(params) != 1 {
						return "invalid: IPAddress/Hostname"
					}
					addr := params[0]
					h := db.GetHost(addr)
					if h != nil {
						db.DeleteHost(h)
						if h.Name != "" {
							return fmt.Sprintf("remove host %s with hostname %s", h.Address, h.Name)
						}
						return fmt.Sprintf("remove host %s", h.Address)
					}
					return "could not found host"
				},
			},
		},
		Action: func(from string, params []string) string {
			msg := "hostnames:\n"
			for _, host := range db.Hosts {
				if host.Lastseen.Year() > 1 {
					got, _ := timeago.TimeAgoFromNowWithTime(host.Lastseen)
					msg = fmt.Sprintf("%s%s - %s (%s)\n", msg, host.Address, host.Name, got)
				} else {
					msg = fmt.Sprintf("%s%s - %s\n", msg, host.Address, host.Name)
				}
			}
			return msg
		},
	}
}
