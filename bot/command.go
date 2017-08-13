package bot

import (
	"fmt"
	"strings"
	"time"

	timeago "github.com/ararog/timeago"

	"github.com/genofire/logmania/log"
)

type commandFunc func(func(string), string, []string)

// list help
func (b *Bot) help(answer func(string), from string, params []string) {
	msg := fmt.Sprintf("Hi %s there are the following commands:\n", from)
	for cmd := range b.commands {
		msg = fmt.Sprintf("%s - !%s\n", msg, cmd)
	}
	answer(msg)
}

// add a chat to send log to a chat
func (b *Bot) sendTo(answer func(string), from string, params []string) {
	if len(params) < 1 {
		answer("invalid: CMD IPAddress\n or\n CMD IPAddress to")
		return
	}
	host := params[0]
	to := from
	if len(params) > 1 {
		to = params[1]
	}

	if _, ok := b.state.HostTo[host]; !ok {
		b.state.HostTo[host] = make(map[string]bool)
	}
	b.state.HostTo[host][to] = true

	answer(fmt.Sprintf("added %s in list of %s", to, host))
}

//add a chat to send log to a chat
func (b *Bot) sendRemove(answer func(string), from string, params []string) {
	if len(params) < 1 {
		answer("invalid: CMD IPAddress\n or\n CMD IPAddress to")
		return
	}
	host := params[0]
	to := from
	if len(params) > 1 {
		to = params[1]
	}

	if list, ok := b.state.HostTo[host]; ok {
		delete(list, to)
		b.state.HostTo[host] = list
		answer(fmt.Sprintf("added %s in list of %s", to, host))
	} else {
		answer("not found host")
	}

}

// list all hostname with the chat where it send to
func (b *Bot) sendList(answer func(string), from string, params []string) {
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
	for ip, toMap := range b.state.HostTo {
		toList := ""
		show := all
		for to := range toMap {
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
		if hostname, ok := b.state.Hostname[ip]; ok {
			msg = fmt.Sprintf("%s%s (%s): %s\n", msg, ip, hostname, toList)
		} else {
			msg = fmt.Sprintf("%s%s: %s\n", msg, ip, toList)
		}
	}

	answer(msg)
}

// list all host with his ip
func (b *Bot) listHostname(answer func(string), from string, params []string) {
	msg := "hostnames:\n"
	for ip, hostname := range b.state.Hostname {
		if last, ok := b.state.Lastseen[ip]; ok {
			got, _ := timeago.TimeAgoWithTime(time.Now(), last)
			msg = fmt.Sprintf("%s%s - %s (%s)\n", msg, ip, hostname, got)
		} else {
			msg = fmt.Sprintf("%s%s - %s\n", msg, ip, hostname)
		}
	}
	answer(msg)
}

// list all hostname to a ip
func (b *Bot) setHostname(answer func(string), from string, params []string) {
	if len(params) < 2 {
		answer("invalid: CMD IPAddress NewHostname")
		return
	}
	host := params[0]
	name := params[1]

	b.state.Hostname[host] = name

	answer(fmt.Sprintf("set for %s the hostname %s", host, name))
}

// set a filter by max
func (b *Bot) listMaxfilter(answer func(string), from string, params []string) {
	msg := "filters: "
	if len(params) > 0 && params[0] == "all" {
		msg = fmt.Sprintf("%s\n", msg)
		for to, filter := range b.state.MaxPrioIn {
			msg = fmt.Sprintf("%s%s - %s\n", msg, to, filter.String())
		}
	} else {
		of := from
		if len(params) > 0 {
			of = params[0]
		}
		if filter, ok := b.state.MaxPrioIn[of]; ok {
			msg = fmt.Sprintf("%s of %s is %s", msg, of, filter)
		}
	}
	answer(msg)
}

// set a filter to a mix
func (b *Bot) setMaxfilter(answer func(string), from string, params []string) {
	if len(params) < 1 {
		answer("invalid: CMD Priority\n or\n CMD IPAddress Priority")
		return
	}
	to := from
	max := log.NewLoglevel(params[0])

	if len(params) > 1 {
		to = params[0]
		max = log.NewLoglevel(params[1])
	}

	b.state.MaxPrioIn[to] = max

	answer(fmt.Sprintf("set filter for %s to %s", to, max.String()))
}

// list of regex filter
func (b *Bot) listRegex(answer func(string), from string, params []string) {
	msg := "regexs:\n"
	if len(params) > 0 && params[0] == "all" {
		for to, regexs := range b.state.RegexIn {
			msg = fmt.Sprintf("%s%s\n-------------\n", msg, to)
			for expression := range regexs {
				msg = fmt.Sprintf("%s - %s\n", msg, expression)
			}
		}
	} else {
		of := from
		if len(params) > 0 {
			of = params[0]
		}
		if regexs, ok := b.state.RegexIn[of]; ok {
			msg = fmt.Sprintf("%s%s\n-------------\n", msg, from)
			for expression := range regexs {
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
	regex := strings.Join(params, " ")

	if err := b.state.AddRegex(from, regex); err == nil {
		answer(fmt.Sprintf("add regex for \"%s\" to %s", from, regex))
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
	regex := strings.Join(params, " ")
	delete(b.state.RegexIn[from], regex)
	b.listRegex(answer, from, []string{})
}
