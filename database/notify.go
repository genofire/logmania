package database

import (
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

type Notify struct {
	Protocol               string                    `json:"proto"`
	To                     string                    `json:"to"`
	RegexIn                map[string]*regexp.Regexp `json:"regexIn"`
	RegexReplace           map[string]string         `json:"regexReplace"`
	MaxPrioIn              log.Level                 `json:"maxLevel"`
	regexReplaceExpression map[string]*regexp.Regexp
}

func (n *Notify) Init() {
	if n.RegexIn == nil {
		n.RegexIn = make(map[string]*regexp.Regexp)
	}
	if n.RegexReplace == nil {
		n.RegexReplace = make(map[string]string)
	}
	if n.regexReplaceExpression == nil {
		n.regexReplaceExpression = make(map[string]*regexp.Regexp)
	}
	for exp := range n.RegexIn {
		regex, err := regexp.Compile(exp)
		if err == nil {
			n.RegexIn[exp] = regex
		}
	}
	for exp := range n.RegexReplace {
		regex, err := regexp.Compile(exp)
		if err == nil {
			n.regexReplaceExpression[exp] = regex
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
func (n *Notify) AddRegexReplace(expression, value string) error {
	regex, err := regexp.Compile(expression)
	if err == nil {
		n.regexReplaceExpression[expression] = regex
		n.RegexReplace[expression] = value
	}
	return err
}
func (n *Notify) RunReplace(msg string) string {
	for key, re := range n.regexReplaceExpression {
		value := n.RegexReplace[key]
		msg = re.ReplaceAllString(msg, value)
	}
	return msg
}

func (n *Notify) Address() string {
	return n.Protocol + ":" + n.To
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
	if len(addr) != 2 {
		return nil
	}
	n := &Notify{
		Protocol:  addr[0],
		To:        addr[1],
		RegexIn:   make(map[string]*regexp.Regexp),
		MaxPrioIn: log.DebugLevel,
	}
	db.AddNotify(n)
	return n
}
