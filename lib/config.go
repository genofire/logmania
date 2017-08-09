package lib

import (
	"io/ioutil"

	"github.com/BurntSushi/toml"

	"github.com/genofire/logmania/log"
)

// Struct of the configuration
// e.g. under github.com/genofire/logmania/logmania_example.conf
type Config struct {
	Notify  NotifyConfig  `toml:"notify"`
	Receive ReceiveConfig `toml:"receive"`
}

type NotifyConfig struct {
	XMPP struct {
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
	IRC struct {
	} `toml:"irc"`
}

type ReceiveConfig struct {
	Syslog struct {
		Bind string `toml:"bind"`
	} `toml:"syslog"`
}

// read configuration from a file (use toml as file-format)
func ReadConfig(path string) (*Config, error) {
	log.Infof("load of configfile: %s", path)
	var config Config
	file, _ := ioutil.ReadFile(path)
	err := toml.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
