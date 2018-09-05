package circleci

import (
	"fmt"
	"net/http"
	"time"

	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"

	"dev.sum7.eu/genofire/logmania/input/webhook"
)

type requestBody struct {
	Payload struct {
		VCSURL    string  `mapstructure:"vcs_url"`
		Status    string  `mapstructure:"status"`
		BuildNum  float64 `mapstructure:"build_num"`
		BuildURL  string  `mapstructure:"build_url"`
		BuildTime float64 `mapstructure:"build_time_millis"`
		Subject   string  `mapstructure:"subject"`
	} `mapstructure:"payload"`
}

const webhookType = "circleci"

var HookstatusMap = map[string]log.Level{
	"failed":  log.ErrorLevel,
	"success": log.InfoLevel,
}

var logger = log.WithField("input", webhook.InputType).WithField("hook", webhookType)

func handler(_ http.Header, body interface{}) *log.Entry {
	var request requestBody
	if err := mapstructure.Decode(body, &request); err != nil {
		logger.Warnf("not able to decode data: %s", err)
		return nil
	}

	if request.Payload.VCSURL == "" {
		return nil
	}

	entry := log.NewEntry(nil)
	entry = entry.WithField("hostname", request.Payload.VCSURL)
	entry.Time = time.Now()
	entry.Level = HookstatusMap[request.Payload.Status]
	entry.Message = fmt.Sprintf("#%0.f (%0.fs): %s - %s", request.Payload.BuildNum, request.Payload.BuildTime/1000, request.Payload.Subject, request.Payload.BuildURL)
	return entry
}

func init() {
	webhook.AddHandler(webhookType, handler)
}
