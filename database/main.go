package database

import (
	"regexp"
	"time"

	"github.com/genofire/golang-lib/file"
	log "github.com/sirupsen/logrus"
)

const AlertMsg = "alert service from logmania, device did not send new message for a while"

type DB struct {
	Hostname       map[string]string                    `json:"hostname"`
	HostTo         map[string]map[string]bool           `json:"host_to"`
	MaxPrioIn      map[string]log.Level                 `json:"maxLevel"`
	RegexIn        map[string]map[string]*regexp.Regexp `json:"regexIn"`
	Lastseen       map[string]time.Time                 `json:"lastseen,omitempty"`
	LastseenNotify map[string]time.Time                 `json:"-"`
}

func (db *DB) SendTo(e *log.Entry) (*log.Entry, []string) {
	hostname, ok := e.Data["hostname"].(string)
	if !ok {
		return e, nil
	}
	if to, ok := db.HostTo[hostname]; ok {
		if e.Message != AlertMsg && hostname != "" {
			db.Lastseen[hostname] = time.Now()
		}
		var toList []string
		for toEntry, _ := range to {
			if lvl := db.MaxPrioIn[toEntry]; e.Level >= lvl {
				continue
			}
			if regex, ok := db.RegexIn[toEntry]; ok {
				stopForTo := false
				for _, expr := range regex {
					if expr.MatchString(e.Message) {
						stopForTo = true
						continue
					}
				}
				if stopForTo {
					continue
				}
			}
			toList = append(toList, toEntry)
		}
		if replaceHostname, ok := db.Hostname[hostname]; ok {
			entry := e.WithField("hostname", replaceHostname)
			entry.Level = e.Level
			entry.Message = e.Message
			return entry, toList
		}
		return e, toList
	} else {
		db.HostTo[hostname] = make(map[string]bool)
	}
	return e, nil
}

func (db *DB) AddRegex(to, expression string) error {
	regex, err := regexp.Compile(expression)
	if err == nil {
		if _, ok := db.RegexIn[to]; !ok {
			db.RegexIn[to] = make(map[string]*regexp.Regexp)
		}
		db.RegexIn[to][expression] = regex
		return nil
	}
	return err
}

func ReadDBFile(path string) *DB {
	var db DB

	if err := file.ReadJSON(path, &db); err == nil {
		log.Infof("loaded %d hosts", len(db.HostTo))
		if db.Lastseen == nil {
			db.Lastseen = make(map[string]time.Time)
		}
		if db.LastseenNotify == nil {
			db.LastseenNotify = make(map[string]time.Time)
		}
		if db.RegexIn == nil {
			db.RegexIn = make(map[string]map[string]*regexp.Regexp)
		} else {
			for to, regexs := range db.RegexIn {
				for exp, _ := range regexs {
					db.AddRegex(to, exp)
				}
			}
		}
		return &db
	} else {
		log.Error("failed to open db file: ", path, ":", err)
	}
	return &DB{
		Hostname:       make(map[string]string),
		HostTo:         make(map[string]map[string]bool),
		MaxPrioIn:      make(map[string]log.Level),
		RegexIn:        make(map[string]map[string]*regexp.Regexp),
		Lastseen:       make(map[string]time.Time),
		LastseenNotify: make(map[string]time.Time),
	}
}

func (db *DB) Alert(expired time.Duration, send func(e *log.Entry) error) {
	c := time.Tick(time.Minute)

	for range c {
		now := time.Now()
		for host, time := range db.Lastseen {
			if time.Before(now.Add(expired * -2)) {
				if timeNotify, ok := db.LastseenNotify[host]; !ok || !time.Before(timeNotify) {
					db.LastseenNotify[host] = now
					entry := log.NewEntry(log.New())
					entry.Level = log.ErrorLevel
					entry.Message = AlertMsg
					entry.WithField("hostname", host)
					send(entry)
				}
			}
		}
	}
}
