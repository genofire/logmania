package git

import (
	"net/http"
	"time"

	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"

	"dev.sum7.eu/genofire/logmania/input/webhook"
)

type requestBody struct {
	Repository struct {
		HTMLURL  string `mapstructure:"html_url"`
		FullName string `mapstructure:"full_name"`
	} `mapstructure:"repository"`
	//push
	Pusher struct {
		Name string `mapstructure:"name"`
	} `mapstructure:"pusher"`
	Commits []struct {
		Added    []interface{} `mapstructure:"added"`
		Removed  []interface{} `mapstructure:"removed"`
		Modified []interface{} `mapstructure:"modified"`
	} `mapstructure:"commits"`
	Compare string `mapstructure:"compare"`
	Ref     string `mapstructure:"ref"`
	// issue + fallback
	Sender struct {
		Login string `mapstructure:"login"`
	} `mapstructure:"sender"`
	// issue
	Action string `mapstructure:"action"`
	Issue  struct {
		HTMLURL string  `mapstructure:"html_url"`
		Number  float64 `mapstructure:"number"`
		Title   string  `mapstructure:"title"`
	} `mapstructure:"issue"`
}

const webhookType = "git"

var eventHeader = []string{"X-GitHub-Event", "X-Gogs-Event"}

var logger = log.WithField("input", webhook.InputType).WithField("hook", webhookType)

func handler(header http.Header, body interface{}) *log.Entry {
	event := ""
	for _, head := range eventHeader {
		event = header.Get(head)

		if event != "" {
			break
		}
	}

	if event == "status" {
		return nil
	}
	var request requestBody
	if err := mapstructure.Decode(body, &request); err != nil {
		logger.Warnf("not able to decode data: %s", err)
		return nil
	}

	if request.Repository.HTMLURL == "" {
		return nil
	}

	entry := log.NewEntry(nil)
	entry = entry.WithField("hostname", request.Repository.HTMLURL)
	entry.Time = time.Now()
	entry.Level = log.InfoLevel
	entry.Message = RequestToString(event, request)
	return entry
}

func init() {
	webhook.AddHandler(webhookType, handler)
}
