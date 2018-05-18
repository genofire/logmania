package database

import (
	"regexp"
	"time"

	log "github.com/sirupsen/logrus"
)

const AlertMsg = "alert service from logmania, device did not send new message for a while"

type DB struct {
	// depraced Format -> transformation to new format by db.update()
	Hostname       map[string]string                    `json:"hostname,omitempty"`
	HostTo         map[string]map[string]bool           `json:"host_to,omitempty"`
	MaxPrioIn      map[string]log.Level                 `json:"maxLevel,omitempty"`
	RegexIn        map[string]map[string]*regexp.Regexp `json:"regexIn,omitempty"`
	Lastseen       map[string]time.Time                 `json:"lastseen,omitempty"`
	LastseenNotify map[string]time.Time                 `json:"-"`
	// new format
	Hosts             []*Host            `json:"hosts"`
	HostsByAddress    map[string]*Host   `json:"-"`
	HostsByName       map[string]*Host   `json:"-"`
	Notifies          []*Notify          `json:"notifies"`
	NotifiesByAddress map[string]*Notify `json:"-"`
}

func (db *DB) SendTo(e *log.Entry) (*log.Entry, *Host, []*Notify) {
	addr, ok := e.Data["hostname"].(string)
	if !ok {
		return e, nil, nil
	}
	var host *Host
	if host, ok := db.HostsByAddress[addr]; ok {
		if e.Message != AlertMsg {
			host.Lastseen = time.Now()
		}
		var toList []*Notify
		for _, notify := range host.NotifiesByAddress {
			if lvl := notify.MaxPrioIn; e.Level >= lvl {
				continue
			}
			stopForTo := false
			for _, expr := range notify.RegexIn {
				if expr.MatchString(e.Message) {
					stopForTo = true
					continue
				}
			}
			if stopForTo {
				continue
			}
			toList = append(toList, notify)
		}
		if host.Name != "" {
			entry := e.WithField("hostname", host.Name)
			entry.Level = e.Level
			entry.Message = e.Message
			return entry, host, toList
		}
		return e, host, toList
	} else {
		host = db.NewHost(addr)
	}
	return e, host, nil
}

func (db *DB) Alert(expired time.Duration, send func(e *log.Entry, n *Notify) bool) {
	c := time.Tick(time.Minute)

	for range c {
		now := time.Now()
		for _, h := range db.Hosts {
			if !h.Lastseen.Before(now.Add(expired * -2)) {
				continue
			}
			if h.LastseenNotify.Year() <= 1 && h.Lastseen.Before(h.LastseenNotify) {
				continue
			}
			h.LastseenNotify = now
			entry := log.NewEntry(log.New())
			entry.Level = log.ErrorLevel
			entry.Message = AlertMsg
			entry.WithField("hostname", h.Address)
			send(entry, nil)
		}
	}
}
