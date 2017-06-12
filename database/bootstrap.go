package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/genofire/logmania/log"
)

var (
	dbType, connect string
	db              *gorm.DB
)

func Connect(initDBType, initConnect string) {
	var err error
	db, err = gorm.Open(initDBType, initConnect)
	if err != nil {
		log.Panic("failed to connect to database", err)
	}
	bootstrap()
	dbType = initDBType
	connect = initConnect
}

func ReplaceConnect(initDBType, initConnect string) bool {
	if dbType == initDBType && connect == initConnect {
		return false
	}
	dbTemp, err := gorm.Open(initDBType, initConnect)
	if err != nil {
		log.Error("failed to setup new database connection", err)
		return false
	}

	err = db.Close()
	if err != nil {
		log.Error("failed to close old database connection", err)
		return false
	}
	db = dbTemp
	bootstrap()
	dbType = initDBType
	connect = initConnect
	return true
}

func bootstrap() {

	var user User
	var app Application
	db.AutoMigrate(&user)
	db.AutoMigrate(&app)
	db.AutoMigrate(&Entry{})
	if resultUser := db.First(&user); resultUser.RecordNotFound() {
		user.Name = "root"
		if resultApp := db.First(app); resultApp.RecordNotFound() {
			app.Name = "TestSoftware"
			app.Token = "example"
			db.Create(&app)
			user.Permissions = []Application{app}
		}
		db.Create(&user)
	}
}
