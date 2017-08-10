package syslog

import (
	"net"

	"github.com/genofire/logmania/lib"
	"github.com/genofire/logmania/log"
	"github.com/genofire/logmania/receive"
)

type Receiver struct {
	receive.Receiver
	exportChannel chan *log.Entry
	serverSocket  *net.UDPConn
}

func Init(config *lib.ReceiveConfig, exportChannel chan *log.Entry) receive.Receiver {
	addr, err := net.ResolveUDPAddr(config.Syslog.Type, config.Syslog.Address)
	ln, err := net.ListenUDP(config.Syslog.Type, addr)

	if err != nil {
		log.Error("syslog init ", err)
		return nil
	}
	recv := &Receiver{
		serverSocket:  ln,
		exportChannel: exportChannel,
	}

	log.Info("syslog init")

	return recv
}

const maxDataGramSize = 8192

func (rc *Receiver) Listen() {
	log.Info("syslog listen")
	for {
		buf := make([]byte, maxDataGramSize)
		n, src, err := rc.serverSocket.ReadFromUDP(buf)
		if err != nil {
			log.Warn("failed to accept connection", err)
			continue
		}

		raw := make([]byte, n)
		copy(raw, buf)
		rc.exportChannel <- toLogEntry(raw, src.IP.String())
	}
}

func (rc *Receiver) Close() {
	rc.serverSocket.Close()
}

func init() {
	receive.AddReceiver("syslog", Init)
}
