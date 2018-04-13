package lib

// Struct of the configuration
// e.g. under dev.sum7.eu/genofire/logmania/logmania_example.conf
type Config struct {
	Notify      NotifyConfig  `toml:"notify"`
	Receive     ReceiveConfig `toml:"receive"`
	DB          string        `toml:"database"`
	HTTPAddress string        `toml:"http_address"`
}

type NotifyConfig struct {
	AlertCheck Duration `toml:"alert_check"`
	Console    bool     `toml:"debug"`
	XMPP       struct {
		JID      string `toml:"jid"`
		Password string `toml:"password"`
	} `toml:"xmpp"`
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
