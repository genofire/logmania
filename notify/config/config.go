package config

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/genofire/logmania/log"
)

type NotifyState struct {
	Hostname  map[string]string                    `json:"hostname"`
	HostTo    map[string]map[string]bool           `json:"host_to"`
	MaxPrioIn map[string]log.LogLevel              `json:"maxLevel"`
	RegexIn   map[string]map[string]bool           `json:"regexIn"`
	regexIn   map[string]map[string]*regexp.Regexp `json:"-"`
}

func (state *NotifyState) SendTo(e *log.Entry) []string {
	if to, ok := state.HostTo[e.Hostname]; ok {
		var toList []string
		for toEntry, _ := range to {
			if lvl := state.MaxPrioIn[toEntry]; e.Level < lvl {
				continue
			}
			toList = append(toList, toEntry)
		}
		if hostname, ok := state.Hostname[e.Hostname]; ok {
			e.Hostname = hostname
		}
		return toList
	} else {
		state.HostTo[e.Hostname] = make(map[string]bool)
	}
	return nil
}

func ReadStateFile(path string) *NotifyState {
	var state NotifyState
	if f, err := os.Open(path); err == nil { // transform data to legacy meshviewer
		if err = json.NewDecoder(f).Decode(&state); err == nil {
			fmt.Println("loaded", len(state.HostTo), "hosts")
			state.regexIn = make(map[string]map[string]*regexp.Regexp)
			return &state
		} else {
			fmt.Println("failed to unmarshal nodes:", err)
		}
	} else {
		fmt.Println("failed to open state notify file: ", path, ":", err)
	}
	return &NotifyState{
		Hostname:  make(map[string]string),
		HostTo:    make(map[string]map[string]bool),
		MaxPrioIn: make(map[string]log.LogLevel),
		RegexIn:   make(map[string]map[string]bool),
		regexIn:   make(map[string]map[string]*regexp.Regexp),
	}
}

func (state *NotifyState) Saver(path string) {
	c := time.Tick(time.Minute)

	for range c {
		state.SaveJSON(path)
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
