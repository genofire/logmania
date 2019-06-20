package database

import (
	"dev.sum7.eu/genofire/golang-lib/file"
	"github.com/bdlm/log"
)

func ReadDBFile(path string) *DB {
	var db DB

	if err := file.ReadJSON(path, &db); err == nil {

		db.InitNotify()
		db.InitHost()
		// import
		db.update()
		log.Infof("loaded %d hosts and %d notifies", len(db.Hosts), len(db.Notifies))
		return &db
	} else {
		log.Error("failed to open db file: ", path, ":", err)
	}
	adb := &DB{}
	adb.InitNotify()
	adb.InitHost()
	return adb
}
