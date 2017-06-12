package database

import "github.com/genofire/logmania/log"

type User struct {
	ID                  int `json:"id"`
	Name                string
	Mail                string
	XMPP                string
	NotifyMail          bool
	NotifyXMPP          bool
	NotifyAfterLoglevel log.LogLevel
	Permissions         []Application `gorm:"many2many:user_permissions;"`
}
