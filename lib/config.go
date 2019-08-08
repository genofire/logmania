package lib

// Struct of the configuration
// e.g. under dev.sum7.eu/sum7/logmania/logmania_example.conf
type Config struct {
	Debug      bool                   `toml:"debug"`
	DB         string                 `toml:"database"`
	AlertCheck Duration               `toml:"alert_check"`
	Output     map[string]interface{} `toml:"output"`
	Input      map[string]interface{} `toml:"input"`
}
