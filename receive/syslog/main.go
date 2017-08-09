package syslog

import (
	"gopkg.in/mcuadros/go-syslog.v2"

	"github.com/genofire/logmania/lib"
	"github.com/genofire/logmania/log"
	"github.com/genofire/logmania/receive"
)

type Receiver struct {
	channel       syslog.LogPartsChannel
	exportChannel chan *log.Entry
	server        *syslog.Server
	receive.Receiver
}

func Init(config *lib.ReceiveConfig, exportChannel chan *log.Entry) receive.Receiver {
	channel := make(syslog.LogPartsChannel)
	handler := syslog.NewChannelHandler(channel)

	server := syslog.NewServer()
	server.SetFormat(syslog.RFC5424)
	server.SetHandler(handler)
	server.ListenUDP(config.Syslog.Bind)

	log.Info("syslog binded to: ", config.Syslog.Bind)

	return &Receiver{
		channel:       channel,
		server:        server,
		exportChannel: exportChannel,
	}
}

func (rc *Receiver) Listen() {
	rc.server.Boot()
	log.Info("boot syslog")
	go func(channel syslog.LogPartsChannel) {
		for logParts := range channel {
			rc.exportChannel <- toLogEntry(logParts)
		}
	}(rc.channel)
}

func (rc *Receiver) Close() {
	rc.server.Kill()
	rc.server.Wait()
}

func init() {
	receive.AddReceiver("syslog", Init)
}
