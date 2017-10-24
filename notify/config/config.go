package config

import (
	"encoding/json"
	"os"
	"regexp"
	"time"

	log "github.com/sirupsen/logrus"
)

const AlertMsg = "alert service from logmania, device did not send new message for a while"

type NotifyState struct {
	Hostname       map[string]string                    `json:"hostname"`
	HostTo         map[string]map[string]bool           `json:"host_to"`
	MaxPrioIn      map[string]log.Level                 `json:"maxLevel"`
	RegexIn        map[string]map[string]*regexp.Regexp `json:"regexIn"`
	Lastseen       map[string]time.Time                 `json:"lastseen,omitempty"`
	LastseenNotify map[string]time.Time                 `json:"-"`
}

func (state *NotifyState) SendTo(e *log.Entry) []string {
	hostname, ok := e.Data["hostname"].(string)
	if !ok {
		return nil
	}
	if to, ok := state.HostTo[hostname]; ok {
		if e.Message != AlertMsg && hostname != "" {
			state.Lastseen[hostname] = time.Now()
		}
		var toList []string
		for toEntry, _ := range to {
			if lvl := state.MaxPrioIn[toEntry]; e.Level >= lvl {
				continue
			}
			if regex, ok := state.RegexIn[toEntry]; ok {
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
		if replaceHostname, ok := state.Hostname[hostname]; ok {
			e.WithField("hostname", replaceHostname)
		}
		return toList
	} else {
		state.HostTo[hostname] = make(map[string]bool)
	}
	return nil
}

func (state *NotifyState) AddRegex(to, expression string) error {
	regex, err := regexp.Compile(expression)
	if err == nil {
		if _, ok := state.RegexIn[to]; !ok {
			state.RegexIn[to] = make(map[string]*regexp.Regexp)
		}
		state.RegexIn[to][expression] = regex
		return nil
	}
	return err
}

func ReadStateFile(path string) *NotifyState {
	var state NotifyState
	if f, err := os.Open(path); err == nil { // transform data to legacy meshviewer
		if err = json.NewDecoder(f).Decode(&state); err == nil {
			log.Infof("loaded %d hosts", len(state.HostTo))
			if state.Lastseen == nil {
				state.Lastseen = make(map[string]time.Time)
			}
			if state.LastseenNotify == nil {
				state.LastseenNotify = make(map[string]time.Time)
			}
			if state.RegexIn == nil {
				state.RegexIn = make(map[string]map[string]*regexp.Regexp)
			} else {
				for to, regexs := range state.RegexIn {
					for exp, _ := range regexs {
						state.AddRegex(to, exp)
					}
				}
			}
			return &state
		} else {
			log.Error("failed to unmarshal nodes:", err)
		}
	} else {
		log.Error("failed to open state notify file: ", path, ":", err)
	}
	return &NotifyState{
		Hostname:       make(map[string]string),
		HostTo:         make(map[string]map[string]bool),
		MaxPrioIn:      make(map[string]log.Level),
		RegexIn:        make(map[string]map[string]*regexp.Regexp),
		Lastseen:       make(map[string]time.Time),
		LastseenNotify: make(map[string]time.Time),
	}
}

func (state *NotifyState) Saver(path string) {
	c := time.Tick(time.Minute)

	for range c {
		state.SaveJSON(path)
	}
}

func (state *NotifyState) Alert(expired time.Duration, send func(e *log.Entry) error) {
	c := time.Tick(time.Minute)

	for range c {
		now := time.Now()
		for host, time := range state.Lastseen {
			if time.Before(now.Add(expired * -2)) {
				if timeNotify, ok := state.LastseenNotify[host]; !ok || !time.Before(timeNotify) {
					state.LastseenNotify[host] = now
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

// SaveJSON to path
func (state *NotifyState) SaveJSON(outputFile string) {
	tmpFile := outputFile + ".tmp"

	f, err := os.OpenFile(tmpFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Panic(err)
	}

	err = json.NewEncoder(f).Encode(state)
	if err != nil {
		log.Panic(err)
	}

	f.Close()
	if err := os.Rename(tmpFile, outputFile); err != nil {
		log.Panic(err)
	}
}
