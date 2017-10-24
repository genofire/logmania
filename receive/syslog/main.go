package syslog

import (
	"net"

	log "github.com/sirupsen/logrus"

	"github.com/genofire/logmania/lib"
	"github.com/genofire/logmania/receive"
)

var logger = log.WithField("receive", "syslog")

type Receiver struct {
	receive.Receiver
	exportChannel chan *log.Entry
	serverSocket  *net.UDPConn
}

func Init(config *lib.ReceiveConfig, exportChannel chan *log.Entry) receive.Receiver {
	addr, err := net.ResolveUDPAddr(config.Syslog.Type, config.Syslog.Address)
	ln, err := net.ListenUDP(config.Syslog.Type, addr)

	if err != nil {
		logger.Error("init ", err)
		return nil
	}
	recv := &Receiver{
		serverSocket:  ln,
		exportChannel: exportChannel,
	}

	logger.Info("init")

	return recv
}

const maxDataGramSize = 8192

func (rc *Receiver) Listen() {
	logger.Info("listen")
	for {
		buf := make([]byte, maxDataGramSize)
		n, src, err := rc.serverSocket.ReadFromUDP(buf)
		if err != nil {
			logger.Warn("failed to accept connection", err)
			continue
		}

		raw := make([]byte, n)
		copy(raw, buf)
		entry := toLogEntry(raw, src.IP.String())
		if entry != nil {
			rc.exportChannel <- entry
		}
	}
}

func (rc *Receiver) Close() {
	rc.serverSocket.Close()
}

func init() {
	receive.AddReceiver("syslog", Init)
}
