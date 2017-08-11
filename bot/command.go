package bot

import (
	"fmt"

	"github.com/genofire/logmania/log"
)

type commandFunc func(func(string), string, []string)

func (b *Bot) help(answer func(string), from string, params []string) {
	msg := fmt.Sprintf("Hi %s there are the following commands:\n", from)
	for cmd := range b.commands {
		msg = fmt.Sprintf("%s - !%s\n", msg, cmd)
	}
	answer(msg)
}

func (b *Bot) sendTo(answer func(string), from string, params []string) {
	host := params[0]
	to := from
	if len(params) > 1 {
		to = params[1]
	}

	if list, ok := b.state.HostTo[host]; ok {
		b.state.HostTo[host] = append(list, to)
	} else {
		b.state.HostTo[host] = []string{to}
	}

	answer(fmt.Sprintf("added %s in list of %s", to, from))
}

func (b *Bot) setHostname(answer func(string), from string, params []string) {
	host := params[0]
	name := params[1]

	b.state.Hostname[host] = name

	answer(fmt.Sprintf("set for %s the hostname %s", host, name))
}

func (b *Bot) listHostname(answer func(string), from string, params []string) {
	msg := "hostnames:\n"
	for ip, hostname := range b.state.Hostname {
		msg = fmt.Sprintf("%s%s - %s", msg, ip, hostname)
	}
	answer(msg)
}

func (b *Bot) listMaxfilter(answer func(string), from string, params []string) {
	msg := "filters:\n"
	for to, filter := range b.state.MaxPrioIn {
		msg = fmt.Sprintf("%s%s - %s", msg, to, filter.String())
	}
	answer(msg)
}

func (b *Bot) setMaxfilter(answer func(string), from string, params []string) {
	to := from
	max := log.NewLoglevel(params[0])

	if len(params) > 1 {
		to = params[0]
		max = log.NewLoglevel(params[1])
	}

	b.state.MaxPrioIn[to] = max

	answer(fmt.Sprintf("set filter for %s to %s", to, max.String()))
}
