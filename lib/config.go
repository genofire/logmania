package lib

// Struct of the configuration
// e.g. under dev.sum7.eu/genofire/logmania/logmania_example.conf
type Config struct {
	Notify  NotifyConfig  `toml:"notify"`
	Receive ReceiveConfig `toml:"receive"`
	DB      string        `toml:"database"`
}

type NotifyConfig struct {
	AlertCheck Duration `toml:"alert_check"`
	Console    bool     `toml:"debug"`
	XMPP       struct {
		JID      string          `toml:"jid"`
		Password string          `toml:"password"`
		Defaults map[string]bool `toml:"default"`
	} `toml:"xmpp"`
	Websocket struct {
		Address string `toml:"address"`
		Webroot string `toml:"webroot"`
		Default string `toml:"default"`
	} `toml:"websocket"`
	File struct {
		Directory string `toml:"directory"`
		Default   string `toml:"default"`
	} `toml:"file"`
}

type ReceiveConfig struct {
	Syslog struct {
		Type    string `toml:"type"`
		Address string `toml:"address"`
	} `toml:"syslog"`
	JournaldJSON struct {
		Type    string `toml:"type"`
		Address string `toml:"address"`
	} `toml:"journald_json"`
}
