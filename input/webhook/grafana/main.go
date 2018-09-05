package grafana

import (
	"fmt"
	"net/http"
	"time"

	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"

	"dev.sum7.eu/genofire/logmania/input/webhook"
)

type evalMatch struct {
	Tags   map[string]string `mapstructure:"tags,omitempty"`
	Metric string            `mapstructure:"metric"`
	Value  float64           `mapstructure:"value"`
}

type requestBody struct {
	Title       string      `mapstructure:"title"`
	State       string      `mapstructure:"state"`
	RuleID      int64       `mapstructure:"ruleId"`
	RuleName    string      `mapstructure:"ruleName"`
	RuleURL     string      `mapstructure:"ruleUrl"`
	EvalMatches []evalMatch `mapstructure:"evalMatches"`
	ImageURL    string      `mapstructure:"imageUrl"`
	Message     string      `mapstructure:"message"`
}

const webhookType = "grafana"

var HookstateMap = map[string]log.Level{
	"no_data":  log.ErrorLevel,
	"paused":   log.InfoLevel,
	"alerting": log.ErrorLevel,
	"ok":       log.InfoLevel,
	"pending":  log.WarnLevel,
}

var logger = log.WithField("input", webhook.InputType).WithField("hook", webhookType)

func handler(_ http.Header, body interface{}) *log.Entry {
	var request requestBody
	if err := mapstructure.Decode(body, &request); err != nil {
		logger.Warnf("not able to decode data: %s", err)
		return nil
	}
	if request.RuleURL == "" {
		return nil
	}

	entry := log.NewEntry(nil)
	entry = entry.WithField("hostname", request.RuleURL)
	entry.Time = time.Now()
	entry.Level = HookstateMap[request.State]
	if request.Message == "" {
		entry.Message = fmt.Sprintf("%s - %s: %s", request.Title, request.State, request.RuleURL)
	} else {
		entry.Message = fmt.Sprintf("%s - %s: %s %s", request.Title, request.State, request.Message, request.RuleURL)
	}
	return entry
}

func init() {
	webhook.AddHandler(webhookType, handler)
}
