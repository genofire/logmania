package bot

import (
	"fmt"

	timeago "github.com/ararog/timeago"
	log "github.com/sirupsen/logrus"
)

type commandFunc func(func(string), string, []string)

// list help
func (b *Bot) help(answer func(string), from string, params []string) {
	msg := fmt.Sprintf("Hi %s there are the following commands:\n", from)
	for _, cmd := range b.commands {
		msg = fmt.Sprintf("%s - !%s\n", msg, cmd)
	}
	answer(msg)
}

// add a chat to send log to a chat
func (b *Bot) addSend(answer func(string), from string, params []string) {
	if len(params) < 1 {
		answer("invalid: CMD IPAddress/Hostname\n or\n CMD IPAddress/Hostname to")
		return
	}
	host := params[0]
	to := from
	if len(params) > 1 {
		to = params[1]
	}

	h := b.db.GetHost(host)
	if h == nil {
		h = b.db.NewHost(host)
	}
	n, ok := b.db.NotifiesByAddress[to]
	if !ok {
		n = b.db.NewNotify(to)
	}
	h.AddNotify(n)

	answer(fmt.Sprintf("added %s in list of %s", to, host))
}

//add a chat to send log to a chat
func (b *Bot) delSend(answer func(string), from string, params []string) {
	if len(params) < 1 {
		answer("invalid: CMD IPAddress/Hostname\n or\n CMD IPAddress/Hostname to")
		return
	}
	host := params[0]
	to := from
	if len(params) > 1 {
		to = params[1]
	}

	if h := b.db.GetHost(host); h != nil {
		h.DeleteNotify(to)
		answer(fmt.Sprintf("removed %s in list of %s", to, host))
	} else {
		answer("not found host")
	}

}

// list all hostname with the chat where it send to
func (b *Bot) listSend(answer func(string), from string, params []string) {
	msg := "sending:\n"
	all := false
	of := from
	if len(params) > 0 {
		if params[0] == "all" {
			all = true
		} else {
			of = params[0]
		}
	}
	for _, host := range b.db.Hosts {
		toList := ""
		show := all
		for _, to := range host.Notifies {
			if all {
				toList = fmt.Sprintf("%s , %s", toList, to)
			} else if to == of {
				show = true
			}
		}
		if !show {
			continue
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

	answer(msg)
}

// add hostname
func (b *Bot) addHostname(answer func(string), from string, params []string) {
	if len(params) < 2 {
		answer("invalid: CMD IPAddress/Hostname NewHostname")
		return
	}
	addr := params[0]
	name := params[1]

	h := b.db.GetHost(addr)
	if h == nil {
		h = b.db.NewHost(addr)
	}
	b.db.ChangeHostname(h, name)

	answer(fmt.Sprintf("set for %s the hostname %s", addr, name))
}

func (b *Bot) delHostname(answer func(string), from string, params []string) {
	if len(params) < 2 {
		answer("invalid: CMD IPAddress/Hostname")
		return
	}
	addr := params[0]
	h := b.db.GetHost(addr)
	if h != nil {
		b.db.DeleteHost(h)
		if h.Name != "" {
			answer(fmt.Sprintf("remove host %s with hostname %s", h.Address, h.Name))
		} else {
			answer(fmt.Sprintf("remove host %s", h.Address))
		}
	} else {
		answer("could not found host")
	}
}

// list all host with his ip
func (b *Bot) listHostname(answer func(string), from string, params []string) {
	msg := "hostnames:\n"
	for _, host := range b.db.Hosts {
		if host.Lastseen.Year() > 1 {
			got, _ := timeago.TimeAgoFromNowWithTime(host.Lastseen)
			msg = fmt.Sprintf("%s%s - %s (%s)\n", msg, host.Address, host.Name, got)
		} else {
			msg = fmt.Sprintf("%s%s - %s\n", msg, host.Address, host.Name)
		}
	}
	answer(msg)
}

// set a filter by max
func (b *Bot) listMaxfilter(answer func(string), from string, params []string) {
	msg := "filters: "
	if len(params) > 0 && params[0] == "all" {
		msg = fmt.Sprintf("%s\n", msg)
		for _, n := range b.db.Notifies {
			msg = fmt.Sprintf("%s%s - %s\n", msg, n.Address(), n.MaxPrioIn.String())
		}
	} else {
		of := from
		if len(params) > 0 {
			of = params[0]
		}
		if filter, ok := b.db.NotifiesByAddress[of]; ok {
			msg = fmt.Sprintf("%s of %s is %s", msg, of, filter)
		}
	}
	answer(msg)
}

// set a filter to a mix
func (b *Bot) setMaxfilter(answer func(string), from string, params []string) {
	if len(params) < 1 {
		answer("invalid: CMD Priority\n or\n CMD Channel Priority")
		return
	}
	to := from
	var max log.Level
	var err error

	if len(params) > 1 {
		to = params[0]
		max, err = log.ParseLevel(params[1])
	} else {
		max, err = log.ParseLevel(params[0])
	}
	if err != nil {
		answer("invalid priority: CMD Priority\n or\n CMD Channel Priority")
		return
	}
	n, ok := b.db.NotifiesByAddress[to]
	if !ok {
		n = b.db.NewNotify(to)
	}

	n.MaxPrioIn = max

	answer(fmt.Sprintf("set filter for %s to %s", to, max.String()))
}

// list of regex filter
func (b *Bot) listRegex(answer func(string), from string, params []string) {
	msg := "regexs:\n"
	if len(params) > 0 && params[0] == "all" {
		for _, n := range b.db.Notifies {
			msg = fmt.Sprintf("%s%s\n-------------\n", msg, n.Address())
			for expression := range n.RegexIn {
				msg = fmt.Sprintf("%s - %s\n", msg, expression)
			}
		}
	} else {
		of := from
		if len(params) > 0 {
			of = params[0]
		}
		if n, ok := b.db.NotifiesByAddress[of]; ok {
			msg = fmt.Sprintf("%s%s\n-------------\n", msg, of)
			for expression := range n.RegexIn {
				msg = fmt.Sprintf("%s - %s\n", msg, expression)
			}
		}
	}
	answer(msg)
}

// add a regex filter
func (b *Bot) addRegex(answer func(string), from string, params []string) {
	if len(params) < 1 {
		answer("invalid: CMD regex\n or\n CMD channel regex")
		return
	}
	of := from
	regex := params[0]
	if len(params) > 1 {
		of = params[0]
		regex = params[1]
	}

	n := b.db.NotifiesByAddress[of]
	if err := n.AddRegex(regex); err == nil {
		answer(fmt.Sprintf("add regex for \"%s\" to %s", of, regex))
	} else {
		answer(fmt.Sprintf("\"%s\" is no valid regex expression: %s", regex, err))
	}
}

// del a regex filter
func (b *Bot) delRegex(answer func(string), from string, params []string) {
	if len(params) < 1 {
		answer("invalid: CMD regex\n or\n CMD channel regex")
		return
	}
	of := from
	regex := params[0]
	if len(params) > 1 {
		of = params[0]
		regex = params[1]
	}
	n := b.db.NotifiesByAddress[of]
	delete(n.RegexIn, regex)
	b.listRegex(answer, of, []string{})
}

// list of regex replace
func (b *Bot) listRegexReplace(answer func(string), from string, params []string) {
	msg := "replaces:\n"
	if len(params) > 0 && params[0] == "all" {
		for _, n := range b.db.Notifies {
			msg = fmt.Sprintf("%s%s\n-------------\n", msg, n.Address())
			for expression, value := range n.RegexReplace {
				msg = fmt.Sprintf("%s - \"%s\" : \"%s\"\n", msg, expression, value)
			}
		}
	} else {
		of := from
		if len(params) > 0 {
			of = params[0]
		}
		if n, ok := b.db.NotifiesByAddress[of]; ok {
			msg = fmt.Sprintf("%s%s\n-------------\n", msg, of)
			for expression, value := range n.RegexReplace {
				msg = fmt.Sprintf("%s - \"%s\" : \"%s\"\n", msg, expression, value)
			}
		}
	}
	answer(msg)
}

// add a regex replace
func (b *Bot) addRegexReplace(answer func(string), from string, params []string) {
	if len(params) < 1 {
		answer("invalid: CMD regex replace\n or\n CMD channel regex replace")
		return
	}
	of := from
	regex := params[0]
	value := params[1]
	if len(params) > 2 {
		of = params[0]
		regex = params[1]
		value = params[2]
	}

	n := b.db.NotifiesByAddress[of]
	if err := n.AddRegexReplace(regex, value); err == nil {
		answer(fmt.Sprintf("add replace in \"%s\" for \"%s\" to \"%s\"", of, regex, value))
	} else {
		answer(fmt.Sprintf("\"%s\" to \"%s\" is no valid regex replace expression: %s", regex, value, err))
	}
}

// del a regex replace
func (b *Bot) delRegexReplace(answer func(string), from string, params []string) {
	if len(params) < 1 {
		answer("invalid: CMD regex\n or\n CMD channel regex")
		return
	}
	of := from
	regex := params[0]
	if len(params) > 1 {
		of = params[0]
		regex = params[1]
	}
	n := b.db.NotifiesByAddress[of]

	delete(n.RegexReplace, regex)
	b.listRegexReplace(answer, of, []string{})
}
