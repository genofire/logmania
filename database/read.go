package database

import (
	"dev.sum7.eu/genofire/golang-lib/file"
	log "github.com/sirupsen/logrus"
)

func ReadDBFile(path string) *DB {
	var db DB

	if err := file.ReadJSON(path, &db); err == nil {
		log.Infof("loaded %d hosts", len(db.HostTo))

		db.InitNotify()
		db.InitHost()
		// import
		db.update()
		return &db
	} else {
		log.Error("failed to open db file: ", path, ":", err)
	}
	adb := &DB{}
	adb.InitNotify()
	adb.InitHost()
	return adb
}
