package xmpp

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/bdlm/log"
)

var tempLog = template.Must(template.New("log").Parse(
	"{{$color := .Color}}<span>" +
		// Hostname
		"{{ if .Hostname }}<span style=\"color: {{ $color.Hostname }};\">{{ .Hostname}}</span>{{ end }}" +
		// Level
		"<span style=\"font-weight: bold; color: {{$color.Level}};\">{{printf \" %5s\" .Level}}</span>" +
		// Message
		"{{printf \" %s\" .Message}}" +
		// Data
		"{{if .Data}}{{range $k, $v := .Data}}" +
		"<span style=\"color: {{$color.DataLabel}};\">{{printf \" %s\" $k}}</span>" +
		"=" +
		"<span style=\"color: {{$color.DataValue}};\">{{$v}}</span>" +
		"{{end}}{{end}}" +
		"</span>",
))

var (
	// DEFAULTColor is the default html 'level' color.
	DEFAULTColor = "#00ff00"
	// ERRORColor is the html 'level' color for error messages.
	ERRORColor = "#ff8700"
	// FATALColor is the html 'level' color for fatal messages.
	FATALColor = "#af0000"
	// PANICColor is the html 'level' color for panic messages.
	PANICColor = "#ff0000"
	// WARNColor is the html 'level' color for warning messages.
	WARNColor = "#ffff00"
	// DEBUGColor is the html 'level' color for debug messages.
	DEBUGColor = "#8a8a8a"

	// DataLabelColor is the html data label color.
	DataLabelColor = "#87afff"
	// DataValueColor is the html data value color.
	DataValueColor = "#d7af87"
	// HostnameColor is the html hostname color.
	HostnameColor = "#00afff"
	// TimestampColor is the html timestamp color.
	TimestampColor = "#5faf87"
)

type logData struct {
	Color     colors                 `json:"-"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Hostname  string                 `json:"host,omitempty"`
	Level     string                 `json:"level,omitempty"`
	Message   string                 `json:"msg,omitempty"`
	Timestamp string                 `json:"time,omitempty"`
}

type colors struct {
	DataLabel string
	DataValue string
	Hostname  string
	Level     string
	Reset     string
	Timestamp string
}

func formatLog(entry *log.Entry) (string, string) {
	var levelColor string

	var logLine *bytes.Buffer
	if entry.Buffer != nil {
		logLine = entry.Buffer
	} else {
		logLine = &bytes.Buffer{}
	}

	data := &logData{
		Data:      make(map[string]interface{}),
		Level:     log.LevelString(entry.Level),
		Message:   entry.Message,
		Timestamp: entry.Time.Format(log.RFC3339Milli),
	}
	switch entry.Level {
	case log.DebugLevel:
		levelColor = DEBUGColor
	case log.WarnLevel:
		levelColor = WARNColor
	case log.ErrorLevel:
		levelColor = ERRORColor
	case log.FatalLevel:
		levelColor = FATALColor
	case log.PanicLevel:
		levelColor = PANICColor
	default:
		levelColor = DEFAULTColor
	}
	data.Color = colors{
		DataLabel: DataLabelColor,
		DataValue: DataValueColor,
		Hostname:  HostnameColor,
		Level:     levelColor,
		Timestamp: TimestampColor,
	}

	for k, v := range entry.Data {
		if k == "hostname" {
			if data.Hostname == "" {
				data.Hostname = v.(string)
			}
			continue
		}
		if str, ok := v.(string); ok {
			data.Data[k] = "'" + str + "'"
		} else {
			data.Data[k] = v
		}
	}

	if err := tempLog.Execute(logLine, data); err != nil {
		return "formating error", "formating error"
	}
	return logLine.String(), fmt.Sprintf("[%s] %s > %s", data.Hostname, log.LevelString(entry.Level), entry.Message)
}
