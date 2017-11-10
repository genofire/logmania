package lib

// Struct of the configuration
// e.g. under github.com/genofire/logmania/logmania_example.conf
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
		Host          string `toml:"host"`
		Username      string `toml:"username"`
		Password      string `toml:"password"`
		Debug         bool   `toml:"debug"`
		NoTLS         bool   `toml:"no_tls"`
		Session       bool   `toml:"session"`
		Status        string `toml:"status"`
		StatusMessage string `toml:"status_message"`
		StartupNotify string `toml:"startup_notify"`
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
