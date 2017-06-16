package database

import (
	"encoding/json"
	"time"

	"github.com/genofire/logmania/log"
)

type Entry struct {
	ID            int `json:"id"`
	Time          time.Time
	ApplicationID int
	Fields        string `sql:"type:json"`
	Text          string
	Level         int
}

func transformToDB(dbEntry *log.Entry) *Entry {
	jsonData, err := json.Marshal(dbEntry.Fields)
	if err != nil {
		return nil
	}
	return &Entry{
		Level:  int(dbEntry.Level),
		Text:   dbEntry.Text,
		Fields: string(jsonData),
	}
}

func InsertEntry(token string, entryLog *log.Entry) *Entry {
	app := Application{}
	db.Where("token = ?", token).First(&app)
	entry := transformToDB(entryLog)
	entry.Time = time.Now()
	entry.ApplicationID = app.ID
	result := db.Create(&entry)
	if result.Error != nil {
		log.Error("saving log entry to database", result.Error)
	}
	return entry
}
