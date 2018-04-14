package database

import (
	"time"
)

type Host struct {
	Name              string             `json:"name"`
	Address           string             `json:"address"`
	Lastseen          time.Time          `json:"lastseen,omitempty"`
	LastseenNotify    time.Time          `json:"lastseen_notify,omitempty"`
	Notifies          []string           `json:"notifies"`
	NotifiesByAddress map[string]*Notify `json:"-"`
}

func (h *Host) AddNotify(n *Notify) {
	if _, ok := h.NotifiesByAddress[n.Address()]; !ok {
		h.Notifies = append(h.Notifies, n.Address())
		h.NotifiesByAddress[n.Address()] = n
	}
}
func (h *Host) DeleteNotify(to string) {
	delete(h.NotifiesByAddress, to)
	for i, v := range h.Notifies {
		if v == to {
			copy(h.Notifies[i:], h.Notifies[i+1:])
			h.Notifies = h.Notifies[:len(h.Notifies)-1]
			return
		}
	}
	return
}

// -- global notify

func (db *DB) InitHost() {
	if db.HostsByAddress == nil {
		db.HostsByAddress = make(map[string]*Host)
	}
	if db.HostsByName == nil {
		db.HostsByName = make(map[string]*Host)
	}
	for _, h := range db.Hosts {
		if h.NotifiesByAddress == nil {
			h.NotifiesByAddress = make(map[string]*Notify)
		}
		for _, nName := range h.Notifies {
			h.NotifiesByAddress[nName] = db.NotifiesByAddress[nName]
		}
		db.HostsByAddress[h.Address] = h
		db.HostsByName[h.Name] = h
	}
}

func (db *DB) AddHost(h *Host) {
	db.Hosts = append(db.Hosts, h)
	db.HostsByAddress[h.Address] = h
	db.HostsByName[h.Name] = h
}

func (db *DB) GetHost(str string) *Host {
	h, ok := db.HostsByAddress[str]
	if ok {
		return h
	}
	return db.HostsByName[str]
}
func (db *DB) DeleteHost(h *Host) {
	delete(db.HostsByAddress, h.Address)
	delete(db.HostsByName, h.Name)
	for i, v := range db.Hosts {
		if v.Address == h.Address {
			copy(db.Hosts[i:], db.Hosts[i+1:])
			db.Hosts[len(db.Hosts)-1] = nil
			db.Hosts = db.Hosts[:len(db.Hosts)-1]
			return
		}
	}
	return
}

func (db *DB) ChangeHostname(h *Host, name string) {
	delete(db.HostsByName, h.Name)
	h.Name = name
	db.HostsByName[name] = h
}

func (db *DB) NewHost(addr string) *Host {
	h := &Host{
		Address:           addr,
		NotifiesByAddress: make(map[string]*Notify),
	}
	db.AddHost(h)
	return h
}
