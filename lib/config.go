package lib

import (
	"io/ioutil"

	"github.com/BurntSushi/toml"
	"github.com/genofire/logmania/log"
)

type Config struct {
	API struct {
		Bind        string `toml:"bind"`
		Interactive bool   `toml:"interactive"`
	} `toml:"api"`
	Database struct {
		Type    string `toml:"type"`
		Connect string `toml:"connect"`
	} `toml:"database"`
	Webserver struct {
		Enable bool   `toml:"enable"`
		Bind   string `toml:"bind"`
	} `toml:"webserver"`
}

func ReadConfig(path string) (*Config, error) {
	log.Debugf("load of configfile: %s", path)
	var config Config
	file, _ := ioutil.ReadFile(path)
	err := toml.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
