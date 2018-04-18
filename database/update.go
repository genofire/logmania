package database

func (db *DB) update() {
	if db.Hostname != nil {
		for addr, name := range db.Hostname {
			h, ok := db.HostsByName[name]
			if !ok {
				h = db.NewHost(addr)
				db.ChangeHostname(h, name)
			}
		}
		db.Hostname = nil
	}
	if db.HostTo != nil {
		for addr, channels := range db.HostTo {
			h, ok := db.HostsByAddress[addr]
			if !ok {
				h = db.NewHost(addr)
			}
			for channel := range channels {
				n, ok := db.NotifiesByAddress[channel]
				if !ok {
					n = db.NewNotify(channel)
				}
				h.AddNotify(n)
			}
		}
		db.HostTo = nil
	}
	if db.MaxPrioIn != nil {
		for to, prio := range db.MaxPrioIn {
			notify, ok := db.NotifiesByAddress[to]
			if !ok {
				notify = db.NewNotify(to)
			}
			notify.MaxPrioIn = prio
		}
		db.MaxPrioIn = nil
	}
	if db.RegexIn != nil {
		for to, regexs := range db.RegexIn {
			notify, ok := db.NotifiesByAddress[to]
			if !ok {
				notify = db.NewNotify(to)
			}
			for exp := range regexs {
				notify.AddRegex(exp)
			}
		}
		db.RegexIn = nil
	}
	if db.Lastseen != nil {
		for addr, t := range db.Lastseen {
			h, ok := db.HostsByAddress[addr]
			if !ok {
				h = db.NewHost(addr)
			}
			h.Lastseen = t
		}
		db.Lastseen = nil
	}
	if db.LastseenNotify != nil {
		for addr, t := range db.LastseenNotify {
			h, ok := db.HostsByAddress[addr]
			if !ok {
				h = db.NewHost(addr)
			}
			h.LastseenNotify = t
		}
		db.LastseenNotify = nil
	}
}
