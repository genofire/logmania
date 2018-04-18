package database

import (
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

type Notify struct {
	Protocoll string                    `json:"proto"`
	To        string                    `json:"to"`
	RegexIn   map[string]*regexp.Regexp `json:"regexIn"`
	MaxPrioIn log.Level                 `json:"maxLevel"`
}

func (n *Notify) Init() {
	if n.RegexIn == nil {
		n.RegexIn = make(map[string]*regexp.Regexp)
	}
	for exp := range n.RegexIn {
		regex, err := regexp.Compile(exp)
		if err == nil {
			n.RegexIn[exp] = regex
		}
	}
}

func (n *Notify) AddRegex(expression string) error {
	regex, err := regexp.Compile(expression)
	if err == nil {
		n.RegexIn[expression] = regex
	}
	return err
}

func (n *Notify) Address() string {
	return n.Protocoll + ":" + n.To
}

// -- global notify

func (db *DB) InitNotify() {
	if db.NotifiesByAddress == nil {
		db.NotifiesByAddress = make(map[string]*Notify)
	}
	for _, n := range db.Notifies {
		n.Init()
		db.NotifiesByAddress[n.Address()] = n
	}
}

func (db *DB) AddNotify(n *Notify) {
	db.Notifies = append(db.Notifies, n)
	db.NotifiesByAddress[n.Address()] = n
}

func (db *DB) NewNotify(to string) *Notify {
	addr := strings.Split(to, ":")
	n := &Notify{
		Protocoll: addr[0],
		To:        addr[1],
		RegexIn:   make(map[string]*regexp.Regexp),
	}
	db.AddNotify(n)
	return n
}
